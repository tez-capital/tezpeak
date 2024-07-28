package core

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"sync"

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
	closed  bool
	mtx     sync.RWMutex
}

func (c *client) Close() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.closed = true
	close(c.channel)
}

func (c *client) Send(msg string) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	if c.closed {
		return
	}
	c.channel <- msg
}

func newClient(statusChannel chan string) *client {
	return &client{
		channel: statusChannel,
	}
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

	c.clientStoreBase[id] = newClient(statusChannel)
	return func() { c.Remove(id) }, nil
}

func newClientStore() *clientStore {
	return &clientStore{
		clientStoreBase: make(clientStoreBase),
	}
}

var (
	status  = newPeakStatus()
	clients = newClientStore()
)

func createModuleStatusChannel(id string, statusChannel chan<- common.ModuleStatusUpdate) chan<- common.StatusUpdate {
	moduleStatusChannel := make(chan common.StatusUpdate, 100)
	go func() {
		for statusUpdate := range moduleStatusChannel {
			statusChannel <- common.NewModuleStatusUpdate(id, statusUpdate)
		}
	}()

	return moduleStatusChannel
}

func registerStatusEndpoint(app *fiber.Group) {
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

			fmt.Fprintf(w, "data: %v\n\n", status.String())
			w.Flush()

			defer unregisterClient()

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

func notifyClients() {
	serializedStatus := status.String()

	clients.Each(func(_ uuid.UUID, c *client) {
		go c.Send(serializedStatus)
	})
}

// TODO: optimize - diffing, module updates, etc
func runStatusUpdatesProcessing(statusChannel <-chan common.ModuleStatusUpdate) {
	pendingUpdatesChannel := make(chan struct{}, 1)
	defer close(pendingUpdatesChannel)

	for {
		select {
		case statusUpdate, ok := <-statusChannel:
			if !ok {
				return
			}

			module := statusUpdate.GetModule()
			switch statusUpdate := statusUpdate.GetStatusUpdate().(type) {
			case *common.NodeStatusUpdate:
				status.UpdateNodeStatus(statusUpdate.Id, statusUpdate.Status)
			default:
				status.UpdateModuleStatus(module, statusUpdate.GetData())
			}
			// try insert into pendingUpdatesChannel
			select {
			case pendingUpdatesChannel <- struct{}{}:
			default:
			}
		case <-pendingUpdatesChannel:
			notifyClients()
		}
	}
}

func Run(ctx context.Context, config *configuration.Runtime, app *fiber.Group) error {
	status.SetId(config.Id)
	registerStatusEndpoint(app)

	statusChannel := make(chan common.ModuleStatusUpdate, 100)
	go runStatusUpdatesProcessing(statusChannel)

	common.StartNodeStatusProviders(ctx, config.Nodes, createModuleStatusChannel("global", statusChannel))
	// modules
	for id := range config.Modules {
		switch id {
		case constants.TEZBAKE_MODULE_ID:
			ok, configuration := config.GetTezbakeModuleConfiguration()
			if !ok {
				slog.Warn("tezbake module configured but not loaded")
				continue
			}

			err := tezbake.SetupModule(ctx, configuration, app, createModuleStatusChannel(id, statusChannel))
			if err != nil {
				return err
			}
		case constants.TEZPAY_MODULE_ID:
			ok, configuration := config.GetTezpayModuleConfiguration()
			if !ok {
				slog.Warn("tezpay module configured but not loaded")
				continue
			}

			err := tezpay.SetupModule(ctx, configuration, app, createModuleStatusChannel(id, statusChannel))
			if err != nil {
				return err
			}

		}
	}

	return nil

}
