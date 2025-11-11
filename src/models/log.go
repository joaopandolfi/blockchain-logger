package models

import (
	"logger/remotes/blockchain"
	"strings"

	"github.com/google/uuid"
)

const TAG_SEPARATOR = ";"

type Log struct {
	ID       uuid.UUID
	Payload  map[string]interface{}
	SystemID string
	Tags     string
	Block    *blockchain.Block
}

func (m *Log) ParseTags() []string {
	return strings.Split(m.Tags, TAG_SEPARATOR)
}
