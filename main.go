package main

import (
	"fmt"
	"context"
	"crypto/ecdsa"

	"log"
	"math"
	"math/big"

	// "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/alvinwong12/cold-wallet/models"
	"github.com/ethereum/go-ethereum/ethclient"
	// hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	// hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

const (
	LOCAL_NET = "http://127.0.0.1:7545"
)

func main(){
	mnemonic := "tag volcano eight thank tide danger coast health above argue embrace heavy"
	ethWallet := models.NewETHWallet(mnemonic, "Alvin")

	client, err := ethclient.Dial(LOCAL_NET)
	if err != nil {
		log.Fatal(err)
	}

	account := ethWallet.GetAccount()

	// nonce, err := client.PendingNonceAt(context.Background(), account.Address)
	// fmt.Printf("nonce %d\n", nonce)

	weiBalance := ethWallet.GetBalanceInWei(client, account) 
	fmt.Printf("Wei %d\n", weiBalance)
	ethBalance :=  ethWallet.GetBalanceInEth(client, account)
	fmt.Printf("Eth %f\n", ethBalance)

	// ethWallet.Hello()
	PRIV_KEY_1 := "4f3db04ceab3caf460e6ec18b6ee4c3c2341136740058530736f8e35b9f7cc90"
	// ACCOUNT_1 := "0x40a5441e90087b1b2E3A099Da15b9299DEE190f3"
	// // PRIV_KEY_2 := "7cbf100453594f5bc5d1b1778b9e6d6b176759b241230e76d04bce3a88e1a9b6"
	// ACCOUNT_2 := "0xC0d6c55d2fC59C366e8f6020F0B3f2572654628d"
	WEI_FLOAT := big.NewFloat(math.Pow10(18))
	// // WEI_INT := big.NewInt(1000000000000000000)
	GAS_LIMIT := uint64(21000)
	// client, err := ethclient.Dial("http://127.0.0.1:7545");
	// if err != nil {
	// 	log.Fatal(err);
	// }
	
	// // Block

	// header, err := client.HeaderByNumber(context.Background(), nil)
	// // header, err := client.HeaderByNumber(context.Background(), big.NewInt(0))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// blockNumber := big.NewInt(header.Number.Int64());
	// block, err := client.BlockByNumber(context.Background(), blockNumber)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _,tx := range block.Transactions() {
	// 	fmt.Println(tx)
	// }

	// // Account
	// account := common.HexToAddress(ACCOUNT_1)
	// weiBalance, err := client.BalanceAt(context.Background(), account, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(weiBalance)
	// ethBalance := new(big.Float)
	// ethBalance.SetString(weiBalance.String())
	// ethBalance = new(big.Float).Quo(ethBalance, WEI_FLOAT)
	// fmt.Println(ethBalance)

	// // Send ETH
	privKey1, err := crypto.HexToECDSA(PRIV_KEY_1)
	if err != nil {
		log.Fatal(err)
	}

	publKey1ECDSA, ok := privKey1.Public().(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Error casting public key to ESDSA")
	}

	// // fmt.Println(publKey1ECDSA)
	fromAddr := crypto.PubkeyToAddress(*publKey1ECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		log.Fatal(err)
	}

	amount, _ := new(big.Float).Mul(big.NewFloat(float64(2)), WEI_FLOAT).Int(new(big.Int))
	// // amount := new(big.Int).Mul(big.NewInt(2), WEI_INT)
	fmt.Printf("Amount sending %s\n", amount.String())
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("Gas Price %s\n", gasPrice.String())
	// toAddr := common.HexToAddress(ACCOUNT_2)
	toAddr := account.Address
	fmt.Printf("To address %s\n", toAddr.Hex())
	tx := types.NewTransaction(nonce, toAddr, amount, GAS_LIMIT, gasPrice, nil)

	chainID, err := client.NetworkID(context.Background())
	if err != nil{
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privKey1)

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("signedtx sent: %s\n", signedTx.Hash().Hex())
	// tx, isPending, err := client.TransactionByHash(context.Background(), signedTx.Hash())

	// fmt.Println(isPending)
	
	// receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	// // receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash("0x13b664a3bfb931f20cd0dd3551a5261253defb69d45aa1cb3914eb1443fa8721"))
	// if err != nil {
	//   log.Fatal(err)
	// }

	// fmt.Println(tx.Hash())
	// fmt.Println(receipt.Status)
	// fmt.Println(receipt.Logs)
}
