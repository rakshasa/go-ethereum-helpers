package ethhelpers

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockNumberReader interface {
	BlockNumber(ctx context.Context) (uint64, error)
}

type FilterLogsReader interface {
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
}
