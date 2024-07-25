package core

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tez-capital/tezpeak/configuration"
	"github.com/tez-capital/tezpeak/constants"
	"github.com/tez-capital/tezpeak/core/common"
	"github.com/tez-capital/tezpeak/core/providers/tezbake"
	"github.com/tez-capital/tezpeak/core/providers/tezpay"
)

type client struct {
	channel chan string
	closed  atomic.Bool
}

func (c *client) Close() {
	c.closed.Store(true)
	close(c.channel)
}

type clientStoreBase map[uuid.UUID]*client

type clientStore struct {
	clientStoreBase
	mtx sync.RWMutex
}

func (c *clientStore) Remove(id uuid.UUID) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	client := c.clientStoreBase[id]
	client.Close()
	delete(c.clientStoreBase, id)
}

func (c *clientStore) Each(f func(id uuid.UUID, client *client)) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	for id, client := range c.clientStoreBase {
		f(id, client)
	}
}

func (c *clientStore) Add(statusChannel chan string) (close func(), err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.clientStoreBase[id] = &client{
		channel: statusChannel,
	}
	return func() { c.Remove(id) }, nil
}

func newClientStore() *clientStore {
	return &clientStore{
		clientStoreBase: make(clientStoreBase),
	}
}

var (
	status = PeakStatus{
		Modules: make(map[string]any),
		Nodes:   make(map[string]common.NodeStatus),
	}
	statusMtx = sync.RWMutex{}
	clients   = newClientStore()
)

func createModuleStatusChannel(id string, statusChannel chan<- common.ModuleStatusUpdatedReport) chan<- common.StatusUpdatedReport {
	moduleStatusChannel := make(chan common.StatusUpdatedReport, 100)
	go func() {
		for statusUpdate := range moduleStatusChannel {
			statusChannel <- common.ModuleStatusUpdatedReport{
				Id:     id,
				Report: statusUpdate,
			}
		}
	}()

	return moduleStatusChannel
}

func registerStatusEndpoint(app *fiber.App) {
	app.Get("/sse", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			statusUpdateChannel := make(chan string, 100) // Buffer to avoid blocking
			unregisterClient, err := clients.Add(statusUpdateChannel)
			if err != nil {
				c.Status(500).SendString("Failed to generate UUID")
				return
			}
			statusMtx.RLock()
			fmt.Fprintf(w, "data: %v\n\n", status.MarshalJSON())
			statusMtx.RUnlock()

			w.Flush()

			defer func() {
				unregisterClient()
			}()

			for msg := range statusUpdateChannel {
				if _, err := fmt.Fprintf(w, "data: %v\n\n", msg); err != nil {
					// Handle client disconnection or error in sending message
					slog.Warn("error sending message to client", "error", err.Error())
					return
				}
				w.Flush()
			}

		})

		return nil
	})

}

func updateStatus(moduleId string, data any) {
	statusMtx.Lock()
	defer statusMtx.Unlock()
	status.Modules[moduleId] = data
}

func updateNodeStatus(nodeId string, data common.NodeStatus) {
	statusMtx.Lock()
	defer statusMtx.Unlock()
	status.Nodes[nodeId] = data
}

func notifyClients() {
	statusMtx.RLock()
	serializedStatus := status.MarshalJSON()
	statusMtx.RUnlock()

	clients.Each(func(_ uuid.UUID, c *client) {
		go func() {
			if c.closed.Load() {
				return
			}
			c.channel <- serializedStatus
		}()
	})
}

// TODO: optimize - diffing, module updates, etc
func runStatusUpdatesProcessing(moduleStatusChannel <-chan common.ModuleStatusUpdatedReport) {
	pendingUpdatesChannel := make(chan struct{}, 1)
	defer close(pendingUpdatesChannel)
	pendingUpdatesCounter := 0
	for {
		select {
		case statusUpdate, ok := <-moduleStatusChannel:
			if !ok {
				return
			}
			if pendingUpdatesCounter > 10 {
				notifyClients()
				pendingUpdatesCounter = 0
			}

			id := statusUpdate.Id
			switch data := statusUpdate.Report.GetData().(type) {
			case common.NodeStatus:
				updateNodeStatus(statusUpdate.Report.GetId(), data)
			default:
				updateStatus(id, data)
			}
			pendingUpdatesCounter++
			// try insert into pendingUpdatesChannel
			select {
			case pendingUpdatesChannel <- struct{}{}:
			default:
			}
		case <-pendingUpdatesChannel:
			if pendingUpdatesCounter > 0 {
				notifyClients()
				pendingUpdatesCounter = 0
			}
		}
	}
}

func Run(ctx context.Context, config *configuration.Runtime, app *fiber.App) error {
	status.Id = config.Id
	registerStatusEndpoint(app)

	moduleStatusChannel := make(chan common.ModuleStatusUpdatedReport, 100)
	go runStatusUpdatesProcessing(moduleStatusChannel)

	common.StartNodeStatusProviders(ctx, config.Nodes, createModuleStatusChannel("global", moduleStatusChannel))
	// modules
	for id := range config.Modules {
		switch id {
		case constants.TEZBAKE_MODULE_ID:
			ok, configuration := config.GetTezbakeModuleConfiguration()
			if !ok {
				slog.Warn("tezbake module configured but not loaded")
				continue
			}

			err := tezbake.SetupModule(ctx, configuration, app, createModuleStatusChannel(id, moduleStatusChannel))
			if err != nil {
				return err
			}
		case constants.TEZPAY_MODULE_ID:
			ok, configuration := config.GetTezpayModuleConfiguration()
			if !ok {
				slog.Warn("tezpay module configured but not loaded")
				continue
			}

			err := tezpay.SetupModule(ctx, configuration, app, createModuleStatusChannel(id, moduleStatusChannel))
			if err != nil {
				return err
			}

		}
	}

	return nil

}
