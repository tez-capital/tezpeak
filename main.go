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

	err = core.Run(context.Background(), config, app)
	if err != nil {
		panic(err)
	}

	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(staticFiles),
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
