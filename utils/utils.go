package utils

import (
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gofiber/fiber/v2"
)

func StringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func GetFromToTransaction(signed_tx string) (*common.Address, *common.Address, error) {
	raw_tx_data, err := hex.DecodeString(signed_tx[2:]) // Remove hex prefix "0x..."
	if err != nil {
		return nil, nil, err
	}

	var tx types.Transaction
	err = rlp.DecodeBytes(raw_tx_data, &tx)
	if err != nil {
		return nil, nil, err
	}

	signer := types.NewEIP155Signer(tx.ChainId())
	sender, err := signer.Sender(&tx)
	if err != nil {
		return nil, nil, err
	}

	return tx.To(), &sender, nil
}

func ParseAddr(raw string, network string) (string, string) { //nolint:revive // Returns (host, port)
	host := ""
	port := ""

	if i := strings.LastIndex(raw, ":"); i != -1 {
		port = raw[i+1:]
		host = raw[:i]
	}

	if host == "" {
		if network == fiber.NetworkTCP6 {
			host = "[::1]"
		} else {
			host = "0.0.0.0"
		}
	}

	return port, host
}
