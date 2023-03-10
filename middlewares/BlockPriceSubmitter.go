package middlewares

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	block_tx_logger "flare-node-proxy/logging"
	"flare-node-proxy/utils"
	"flare-node-proxy/whitelist"
)

// Example Blocked Transaction request
/*
  {
	from: "0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8",
	to: "0x1000000000000000000000000000000000000003",
	gasLimit: "21000",
	maxFeePerGas: "300",
	maxPriorityFeePerGas: "10",
	nonce: "0",
	value: "10000000000"
  }
*/

type CChainPOSTRequestMethodOnly struct {
	Method string `json:"method" xml:"method" form:"method"`
}
type UnsignedTransactionRequestParams []struct {
	From string `json:"from" xml:"from" form:"from"`
	To   string `json:"to" xml:"to" form:"to"`
}
type UnsignedTransactionRequest struct {
	Params UnsignedTransactionRequestParams `json:"params" xml:"params" form:"params"`
}
type SignedTransactionRequest struct {
	Params []string `json:"params" xml:"params" form:"params"`
}

var PRICE_SUBMITTER_ADDRESS = "0x1000000000000000000000000000000000000003"
var BLOCKED_METHODS = []string{"eth_signTransaction", "eth_sendTransaction", "eth_sendRawTransaction"}

var blocked_response = fiber.Map{
	"jsonrpc": "2.0",
	"id":      2,
	"error": fiber.Map{
		"code":    -32600,
		"message": "transaction blocked",
	},
}

// Check if the body of the request contains field called 'to' with the PriceSubmitter contract addres
// 0x1000000000000000000000000000000000000003
func BlockPriceSubmitter(c *fiber.Ctx) error {
	req := new(CChainPOSTRequestMethodOnly)

	if err := c.BodyParser(req); err != nil {
		log.Errorf(err.Error())
		return c.Status(400).SendString("bad request")
	}

	if utils.StringInSlice(req.Method, BLOCKED_METHODS) { // if method is "eth_sendRawTransaction"
		if req.Method == BLOCKED_METHODS[2] {
			tx := new(SignedTransactionRequest)
			if err := c.BodyParser(tx); err != nil {
				log.Errorf(err.Error())
				return c.Status(400).SendString("bad request")
			}
			for _, params := range tx.Params {
				to, from, err := utils.GetFromToTransaction_EIP1559(params)
				if err != nil {
					log.Errorf(err.Error())
					return c.Status(400).SendString("bad request")
				}
				if to.String() == PRICE_SUBMITTER_ADDRESS && !whitelist.CheckWhitelist(from.String()) {
					log.Infof("Blocked tx to PriceSubmitter from address %s", from.String())
					block_tx_logger.BlockedIpsAndAddressesLogger.WithFields(log.Fields{"ip": c.IP(), "from": from.String()}).Info("transaction blocked")
					return c.Status(200).JSON(blocked_response)
				}
			}
		} else { // if method is "eth_sendTransaction" or "eth_signTransaction"
			tx := new(UnsignedTransactionRequest)
			if err := c.BodyParser(tx); err != nil {
				log.Errorf(err.Error())
				return c.Status(400).SendString("bad request")
			}
			for _, params := range tx.Params {
				if params.To == PRICE_SUBMITTER_ADDRESS && !whitelist.CheckWhitelist(params.From) {
					log.Infof("Blocked tx to PriceSubmitter from address %s", params.From)
					block_tx_logger.BlockedIpsAndAddressesLogger.WithFields(log.Fields{"ip": c.IP(), "from": params.From}).Info("transaction blocked")
					return c.Status(200).JSON(blocked_response)
				}
			}
		}

	}
	return c.Next()
}
