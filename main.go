package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"

	"github.com/tez-capital/tezpeak/configuration"
	"github.com/tez-capital/tezpeak/constants"
	"github.com/tez-capital/tezpeak/core"
	"github.com/tez-capital/tezpeak/util"
)

type Message struct {
	Text string
}

//go:embed web/dist/*
var staticFiles embed.FS

type staticFs struct {
	http.FileSystem
}

func (c staticFs) Open(name string) (http.File, error) {
	// try to open .html file first, because filesystem middleware forbids to open folders
	// but we want to provide html if it exists when user tries to access path like /tezpay
	f, err := c.FileSystem.Open(name + ".html")
	if err != nil {
		return c.FileSystem.Open(name)
	}
	return f, nil
}

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

	app := fiber.New()
	group, ok := app.Group("/api").(*fiber.Group)
	if !ok {
		panic("failed to create api group")
	}

	err = core.Run(context.Background(), config, group)
	if err != nil {
		panic(err)
	}

	app.Use("/", filesystem.New(filesystem.Config{
		Root:         staticFs{http.FS(staticFiles)},
		Index:        "index.html",
		NotFoundFile: "/web/dist/index.html",
		PathPrefix:   "/web/dist",
		Browse:       false,
	}))

	err = app.Listen(config.Listen)
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
