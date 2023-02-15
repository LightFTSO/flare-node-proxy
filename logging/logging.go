package logging

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var BlockedIpsAndAddressesLogger = log.New()

func Setup() {
	err := os.MkdirAll("./logs", 0750)
	if err != nil {
		panic(err)
	}

	logfile, err := os.OpenFile("./logs/output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	multioutput := io.MultiWriter(os.Stderr, logfile)
	log.SetFormatter(&prefixed.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	})
	log.SetOutput(multioutput)

	blocked_txs_logfile, err := os.OpenFile("./logs/blocked.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	BlockedIpsAndAddressesLogger.SetFormatter(&log.JSONFormatter{})
	BlockedIpsAndAddressesLogger.SetOutput(blocked_txs_logfile)
}
