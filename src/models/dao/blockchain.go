package dao

import (
	"context"
	"fmt"
	"logger/remotes/blockchain"
	"sync"

	"github.com/joaopandolfi/blackwhale/models/dao"
	"github.com/joaopandolfi/blackwhale/remotes/jaeger"
)

type BlockChain interface {
	AppendBlock(ctx context.Context, b *blockchain.Block) (*blockchain.Block, error)
	GetSegment(init, end int) ([]blockchain.Block, error)
	GetAll() ([]blockchain.Block, error)
}

type blockChain struct {
	dao dao.SQLDAO
	mu  sync.Mutex
}

var singleton *blockChain

func NewBlockChainDao() BlockChain {
	if singleton == nil {
		singleton = &blockChain{
			dao: new(),
		}
	}

	return singleton
}

func (s *blockChain) AppendBlock(ctx context.Context, block *blockchain.Block) (*blockchain.Block, error) {
	_, tracer := jaeger.SpanTrace(ctx, "dao.blockchain.AppendBlock", map[string]interface{}{"id": block.ID})
	defer tracer.Finish()

	s.mu.Lock()
	defer s.mu.Unlock()

	var lastBlock blockchain.Block
	err := s.dao.ListAll(&lastBlock, dao.ListParams{
		Limit: 1,
		Order: "created_at desc",
	})
	if err != nil {
		return nil, fmt.Errorf("recovering last block: %w", err)
	}

	newValidBlock, err := blockchain.Get().ChainBlocks(&lastBlock, block)
	if err != nil {
		return nil, fmt.Errorf("adding block in to chain: %w", err)
	}

	err = s.dao.New(&newValidBlock)
	if err != nil {
		return nil, fmt.Errorf("saving new block on database: %w", err)
	}

	return newValidBlock, nil
}

func (s *blockChain) GetAll() ([]blockchain.Block, error) {
	//TODO: Do it in batches
	return s.GetSegment(0, 0) // entire chain
}

func (s *blockChain) GetSegment(init, end int) ([]blockchain.Block, error) {
	var blocks []blockchain.Block

	filters := dao.ListParams{
		Order: "seq_id asc",
	}

	if end != 0 {
		filters.Limit = end - init
		filters.Offset = init
	}

	err := s.dao.ListAll(&blocks, filters)
	if err != nil {
		return nil, fmt.Errorf("getting blocks: %w", err)
	}

	return blocks, nil
}
