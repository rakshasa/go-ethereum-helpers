package ethtesting

import (
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/log"
)

type pendingLogHandler struct {
	mu      sync.Mutex
	records []*log.Record
}

func (h *pendingLogHandler) Log(r *log.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.records = append(h.records, r)
	return nil
}

// PendingLogHandler returns a new log.Handler that buffers records until
// the returned commit function is called.
func PendingLogHandler(h log.Handler) (log.Handler, func()) {
	p := &pendingLogHandler{
		records: []*log.Record{},
	}

	commit := func() {
		p.mu.Lock()
		defer p.mu.Unlock()

		for _, r := range p.records {
			_ = h.Log(r)
		}

		p.records = []*log.Record{}
	}

	return p, commit
}

func PendingLogHandlerForTesting(t *testing.T, logger log.Logger) func() {
	oldHandler := logger.GetHandler()
	handler, commitLogs := PendingLogHandler(oldHandler)

	logger.SetHandler(handler)

	return func() {
		if t.Failed() {
			commitLogs()
		}

		logger.SetHandler(oldHandler)
	}
}
