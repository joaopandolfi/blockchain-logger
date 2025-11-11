package services

import (
	"context"
	"fmt"
	"logger/models"
	"logger/models/dao"
	"logger/remotes/blockchain"

	"github.com/joaopandolfi/blackwhale/remotes/jaeger"
)

type Logs interface {
	New(ctx context.Context, l *models.Log) (*models.Log, error)
}

type logs struct {
	dao dao.BlockChain
}

func NewLogs() Logs {
	return &logs{
		dao: dao.NewBlockChainDao(),
	}
}

func (s *logs) New(ctx context.Context, l *models.Log) (*models.Log, error) {
	sCtx, tracer := jaeger.SpanTrace(ctx, "service.New", map[string]interface{}{"system": l.SystemID})
	defer tracer.Finish()

	block := blockchain.NewBlock(l.SystemID, l.Payload, l.ParseTags()...)

	signedBlock, err := s.dao.AppendBlock(sCtx, block)
	if err != nil {
		return nil, fmt.Errorf("appending block: %w", err)
	}

	l.Block = signedBlock

	return l, nil
}
