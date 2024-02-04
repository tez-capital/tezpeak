package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/tez-capital/tezpeak/configuration"
	"github.com/tez-capital/tezpeak/constants"
	"github.com/tez-capital/tezpeak/core"
	"github.com/tez-capital/tezpeak/core/common"
	"github.com/tez-capital/tezpeak/util"
)

type Message struct {
	Text string
}

func main() {

	logLevelFlag := flag.String("log-level", "info", "Log level")
	versionFlag := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Println(constants.TEZPEAK_VERSION)
		return
	}

	util.InitLog(*logLevelFlag)
	config, err := configuration.Load()
	if err != nil {
		panic(err)
	}

	statusChannel, err := core.Run(context.Background(), config)
	if err != nil {
		panic(err)
	}

	serializedStatus := ""
	var serializedStatusMutex sync.RWMutex // Use RWMutex for better read performance
	clientChannels := make(map[uuid.UUID]chan string)
	var clientChannelsMutex sync.Mutex // Add mutex for clientChannels map access

	go func() {
		for newStatus := range statusChannel {
			serializedStatusBytes, err := json.Marshal(core.PeakStatusUpdatedeRport{
				Kind: common.FullStatusUpdateKind,
				Data: newStatus,
			})
			if err != nil {
				slog.Warn("failed to serialize status", "error", err.Error())
				continue
			}

			serializedStatusMutex.Lock()
			serializedStatus = string(serializedStatusBytes)
			serializedStatusMutex.Unlock()

			// Notify all connected clients of the new status
			clientChannelsMutex.Lock()
			for _, ch := range clientChannels {
				select {
				case ch <- serializedStatus: // Non-blocking send
				default:
					slog.Warn("skipping client, channel full")
				}
			}
			clientChannelsMutex.Unlock()
		}
	}()

	app := fiber.New()

	app.Get("/sse", func(c *fiber.Ctx) error {
		// Set SSE headers
		// ... (no changes here)

		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			id, err := uuid.NewRandom()
			if err != nil {
				c.Status(500).SendString("Failed to generate UUID")
				return
			}

			serializedStatusMutex.RLock() // Use RLock for reading
			fmt.Fprintf(w, "data: %v\n\n", serializedStatus)
			serializedStatusMutex.RUnlock()

			statusUpdateChannel := make(chan string, 1) // Buffer to avoid blocking
			clientChannelsMutex.Lock()
			clientChannels[id] = statusUpdateChannel
			clientChannelsMutex.Unlock()
			w.Flush()

			defer func() {
				clientChannelsMutex.Lock()
				delete(clientChannels, id)
				clientChannelsMutex.Unlock()
				close(statusUpdateChannel)
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

	fmt.Println(app.Listen(config.Listen))
}
