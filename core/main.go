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
	ctx     context.Context
	mtx     sync.Mutex
}

func (c *client) Send(msg string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.closed {
		return
	}

	select {
	case c.channel <- msg:
	case <-c.ctx.Done():
		c.closed = true
		close(c.channel)
	}
}

func (c *client) Close() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.closed = true
	close(c.channel)
}

func newClient(ctx context.Context, statusChannel chan string) *client {
	return &client{
		channel: statusChannel,
		ctx:     ctx,
	}
}

type clientStore struct {
	m sync.Map
}

func (c *clientStore) Remove(id uuid.UUID) {
	entry, ok := c.m.Load(id)
	if !ok {
		return
	}
	client := entry.(*client)
	defer client.Close()
	c.m.Delete(id)
}

func (c *clientStore) Each(f func(id uuid.UUID, client *client)) {
	c.m.Range(func(key, value any) bool {
		id := key.(uuid.UUID)
		client := value.(*client)
		f(id, client)
		return true
	})
}

func (c *clientStore) Add(ctx context.Context, statusChannel chan string) (close func(), err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	c.m.Store(id, newClient(ctx, statusChannel))
	return func() { c.Remove(id) }, nil
}

func newClientStore() *clientStore {
	return &clientStore{
		m: sync.Map{},
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

		context := c.Context()

		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			statusUpdateChannel := make(chan string, 100) // Buffer to avoid blocking
			unregisterClient, err := clients.Add(context, statusUpdateChannel)
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
					slog.Debug("error sending message to client", "error", err.Error())
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
