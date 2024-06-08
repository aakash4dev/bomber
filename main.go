package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	for {
		if err := sendTransactions(); err != nil {
			fmt.Printf("Error occurred: %s, retrying in 10 seconds...\n", err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		break
	}
}

func sendTransactions() error {
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		return err
	}

	privateKeyHex := "0x2580A14CF3F5DFD7F57FFE14C615F8078AA50930E354DDFA9E45577F6104B620"
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		return err
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return err
	}
	//
	recipientAddress := common.HexToAddress("0xBdb56Cf303763cBAC0F610D5ec6B811Aa5d91693")
	value := big.NewInt(1)
	gasLimit := uint64(210000) // A typical gas limit for a simple transfer
	data := []byte("lets bomb the network with transactions! AMF to the moon : ) 🚀")

	for i := 0; i < 10000; i++ {
		fmt.Println(i)

		nonce, err := client.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(privateKey.PublicKey))
		if err != nil {
			return err
		}

		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			return err
		}

		tx := types.NewTransaction(nonce, recipientAddress, value, gasLimit, gasPrice, data)
		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			return err
		}

		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			return err
		}

		fmt.Printf("Transaction hash: %s\n", signedTx.Hash().Hex())
		time.Sleep(1 * time.Second)
	}

	return nil
}
