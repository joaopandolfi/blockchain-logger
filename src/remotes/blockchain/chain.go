package blockchain

import (
	"fmt"
	"time"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/google/uuid"
)

var chain *BlockChain

type BlockChain struct {
	GenesisBlock Block

	Chain      []Block `gorm:"-"`
	privKey    string  `json:"-" gorm:"-"`
	passphrase string  `json:"-" gorm:"-"`
	PubKey     string
}

// InitChain start a chain with ou whithout a block
func InitChain(pubKey string, blocks ...Block) {
	chain = &BlockChain{
		PubKey: pubKey,
		Chain:  blocks,
	}

	if len(blocks) > 0 {
		chain.GenesisBlock = blocks[0]
	}
}

// Get - return a instance of a initialized chain
func Get() *BlockChain {
	if chain == nil {
		panic("chain not initialized")
	}
	return chain
}

// SetAuth - used to allow a blockchain instance sign new blocks
func (b *BlockChain) SetAuth(privKey, passphrase string) {
	b.passphrase = passphrase
	b.privKey = privKey
}

// HaveAuth - to verify if the auth is setted
func (b *BlockChain) HaveAuth() bool {
	return b.privKey != "" && b.passphrase != ""
}

// Checkacble - to verify if a blokchain can verify the signature from the blocks
func (b *BlockChain) Checkable() bool {
	return b.PubKey != ""
}

// unlockPrivKey - retrieve the privkey given a password
func (b *BlockChain) unlockPrivKey() (*crypto.KeyRing, error) {
	passphrase := []byte(b.passphrase) // Private key passphrase

	privateKeyObj, err := crypto.NewKeyFromArmored(b.privKey)
	if err != nil {
		return nil, fmt.Errorf("reading privKey: %w", err)
	}
	unlockedKey, err := privateKeyObj.Unlock(passphrase)
	if err != nil {
		return nil, fmt.Errorf("opening privkey with passphrase: %w", err)
	}

	signingKeyRing, err := crypto.NewKeyRing(unlockedKey)
	if err != nil {
		return nil, fmt.Errorf("creating keyring")
	}
	return signingKeyRing, nil
}

// getPubKey - return a instance of pubkey from the pubkey string
func (b *BlockChain) getPubKey() (*crypto.KeyRing, error) {
	publicKeyObj, err := crypto.NewKeyFromArmored(b.PubKey)
	if err != nil {
		return nil, fmt.Errorf("reading pubKey: %w", err)
	}
	signingKeyRing, err := crypto.NewKeyRing(publicKeyObj)
	if err != nil {
		return nil, fmt.Errorf("generating keyring: %w", err)
	}
	return signingKeyRing, nil
}

// AppendBlock - verify add and sign new block on the chain
// the chain needs to be initialized first, needs to have at last 1 block
func (b *BlockChain) AppendBlock(block *Block) (*Block, error) {
	size := len(b.Chain)
	if size == 0 {
		return nil, fmt.Errorf("empty chain")
	}

	lastBlock := b.Chain[size-1]
	newBlock, err := b.ChainBlocks(&lastBlock, block)
	if err != nil {
		return nil, fmt.Errorf("chaining block: %w", err)
	}

	b.Chain = append(b.Chain, *newBlock)

	return newBlock, nil
}

// ChainBlocks - Given a last block, the function will check the consistensy and signature
// sign a new block and chain to last
// used to chain 2 blocks block without need to have the entire chain, just having a valid one block
func (b *BlockChain) ChainBlocks(lastBlock, block *Block) (*Block, error) {

	if !b.HaveAuth() || !b.Checkable() {
		return nil, fmt.Errorf("not initialized")
	}

	pubKey, err := b.getPubKey()
	if err != nil {
		return nil, fmt.Errorf("getting pubkey: %w", err)
	}

	privKey, err := b.unlockPrivKey()
	if err != nil {
		return nil, fmt.Errorf("unlocking privKey: %w", err)
	}

	err = b.validateBlock(pubKey, lastBlock)
	if err != nil {
		return nil, fmt.Errorf("validating last block: %w", err)
	}

	block.LastBlockHash = lastBlock.Hash
	block.LastBlockID = lastBlock.ID
	block.SeqID = lastBlock.SeqID + 1

	err = block.HashBlock()
	if err != nil {
		return nil, fmt.Errorf("hashing block: %w", err)
	}

	err = b.signBlock(privKey, block)
	if err != nil {
		return nil, fmt.Errorf("signing block: %w", err)
	}

	return block, nil
}

// GenerateGenesis - create a first valid block of the chain
func (b *BlockChain) GenerateGenesis() error {

	if !b.HaveAuth() {
		return fmt.Errorf("not initialized")
	}

	privKey, err := b.unlockPrivKey()
	if err != nil {
		return fmt.Errorf("unlocking privKey: %w", err)
	}

	genesisID, _ := uuid.Parse(GENESIS_ID_BLOCK)
	genesisBlock := &Block{
		ID:            genesisID,
		LastBlockID:   uuid.Nil,
		LastBlockHash: "",
		Transaction:   map[string]interface{}{"there was light": "and become light"},
		SeqID:         0,
		SystemID:      "genesis",
	}
	err = genesisBlock.HashBlock()
	if err != nil {
		return fmt.Errorf("hashing genesis block: %w", err)
	}

	err = b.signBlock(privKey, genesisBlock)
	if err != nil {
		return fmt.Errorf("validating genesis block: %w", err)
	}

	b.GenesisBlock = *genesisBlock
	b.Chain = append(b.Chain, *genesisBlock)

	return nil
}

// Validate the entire chain
func (b *BlockChain) Validate() error {

	if !b.Checkable() {
		return fmt.Errorf("not initiliazed")
	}

	pubKey, err := b.getPubKey()
	if err != nil {
		return fmt.Errorf("getting PubKey: %w", err)
	}

	lastBlock := b.GenesisBlock

	size := len(b.Chain)

	for i := 0; i < size; i++ {
		block := b.Chain[i]

		err = b.validateBlock(pubKey, &block)
		if err != nil {
			return fmt.Errorf("validating block: %w", err)
		}

		if block.ID.String() == GENESIS_ID_BLOCK && block.Hash == GENESIS_HASH_BLOCK {
			continue
		}

		if block.LastBlockHash != lastBlock.Hash {
			return fmt.Errorf("chain is broken, last block is different! blockID: %s seqBlock: %d Comparation: [%s] [%s] LastBLockId: [%s]",
				block.ID.String(),
				block.SeqID,
				block.LastBlockHash,
				lastBlock.Hash,
				lastBlock.ID.String())
		}

		lastBlock = block
	}

	return nil
}

func (b *BlockChain) Clean() {
	b.Chain = nil
	b.GenesisBlock = Block{}
}

func (b *BlockChain) signBlock(privKey *crypto.KeyRing, block *Block) error {

	message := crypto.NewPlainMessage([]byte(block.Signable()))
	pgpSignature, err := privKey.SignDetached(message)
	if err != nil {
		return fmt.Errorf("generate signature: %w", err)
	}

	block.SignedAt = time.Now()
	block.Signature = formatSignatureToPGP(pgpSignature.Data)

	return nil
}

func (b *BlockChain) validateBlock(pubKey *crypto.KeyRing, block *Block) error {

	message := crypto.NewPlainMessage([]byte(block.Signable()))

	pgpSignature, err := crypto.NewPGPSignatureFromArmored(block.Signature)
	if err != nil {
		return fmt.Errorf("opening signature: %w", err)
	}

	err = pubKey.VerifyDetached(message, pgpSignature, crypto.GetUnixTime())
	if err != nil {
		return fmt.Errorf("checking signature: %w", err)
	}

	hash := block.CalcHash()
	if block.Hash != hash {
		return fmt.Errorf("invalid hash block: %v != %v", block.Hash, hash)
	}

	return nil
}
