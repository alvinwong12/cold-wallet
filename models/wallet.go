package models

import (
	// "fmt"
	"log"
	"math"
	"math/big"
	"context"
	"strconv"

	// "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

const (
	PURPOSE = "44'"
	ETHEUREM = "60'"
)

type HDWallet = hdwallet.Wallet

type ColdWallet struct {
	*HDWallet
	ownerName string
	purpose string
	coinType string
	encryptedMnemonic string
	index uint64
}

type ETHWallet struct {
	*ColdWallet
	wei *big.Float
	gas_limit_simple_tx uint64
}

func NewColdWallet(mnemonic string, ownerName string, coinType string) *ColdWallet {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}
	coldWallet := ColdWallet{
		HDWallet: wallet,
		ownerName: ownerName,
		purpose: PURPOSE,
		coinType: coinType,
		index: uint64(0),
	}
	return &coldWallet
}

func NewETHWallet(mnemonic string, ownerName string) *ETHWallet {
	ethWallet := ETHWallet{
		ColdWallet: NewColdWallet(mnemonic, ownerName, ETHEUREM),
		wei: big.NewFloat(math.Pow10(18)),
		gas_limit_simple_tx: uint64(21000),
	}
	return &ethWallet
}

func(wallet *ColdWallet) SetOwnerName(ownerName string){
	wallet.ownerName = ownerName
}

func(wallet *ColdWallet) GetOwnerName() string{
	return wallet.ownerName
}

func(wallet *ColdWallet) GetDerivationPath() string {
	return "m/" + wallet.purpose + "/" + wallet.coinType + "/" + "0'/0/" + strconv.FormatUint(wallet.index, 10)
}

func(wallet *ColdWallet) GetAccount() *accounts.Account {
	path := hdwallet.MustParseDerivationPath(wallet.GetDerivationPath())
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}
	wallet.index += 1
	return &account
}

func(wallet *ETHWallet) EthToWei(ethAmount *big.Float) *big.Int {
	amount, _ := new(big.Float).Mul(ethAmount, wallet.wei).Int(new(big.Int))
	return amount
}

func(wallet *ETHWallet) WeiToEth(weiAmount *big.Int) *big.Float {
	weiInBigFloat := new(big.Float)
	weiInBigFloat.SetString(weiAmount.String())
	amount := new(big.Float).Quo(weiInBigFloat, wallet.wei)
	return amount
}

func(wallet *ETHWallet) GetBalanceInWei(client *ethclient.Client, account *accounts.Account) *big.Int{
	weiBalance, err := client.BalanceAt(context.Background(), account.Address, nil)

	if err != nil {
		log.Fatal(err)
	}

	return weiBalance
}

func(wallet *ETHWallet) GetBalanceInEth(client *ethclient.Client, account *accounts.Account) *big.Float{
	return wallet.WeiToEth(wallet.GetBalanceInWei(client, account))
}

func(wallet *ETHWallet) GetPendingBalanceInWei(client *ethclient.Client, account *accounts.Account) *big.Int {
	pendingBalance, err := client.PendingBalanceAt(context.Background(), account.Address)
	if err != nil {
		log.Fatal(err)
	}
	return pendingBalance
}

func(wallet *ETHWallet) GetPendingBalanceInEth(client *ethclient.Client, account *accounts.Account) *big.Float{
	return wallet.WeiToEth(wallet.GetPendingBalanceInWei(client, account))
}
