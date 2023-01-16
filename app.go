package main

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/lightftso/flare-proxy/middlewares"

	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	port          = flag.String("port", ":3000", "Port to listen on, can also include a full host to enable listening on multiple interfaces e.g. 0.0.0.0:3000")
	endpoint      = flag.String("endpoint", "http://localhost:9650", "Flare (or Avax) node this program will proxy requests to")
	enablemonitor = flag.Bool("monitor", false, "Enable Fiber Server Monitor on /flareproxy/metrics")
	prod          = flag.Bool("prod", false, "Enable prefork in Production")
)

func main() {
	// Parse command-line flags
	flag.Parse()

	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork:      *prod, // go run app.go -prod
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		IdleTimeout:  time.Second * 5,
		AppName:      "Flare Node Proxy",
	})

	proxy.WithTlsConfig(&tls.Config{
		InsecureSkipVerify: true,
	})

	// Standard middlewares
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} ${latency} - ${method} ${path}\n",
	}))

	// Enable monitor from Fiber
	if *enablemonitor {
		//app.Get("/flareproxy/metrics", monitor.New(monitor.Config{Title: "Flare Node Proxy Metrics Page"}))
		fmt.Println("Monitor enabled")
	}

	// Block transactions to PriceSubmitter
	app.Use("*", middlewares.BlockPriceSubmitter)

	// Proxy requests to specified host
	app.All("*", func(c *fiber.Ctx) error {
		url := *endpoint + c.Path()
		fmt.Println(url)
		if err := proxy.Do(c, url); err != nil {
			return err
		}
		return nil
	})

	// Listen on port 3000 by default or port defined via cli flags
	log.Fatal(app.Listen(*port)) // go run app.go -port=:3000
}
