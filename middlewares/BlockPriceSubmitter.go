package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"flare-node-proxy/utils"
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
		"message": "transaction blocked"},
}

// Check if the body of the request contains field called 'to' with the PriceSubmitter contract addres
// 0x1000000000000000000000000000000000000003
func BlockPriceSubmitter(c *fiber.Ctx) error {
	req := new(CChainPOSTRequestMethodOnly)

	if err := c.BodyParser(req); err != nil {
		return err
	}

	if utils.StringInSlice(req.Method, BLOCKED_METHODS) {
		if req.Method == BLOCKED_METHODS[2] {
			tx := new(SignedTransactionRequest)
			if err := c.BodyParser(tx); err != nil {
				return err
			}
			for _, params := range tx.Params {
				to, from, err := utils.GetFromToTransaction(params)
				if err != nil {
					return err
				}
				if to.String() == PRICE_SUBMITTER_ADDRESS {
					fmt.Printf("Blocked tx to PriceSubmitter from address %s\n", from.String())
					return c.Status(200).JSON(blocked_response)
				}
			}
		} else {
			tx := new(UnsignedTransactionRequest)
			if err := c.BodyParser(tx); err != nil {
				return err
			}
			for _, params := range tx.Params {
				if params.To == PRICE_SUBMITTER_ADDRESS {
					fmt.Printf("Blocked tx to PriceSubmitter from address %s\n", params.From)
					return c.Status(200).JSON(blocked_response)
				}
			}
		}

	}
	return c.Next()
}
