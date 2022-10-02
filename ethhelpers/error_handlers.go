package ethhelpers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

// DefaultErrorHandlerWithLogger is used by e.g. WaitForTransactionReceipt to
// decide if the error is a temporary connection / receipt not found issue and
// it should retry.
//
// This is a work-in-progress.
func DefaultErrorHandlerWithMessages(msgHandler func(txHash common.Hash, msg string)) func(txHash common.Hash, err error) error {
	return func(txHash common.Hash, err error) error {
		var rpcErr rpc.Error

		switch {
		case errors.As(err, &rpcErr):
			// TODO: We should have optional permissive error handling code for
			// temporary errors.
			if err.Error() == "request failed or timed out" {
				msgHandler(txHash, fmt.Sprintf("waiting for transaction receipt, temporary error : %d : %v", rpcErr.ErrorCode(), rpcErr))
				break
			}

			msgHandler(txHash, fmt.Sprintf("waiting for transaction receipt, unknown rpc error : %d : %v", rpcErr.ErrorCode(), rpcErr))
			return err

		case errors.Is(err, context.Canceled):
			msgHandler(txHash, fmt.Sprintf("waiting for transaction receipt, context canceled"))

		case errors.Is(err, context.DeadlineExceeded):
			msgHandler(txHash, fmt.Sprintf("waiting for transaction receipt, context deadline"))

		case errors.Is(err, ethereum.NotFound):
			msgHandler(txHash, fmt.Sprintf("waiting for transaction receipt, not found"))

		case errors.Is(err, syscall.ECONNRESET):
			msgHandler(txHash, fmt.Sprintf("waiting for transaction receipt, connection reset by peer"))

		case os.IsTimeout(err):
			msgHandler(txHash, fmt.Sprintf("waiting for transaction receipt, timeout"))

		// TODO: Check non-temporary errors.
		default:
			msgHandler(txHash, fmt.Sprintf("waiting for transaction receipt, unknown error : %v", err))
			return err
		}

		return nil
	}
}
