package ethhelpers

import (
	"context"
)

type BlockNumberReader interface {
	BlockNumber(ctx context.Context) (uint64, error)
}
