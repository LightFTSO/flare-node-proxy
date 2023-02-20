package utils

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gofiber/fiber/v2"
)

func RemoveHexPrefix(str string) string {
	if (str[:2]) == "0x" { // Remove hex prefix "0x..."
		return str[2:]
	}

	return str
}

func StringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func GetFromToTransaction_EIP1559(signed_tx string) (*common.Address, *common.Address, error) {
	signed_tx = RemoveHexPrefix(signed_tx) // Remove hex prefix "0x..."

	var tx = &types.Transaction{}
	raw_tx_data, err := hex.DecodeString(signed_tx)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, nil, err
	}
	if err := tx.UnmarshalBinary(raw_tx_data); err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, nil, err
	}

	signer := types.LatestSignerForChainID(tx.ChainId())
	sender, err := signer.Sender(tx)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, nil, err
	}

	return tx.To(), &sender, nil
}

func GetFromToTransaction(signed_tx string) (*common.Address, *common.Address, error) {
	signed_tx = RemoveHexPrefix(signed_tx) // Remove hex prefix "0x..."

	raw_tx_data, err := hex.DecodeString(signed_tx)
	if err != nil {
		fmt.Println(err.Error())
		return nil, nil, err
	}

	var tx types.Transaction
	err = rlp.DecodeBytes(raw_tx_data, &tx)
	if err != nil {
		fmt.Println(err.Error())
		return nil, nil, err
	}

	signer := types.NewEIP155Signer(tx.ChainId())
	sender, err := signer.Sender(&tx)
	if err != nil {
		fmt.Println(err.Error())
		return nil, nil, err
	}

	return tx.To(), &sender, nil
}

func ParseAddr(raw string, network string) (string, string) { //nolint:revive // Returns (port, host)
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
