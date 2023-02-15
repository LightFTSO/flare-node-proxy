package main

import (
	"flag"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"flare-node-proxy/flags"
	"flare-node-proxy/logging"
	"flare-node-proxy/middlewares"
	"flare-node-proxy/utils"
	"flare-node-proxy/whitelist"

	log "github.com/sirupsen/logrus"
)

func main() {
	logging.Setup()

	// Use an external setup function in order
	// to configure the app in tests as well
	app := Setup()

	// Listen on port 3000 by default or port defined via cli flags
	go func() {
		log.Fatal(app.Listen(*flags.Port)) // go run app.go -port=:3000
	}()

	whitelist.InitPriceSubmitterWhitelist(flags.WhitelistFilePath)
	<-make(chan struct{})
}

func Setup() *fiber.App {
	// Parse command-line flags
	flag.Parse()

	if !*flags.Prod {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork:               *flags.Prod, // go run app.go -prod
		AppName:               "Flare Node Proxy",
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
	})

	// Standard middlewares
	app.Use(recover.New())

	if !*flags.Prod {
		app.Use(logger.New(logger.Config{
			Format: "[${ip}]:${port} ${status} ${latency} - ${method} ${path}\n",
		}))
	}

	// Enable monitor from Fiber
	if *flags.Enablemonitor {
		app.Get("/flareproxy/metrics", monitor.New(monitor.Config{Title: "Flare Node Proxy Metrics Page"}))
		log.Info("Monitor enabled")
	}

	// Block transactions to PriceSubmitter
	app.Post("*", middlewares.BlockPriceSubmitter)

	// Proxy requests to specified host
	app.All("*", func(c *fiber.Ctx) error {
		url := *flags.Endpoint + c.Path()
		if err := proxy.Do(c, url); err != nil {
			return err
		}
		return nil
	})

	port, host := utils.ParseAddr(*flags.Port, app.Config().Network)
	log.Infof("Server listenning on %s:%s", host, port)
	return app
}
