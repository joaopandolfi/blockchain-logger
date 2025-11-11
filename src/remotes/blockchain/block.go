package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const TRANSACTION_CODE_SYSTEM_ID = "system_id"
const GENESIS_HASH_BLOCK = "01f4913e4f39713b5d2260b443ff30c2b696256d3488731a8b4587ab3fb6983f"
const GENESIS_ID_BLOCK = "6ec9d09f-fee4-494c-9309-f603f275f4df"

type Block struct {
	ID uuid.UUID `gorm:"primarykey"`

	// Used to link this block to the previous
	LastBlockID uuid.UUID

	// Used to link this block to the previous and check the consistency
	LastBlockHash string

	// TransactionStr its the payload to be hashed and signed
	TransactionStr string `gorm:"transaction"`
	// Transaction its the variable to be manipulated
	Transaction map[string]interface{} `gorm:"-"`

	// Metadata used to filter blocks by a system in database
	SystemID string

	// Unverifyed metadata
	// Use only to mark a block in a chain
	Tags string

	// Squential id to incremented when added in a chain
	// Used to verify the block position in a chain
	SeqID uint `gorm:"index"`

	// Hash value from the payload and others metadata
	// Used to verify if the block its consistent
	Hash string

	// Used to sign a block and keeping then trustable
	Signature string
	CreatedAt time.Time
	UpdatedAt time.Time
	SignedAt  time.Time
	HashedAt  string
}

func NewBlock(systemID string, transaction map[string]interface{}, tags ...string) *Block {
	return &Block{
		ID:          uuid.New(),
		Transaction: transaction,
		SystemID:    systemID,
		Tags:        strings.Join(tags, ";"),
	}
}

// Hashable - returns the string to be hashed
// Uses the variables initialized on HashBlock()
func (b *Block) Hashable() string {
	return fmt.Sprintf("%d.%s.%s.%s.%s.%s", b.SeqID, b.ID.String(), b.LastBlockID.String(), b.LastBlockHash, b.TransactionStr, b.HashedAt)
}

// Signable - returns the values to be signed by chain
// Important the block be Hashed first
func (b *Block) Signable() string {
	return fmt.Sprintf("%s.%s", b.ID.String(), b.Hash)
}

// HashBlock - calc the hash for block and add metadata inside them
// After be hashed the block can be signed and will be valid
func (b *Block) HashBlock() error {

	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	if b.Transaction == nil {
		b.Transaction = map[string]interface{}{}
	}

	b.Transaction[TRANSACTION_CODE_SYSTEM_ID] = b.SystemID

	t, err := json.Marshal(&b.Transaction)
	if err != nil {
		return fmt.Errorf("marshaling transaction: %w", err)
	}

	b.TransactionStr = string(t)
	b.HashedAt = time.Now().Format(time.RFC3339Nano)

	if b.ID.String() == GENESIS_ID_BLOCK {
		b.HashedAt = time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC).Format(time.RFC3339Nano)
	}

	b.Hash = b.CalcHash()

	return nil
}

// CalcHash - calc the hash to the block
func (b *Block) CalcHash() string {
	hashable := b.Hashable()

	hash := sha256.Sum256([]byte(hashable))
	return fmt.Sprintf("%x", hash[:])
}

// Unpack - unpack the payload stored on transactionString to transaction struct
func (b *Block) Unpack() error {
	err := json.Unmarshal([]byte(b.TransactionStr), &b.Transaction)
	if err != nil {
		return fmt.Errorf("parsing the transaction string into a struct: %w", err)
	}

	return nil
}
