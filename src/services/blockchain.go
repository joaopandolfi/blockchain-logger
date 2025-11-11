package services

import (
	"context"
	"fmt"
	"logger/models/dao"
	"logger/remotes/blockchain"

	"github.com/google/uuid"
	"github.com/joaopandolfi/blackwhale/remotes/jaeger"
)

type BlockChain interface {
	Validate(ctx context.Context) error
	ValidateSegment(ctx context.Context, init, end int) error
}

type blockChainService struct {
	dao dao.BlockChain
}

func NewBlockChain() BlockChain {
	return &blockChainService{
		dao: dao.NewBlockChainDao(),
	}
}

func (s *blockChainService) Validate(ctx context.Context) error {
	sCtx, tracer := jaeger.SpanTrace(ctx, "service.Validate", nil)
	defer tracer.Finish()

	return s.ValidateSegment(sCtx, 0, 0)
}

func (s *blockChainService) ValidateSegment(ctx context.Context, init, end int) error {
	_, tracer := jaeger.SpanTrace(ctx, "service.ValidateSegment", map[string]interface{}{"init": init, "end": end})
	defer tracer.Finish()

	blocks, err := s.dao.GetSegment(init, end)
	if err != nil {
		return fmt.Errorf("getting blocks (%d, %d): %w", init, end, err)
	}

	chainSize := len(blocks)
	if chainSize == 0 {
		return nil
	}

	chain := blockchain.Get()
	chain.Chain = blocks
	defer chain.Clean()

	if chain.GenesisBlock.ID == uuid.Nil {
		chain.GenesisBlock = chain.Chain[0]
		chain.Chain = chain.Chain[1 : chainSize-1]
	}

	err = chain.Validate()
	if err != nil {
		return fmt.Errorf("validating %d blocks [%d - %d] on chain: %w", init, end, len(chain.Chain), err)
	}

	return nil
}
