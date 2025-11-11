package blockchain_test

import (
	"encoding/json"
	"fmt"
	"logger/remotes/blockchain"
	"testing"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func _mockGenesisBlock() blockchain.Block {
	genesisID, _ := uuid.Parse(blockchain.GENESIS_ID_BLOCK)
	return blockchain.Block{
		ID:            genesisID,
		LastBlockID:   uuid.Nil,
		LastBlockHash: "",
		Transaction:   map[string]interface{}{"there was light": "and become light"},
		SeqID:         0,
		SystemID:      "genesis",
	}
}

func _mockBlock() *blockchain.Block {
	return &blockchain.Block{
		Transaction: map[string]interface{}{
			"table": "user",
			"from": map[string]interface{}{
				"user_id": "123",
			},
			"to": map[string]interface{}{
				"user_id": "1234",
			},
		},
		SystemID: "sauron",
	}
}

func _generateMockKey(passphrase string) (string, string, error) {
	name := "authority mocked"
	email := "mocked@authority.com"
	rsaBits := 2048
	rsaKeyObj, err := crypto.GenerateKey(name, email, "rsa", rsaBits)
	if err != nil {
		return "", "", fmt.Errorf("generating key: %w", err)
	}
	lockedRsaObj, err := rsaKeyObj.Lock([]byte(passphrase))
	if err != nil {
		return "", "", fmt.Errorf("locking pgp: %w", err)
	}

	pk, _ := lockedRsaObj.GetArmoredPublicKey()
	privK, _ := lockedRsaObj.Armor()

	return privK, pk, nil
}

func TestGeneratingKey(t *testing.T) {

	name := "authority mocked"
	email := "mocked@authority.com"
	passphrase := []byte("Long Long Long Mocked Key Secrets")
	rsaBits := 2048
	// // RSA, Key struct
	rsaKeyObj, err := crypto.GenerateKey(name, email, "rsa", rsaBits)
	assert.Nil(t, err)

	lockedRsaObj, err := rsaKeyObj.Lock(passphrase)
	assert.Nil(t, err)
	assert.NotNil(t, lockedRsaObj)

	pk, err := lockedRsaObj.GetArmoredPublicKey()
	assert.Nil(t, err)

	privK, err := lockedRsaObj.Armor()
	assert.Nil(t, err)

	assert.NotNil(t, pk)
	assert.NotNil(t, privK)
}

func TestGenesisBlock(t *testing.T) {

	genesisBlock := _mockGenesisBlock()

	err := genesisBlock.HashBlock()
	assert.Nil(t, err)
	assert.NotNil(t, genesisBlock)
	assert.NotEmpty(t, genesisBlock.Hash)
}

func TestChain(t *testing.T) {
	pass := "very very long long key"
	privKey, pubKey, err := _generateMockKey(pass)
	assert.Nil(t, err)

	blockchain.InitChain(pubKey)

	chain := blockchain.Get()

	chain.SetAuth(privKey, pass)

	err = chain.GenerateGenesis()
	assert.Nil(t, err)

	addedBlock, err := chain.ChainBlocks(&chain.GenesisBlock, _mockBlock())
	chain.Chain = append(chain.Chain, *addedBlock)

	assert.Nil(t, err)
	assert.NotEmpty(t, addedBlock.Hash)
	assert.NotEmpty(t, addedBlock.Signature)

	addedBlock2, err := chain.ChainBlocks(addedBlock, &blockchain.Block{
		Transaction: map[string]interface{}{
			"table": "user",
			"from": map[string]interface{}{
				"user_id": "555",
			},
			"to": map[string]interface{}{
				"user_id": "555",
			},
		},
		SystemID: "sauron",
	})

	assert.Nil(t, err)
	assert.NotEmpty(t, addedBlock2.Hash)

	chain.Chain = append(chain.Chain, *addedBlock2)

	assert.NotEmpty(t, chain.Chain)
	assert.Equal(t, 3, len(chain.Chain), "1 genesis block and 2 transactionals")

	err = chain.Validate()
	assert.Nil(t, err)
}

func TestChainAppendBlock(t *testing.T) {

	pass := "very very long long key"
	privKey, pubKey, err := _generateMockKey(pass)
	assert.Nil(t, err)

	blockchain.InitChain(pubKey)

	chain := blockchain.Get()

	chain.SetAuth(privKey, pass)

	err = chain.GenerateGenesis()
	assert.Nil(t, err)

	addedBlock, err := chain.AppendBlock(_mockBlock())

	assert.Nil(t, err)
	assert.NotEmpty(t, addedBlock.Hash)
	assert.NotEmpty(t, addedBlock.Signature)

	assert.Equal(t, 2, len(chain.Chain), "1 genesis block and 1 transactional")

	err = chain.Validate()
	assert.Nil(t, err)
}

func TestChainBlockFraudSignature(t *testing.T) {
	pass := "very very long long key"
	privKey, pubKey, err := _generateMockKey(pass)
	assert.Nil(t, err)

	blockchain.InitChain(pubKey)

	chain := blockchain.Get()

	chain.SetAuth(privKey, pass)

	err = chain.GenerateGenesis()
	assert.Nil(t, err)

	addedBlock, err := chain.AppendBlock(_mockBlock())

	assert.Nil(t, err)
	assert.NotEmpty(t, addedBlock.Hash)
	assert.NotEmpty(t, addedBlock.Signature)

	assert.Equal(t, 2, len(chain.Chain), "1 genesis block and 1 transactional")

	block1 := chain.Chain[1]
	block1.Transaction["table"] = "admin"
	err = block1.HashBlock()
	assert.Nil(t, err)
	chain.Chain[1] = block1

	err = chain.Validate()
	assert.ErrorContains(t, err, "Signature Verification Error")
}

func TestChainBlockFraudHash(t *testing.T) {
	pass := "very very long long key"
	privKey, pubKey, err := _generateMockKey(pass)
	assert.Nil(t, err)

	blockchain.InitChain(pubKey)

	chain := blockchain.Get()

	chain.SetAuth(privKey, pass)

	err = chain.GenerateGenesis()
	assert.Nil(t, err)

	addedBlock, err := chain.AppendBlock(_mockBlock())

	assert.Nil(t, err)
	assert.NotEmpty(t, addedBlock.Hash)
	assert.NotEmpty(t, addedBlock.Signature)

	assert.Equal(t, 2, len(chain.Chain), "1 genesis block and 1 transactional")

	block1 := chain.Chain[1]
	block1.Transaction["table"] = "admin"
	payload, err := json.Marshal(&block1)
	assert.Nil(t, err)

	block1.TransactionStr = string(payload)
	chain.Chain[1] = block1

	err = chain.Validate()
	assert.ErrorContains(t, err, "invalid hash block")
}
