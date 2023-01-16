package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"flare-node-proxy/flags"
	"flare-node-proxy/middlewares"
)

func main() {
	// Use an external setup function in order
	// to configure the app in tests as well
	app := Setup()

	// Listen on port 3000 by default or port defined via cli flags
	log.Fatal(app.Listen(*flags.Port)) // go run app.go -port=:3000
}

func Setup() *fiber.App {
	// Parse command-line flags
	flag.Parse()

	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork:     *flags.Prod, // go run app.go -prod
		AppName:     "Flare Node Proxy",
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// Standard middlewares
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} ${latency} - ${method} ${path}\n",
	}))

	// Enable monitor from Fiber
	if *flags.Enablemonitor {
		app.Get("/flareproxy/metrics", monitor.New(monitor.Config{Title: "Flare Node Proxy Metrics Page"}))
		fmt.Println("Monitor enabled")
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

	return app
}
