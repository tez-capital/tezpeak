package providers

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/trilitech/tzgo/rpc"
)

func attemptWithClients[T interface{}](clients []*rpc.Client, f func(client *rpc.Client) (T, error)) (T, error) {
	var err error
	var result T

	for _, client := range clients {
		slog.Debug("attempting with client", "client", client.BaseURL.Host)
		result, err = f(client)
		if err != nil {
			continue
		}
		return result, nil
	}
	return result, err
}

var (
	blockSubscribers         = map[uuid.UUID]chan *rpc.BlockHeaderLogEntry{}
	lastProcessedBlockHeight = int64(0)

	rightsSubscribers         = map[uuid.UUID]chan *RightsStatus{}
	lastProcessedRightsHeight = int64(0)
)

func init() {
	go func() {
		for b := range blockHeaderLogEntryChannel {
			block := b // remove in 1.22
			if block.Level < lastProcessedBlockHeight {
				continue
			}
			lastProcessedBlockHeight = block.Level
			for _, subscriber := range blockSubscribers {
				s := subscriber // remove in 1.22
				go func() {
					s <- block
				}()
			}
		}
	}()

	go func() {
		for r := range rightsChannel {
			rights := r // remove in 1.22
			if rights.Level < lastProcessedRightsHeight {
				continue
			}
			lastProcessedRightsHeight = rights.Level
			for _, subscriber := range rightsSubscribers {
				s := subscriber // remove in 1.22
				go func() {
					s <- rights
				}()
			}
		}
	}()
}
