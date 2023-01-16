package flags

import "flag"

var (
	Port          = flag.String("port", ":3000", "Port to listen on, can also include a full host to enable listening on multiple interfaces e.g. 0.0.0.0:3000")
	Endpoint      = flag.String("endpoint", "https://songbird-api.flare.network", "Flare (or Avax) node this program will proxy requests to")
	Enablemonitor = flag.Bool("monitor", false, "Enable Fiber Server Monitor on /flareproxy/metrics")
	Prod          = flag.Bool("prod", false, "Enable prefork in Production")
)
