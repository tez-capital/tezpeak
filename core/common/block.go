package common

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/trilitech/tzgo/rpc"
)

type Block struct {
	Hash      string    `json:"hash"`
	Timestamp time.Time `json:"timestamp"`
	//	Fitness          string                `json:"fitness"` add if relevant
	Protocol         string                `json:"protocol"`
	LevelInfo        *rpc.LevelInfo        `json:"level_info"`
	VotingPeriodInfo *rpc.VotingPeriodInfo `json:"voting_period_info"`
}

type ConnectionStatus string

const (
	Connected    ConnectionStatus = "connected"
	Connecting   ConnectionStatus = "connecting"
	Disconnected ConnectionStatus = "disconnected"
)

type blockMonitor struct {
	client *rpc.Client
	cancel context.CancelFunc
}

type BlockEventSource struct {
	*EventSource[*rpc.BlockHeaderLogEntry]

	blockMonitors    map[uuid.UUID]blockMonitor
	blockMonitorsMtx sync.RWMutex
}

func createMonitor(ctx context.Context, client *rpc.Client) (*rpc.BlockHeaderMonitor, error) {
	mon := rpc.NewBlockHeaderMonitor()
	if err := client.MonitorBlockHeader(ctx, mon); err != nil {
		return nil, err
	}
	slog.Debug("created block monitor", "source", client.BaseURL)
	return mon, nil
}

func (es *BlockEventSource) runBlockHeaderMonitor(ctx context.Context, client *rpc.Client, connectionStateCallback func(state ConnectionStatus), blockCallback func(*Block)) {
	if connectionStateCallback == nil {
		connectionStateCallback = func(state ConnectionStatus) {}
	}

	var err error
	var mon *rpc.BlockHeaderMonitor
	for {
		if mon == nil {
			mon, err = createMonitor(ctx, client)
			connectionStateCallback(Connecting)
			if err != nil {
				slog.Debug("failed to create block monitor", "source", client.BaseURL.String(), "error", err.Error())
				connectionStateCallback(Disconnected)
				select {
				case <-ctx.Done():
					return
				case <-time.After(5 * time.Second):
					// Wait 5 seconds before retrying
				}
			} else {
				connectionStateCallback(Connected)
			}
			continue
		}

		select {
		case <-ctx.Done():
			return
		default:
			h, err := mon.Recv(ctx)
			if err == rpc.ErrMonitorClosed {
				slog.Debug("block monitor closed", "source", client.BaseURL.String())
				connectionStateCallback(Disconnected)
				mon = nil
				select {
				case <-ctx.Done():
					return
				case <-time.After(5 * time.Second):
					// Wait 5 seconds before retrying
				}
				continue
			}
			if err != nil && h == nil {
				slog.Debug("failed to receive block", "source", client.BaseURL.String(), "error", err.Error())
				connectionStateCallback(Disconnected)
				mon = nil
				select {
				case <-ctx.Done():
					break
				case <-time.After(5 * time.Second):
					// Wait 5 seconds before retrying
				}
				continue
			}
			go func() {
				es.source <- h
			}()

			go func() {
				metadata, err := client.GetBlockMetadata(ctx, h.Hash)
				if err != nil {
					slog.Debug("failed to get block metadata", "source", client.BaseURL.String(), "error", err.Error())
					return
				}

				blockCallback(&Block{
					Hash:             h.Hash.String(),
					Timestamp:        h.Timestamp,
					Protocol:         metadata.Protocol.String(),
					LevelInfo:        metadata.LevelInfo,
					VotingPeriodInfo: metadata.VotingPeriodInfo,
				})
			}()
		}
	}

}

func (es *BlockEventSource) AddBlockMonitor(ctx context.Context, client *rpc.Client, connectionStateCallback func(state ConnectionStatus), blockCallback func(*Block)) (uuid.UUID, error) {
	var err error

	id, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, err
	}
	blockMonitorContext, cancel := context.WithCancel(ctx)

	go es.runBlockHeaderMonitor(blockMonitorContext, client, connectionStateCallback, blockCallback)

	es.blockMonitorsMtx.Lock()
	defer es.blockMonitorsMtx.Unlock()
	es.blockMonitors[id] = blockMonitor{
		client: client,
		cancel: cancel,
	}
	return id, err
}

func (es *BlockEventSource) RemoveBlockMonitor(id uuid.UUID) {
	es.blockMonitorsMtx.Lock()
	defer es.blockMonitorsMtx.Unlock()
	if monitor, ok := es.blockMonitors[id]; ok {
		monitor.cancel()
		delete(es.blockMonitors, id)
	}
}

func NewBlockEventSource() *BlockEventSource {
	return &BlockEventSource{
		EventSource:   NewEventSource[*rpc.BlockHeaderLogEntry](nil),
		blockMonitors: make(map[uuid.UUID]blockMonitor),
	}
}

var (
	blockEventSource = NewBlockEventSource()
)

func init() {
	go blockEventSource.Run()
}

func AddBlockMonitor(ctx context.Context, client *rpc.Client, connectionStateCallback func(ConnectionStatus), blockCallback func(*Block)) (uuid.UUID, error) {
	return blockEventSource.AddBlockMonitor(ctx, client, connectionStateCallback, blockCallback)
}

func RemoveBlockMonitor(id uuid.UUID) {
	blockEventSource.RemoveBlockMonitor(id)
}

func SubscribeToBlockHeaderEvents() (uuid.UUID, <-chan *rpc.BlockHeaderLogEntry, error) {
	if len(blockEventSource.blockMonitors) == 0 {
		slog.Warn("no block event providers are running, block event provider will not work until at least one block provider is running")
	}

	return blockEventSource.Subscribe()
}

func UnsubscribeFromBlockHeaderEvents(id uuid.UUID) {
	blockEventSource.Unsubscribe(id)
}
