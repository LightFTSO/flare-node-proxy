package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Example Transaction request
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

type Params []struct {
	From string `json:"from" xml:"from" form:"from"`
	To   string `json:"to" xml:"to" form:"to"`
}
type CChainPOSTRequest struct {
	Id     int    `json:"id" xml:"id" form:"id"`
	Method string `json:"method" xml:"method" form:"method"`
	Params Params `json:"params" xml:"params" form:"params"`
}

var PRICE_SUBMITTER_ADDRESS = "0x1000000000000000000000000000000000000003"

// Check if the body of the request contains field called 'to' with the PriceSubmitter contract addres
// 0x1000000000000000000000000000000000000003
func BlockPriceSubmitter(c *fiber.Ctx) error {
	tx := new(CChainPOSTRequest)

	if err := c.BodyParser(tx); err != nil {
		return err
	}

	if c.Locals("verbose") == true {
		fmt.Printf("from: %s to: %s\n", tx.Params[0].From, tx.Params[0].To)
	}

	for _, params := range tx.Params {
		if params.To == PRICE_SUBMITTER_ADDRESS {
			fmt.Printf("Blocked to PriceSubmitter from address %s", params.From)
			return c.Status(400).SendStatus(400)
		}
	}

	return c.Next()
}
