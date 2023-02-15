package whitelist

import (
	"os"
	"strings"

	"flare-node-proxy/utils"

	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
)

type WhitelistedAddresses []string

var whitelisted_addresses WhitelistedAddresses = make(WhitelistedAddresses, 0)

func readWhitelistFile(path *string) {
	whitelist_file, err := os.ReadFile(*path)
	if err != nil {
		log.Errorf("Error reading whitelist file from %s", *path)
		return
	}

	if err = json.Unmarshal(whitelist_file, &whitelisted_addresses); err != nil {
		log.Errorf("Error parsing white list file from %s", *path)
		return
	}

	log.Infoln("Whitelisted addresses:", whitelisted_addresses)
	for i, address := range whitelisted_addresses {
		whitelisted_addresses[i] = strings.ToLower(address)
	}
}

/*
*
Reads the whitelist file on init and then listens for changes in it, updating the whitelist
*
*/
func InitPriceSubmitterWhitelist(path *string) {
	readWhitelistFile(path)
	go watchFile(*path)
}

/*
*
Checks if an address is on the whiltelist
*
*/
func CheckWhitelist(address string) bool {
	return utils.StringInSlice(strings.ToLower(address), whitelisted_addresses)
}
