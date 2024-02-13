package main

import (
	"bufio"
	"context"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/google/uuid"

	"github.com/tez-capital/tezpeak/configuration"
	"github.com/tez-capital/tezpeak/constants"
	"github.com/tez-capital/tezpeak/core"
	"github.com/tez-capital/tezpeak/core/common"
	"github.com/tez-capital/tezpeak/core/providers"
	"github.com/tez-capital/tezpeak/util"
)

type Message struct {
	Text string
}

//go:embed web/dist/*
var staticFiles embed.FS

func main() {
	logLevelFlag := flag.String("log-level", "info", "Log level")
	versionFlag := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("tezpeak %s - %s \n", constants.TEZPEAK_VERSION, constants.TEZPEAK_CODENAME)
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

	governanceProvider := providers.InitGovernanceProvider(context.Background(), config)
	if governanceProvider == nil {
		slog.Error("failed to get governance provider")
		os.Exit(1)
	}

	app.Get("/api/governance/can-vote", func(c *fiber.Ctx) error {
		return c.JSON(governanceProvider.CanVote())
	})

	app.Get("/api/governance/period-detail", func(c *fiber.Ctx) error {
		if !governanceProvider.CanVote() {
			return c.Status(403).SendString("not allowed")
		}

		periodInfo, err := governanceProvider.GetGovernancePeriodDetail(c.Context())
		if err != nil {
			return c.Status(500).SendString("failed to get governance period detail")
		}

		return c.JSON(periodInfo)
	})

	app.Get("/api/governance/available-pkhs", func(c *fiber.Ctx) error {
		if !governanceProvider.CanVote() {
			return c.Status(403).SendString("not allowed")
		}

		pkhs, err := governanceProvider.GetAvailablePkhs(c.Context())
		if err != nil {
			return c.Status(500).SendString("failed to get available pkhs")
		}
		return c.JSON(pkhs)
	})

	app.Post("/api/governance/vote", func(c *fiber.Ctx) error {
		if !governanceProvider.CanVote() {
			return c.Status(403).SendString("not allowed")
		}

		var params providers.VoteParams
		if err := c.BodyParser(&params); err != nil {
			return c.Status(400).SendString("invalid request")
		}

		opHash, err := governanceProvider.Vote(c.Context(), &params)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(opHash)
	})

	app.Post("/api/governance/upvote", func(c *fiber.Ctx) error {
		if !governanceProvider.CanVote() {
			return c.Status(403).SendString("not allowed")
		}

		var params providers.UpvoteParams
		if err := c.BodyParser(&params); err != nil {
			slog.Error("failed to parse upvote params", "error", err.Error())
			return c.Status(400).SendString("invalid request")
		}

		opHash, err := governanceProvider.Upvote(c.Context(), &params)
		if err != nil {
			slog.Error("failed to upvote", "error", err.Error())
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(opHash)
	})

	app.Post("/api/governance/wait-for-apply", func(c *fiber.Ctx) error {
		if !governanceProvider.CanVote() {
			return c.Status(403).SendString("not allowed")
		}

		var params string
		if err := c.BodyParser(&params); err != nil {
			slog.Error("failed to parse upvote params", "error", err.Error())
			return c.Status(400).SendString("invalid request")
		}

		applied, err := governanceProvider.WaitConfirmation(c.Context(), params)
		if err != nil {
			slog.Error("failed to upvote", "error", err.Error())
			return c.Status(500).SendString("failed to vote")
		}

		return c.JSON(applied)
	})

	app.Get("/sse", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

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

	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(staticFiles),
		Index:        "index.html",
		NotFoundFile: "/web/dist/index.html",
		PathPrefix:   "/web/dist",
		Browse:       false,
	}))

	fmt.Println(app.Listen(config.Listen))
}
