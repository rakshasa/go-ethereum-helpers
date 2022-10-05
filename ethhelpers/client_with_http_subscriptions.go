package ethhelpers

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

type clientWithHTTPSubscriptions struct {
	Client
}

// TODO: Add options.

func NewClientWithHTTPSubscriptions(c Client) Client {
	return &clientWithHTTPSubscriptions{
		Client: c,
	}
}

func (c *clientWithHTTPSubscriptions) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	if ch == nil {
		panic("channel given to SubscribeFilterLogs must not be nil")
	}

	if q.FromBlock == nil {
		blockNum, err := c.BlockNumber(ctx)
		if err != nil {
			return nil, err
		}

		q.FromBlock = new(big.Int).SetUint64(blockNum)
	}

	return nil, fmt.Errorf("not implemented")
}
