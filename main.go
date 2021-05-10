package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/core/types"
)

func main(){
	PRIV_KEY_1 := "42ad2eac07b8e80beeee9f5bb08191425cecf97c493aa5f35525ca90e2b5b923"
	ACCOUNT_1 := "0x0c8e7f549FE5b0c297Ac7535059E2f280C016787"
	// PRIV_KEY_2 := "f580d19ee593a51eb62f5077c86c49d03129feed1c6a99326a61f4e2b97f573c"
	ACCOUNT_2 := "0xf8BA069f4455e4C77E4e0B13d5Ccc616F85F5Ab4"
	WEI_FLOAT := big.NewFloat(math.Pow10(18))
	// WEI_INT := big.NewInt(1000000000000000000)
	GAS_LIMIT := uint64(21000)
	client, err := ethclient.Dial("http://127.0.0.1:7545");
	if err != nil {
		log.Fatal(err);
	}

	// Block

	header, err := client.HeaderByNumber(context.Background(), nil)
	// header, err := client.HeaderByNumber(context.Background(), big.NewInt(0))
	if err != nil {
		log.Fatal(err)
	}

	blockNumber := big.NewInt(header.Number.Int64());
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	for _,tx := range block.Transactions() {
		fmt.Println(tx)
	}

	// Account
	account := common.HexToAddress(ACCOUNT_1)
	weiBalance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(weiBalance)
	ethBalance := new(big.Float)
	ethBalance.SetString(weiBalance.String())
	ethBalance = new(big.Float).Quo(ethBalance, WEI_FLOAT)
	fmt.Println(ethBalance)

	// Send ETH
	privKey1, err := crypto.HexToECDSA(PRIV_KEY_1)
	if err != nil {
		log.Fatal(err)
	}

	publKey1ECDSA, ok := privKey1.Public().(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Error casting public key to ESDSA")
	}

	// fmt.Println(publKey1ECDSA)
	fromAddr := crypto.PubkeyToAddress(*publKey1ECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		log.Fatal(err)
	}

	amount, _ := new(big.Float).Mul(big.NewFloat(float64(2)), WEI_FLOAT).Int(new(big.Int))
	// amount := new(big.Int).Mul(big.NewInt(2), WEI_INT)
	fmt.Printf("Amount sending %s\n", amount.String())
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("Gas Price %s\n", gasPrice.String())
	toAddr := common.HexToAddress(ACCOUNT_2)
	fmt.Printf("To address %s\n", toAddr.Hex())
	tx := types.NewTransaction(nonce, toAddr, amount, GAS_LIMIT, gasPrice, nil)

	chainID, err := client.NetworkID(context.Background())
	if err != nil{
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privKey1)

	client.SendTransaction(context.Background(), signedTx)

	fmt.Printf("signedtx sent: %s\n", signedTx.Hash().Hex())
}
