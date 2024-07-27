package tezbake

import (
	"context"

	"github.com/tez-capital/tezpeak/core/common"
	"github.com/trilitech/tzgo/rpc"
)

type LedgerStatus struct {
	Level int64 `json:"level"`
}

func StartLedgerStatusProvider(ctx context.Context, clients []*rpc.Client, bakers []string, window int64, statusChannel chan<- common.ModuleStatusUpdate) {
	// not implemented
}
