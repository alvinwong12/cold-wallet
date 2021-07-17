package models

import (
	// "fmt"
	"log"
	"math"
	"math/big"
	"context"
	"strconv"
	"sync"

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

func(wallet *ColdWallet) makeDerivationPathFromIndex(index uint64) string {
	return "m/" + wallet.purpose + "/" + wallet.coinType + "/" + "0'/0/" + strconv.FormatUint(index, 10)
}

func(wallet *ColdWallet) GetDerivationPath() string {
	return wallet.makeDerivationPathFromIndex(wallet.index)
}

func(wallet *ColdWallet) GetNewAccount() *accounts.Account {
	account := wallet.GetAccount(wallet.GetDerivationPath());
	wallet.index += 1
	return account
}

func(wallet *ColdWallet) GetAccount(derivationPath string) *accounts.Account {
	path := hdwallet.MustParseDerivationPath(derivationPath)
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}
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

func(wallet *ETHWallet) GetBalanceInWei(client *ethclient.Client) *big.Int{
	balanceChannel := make(chan big.Int)
	var wg sync.WaitGroup

	for i := 0; i < int(wallet.index) ;i++ {
		wg.Add(1)
		go func(index int) {
			account := wallet.GetAccount(wallet.makeDerivationPathFromIndex(uint64(index)))
			weiBalance, err := client.BalanceAt(context.Background(), account.Address, nil)
			if err != nil {
				balanceChannel <- *big.NewInt(int64(0))
			} else {
				balanceChannel <- *weiBalance
			}
			wg.Done()
		}(i)
	}

	go func(){
		wg.Wait()
		close(balanceChannel)
	}()

	totalWeiBalance := big.NewInt(int64(0))
	for balance := range balanceChannel {
		totalWeiBalance.Add(totalWeiBalance, &balance)
	}
	return totalWeiBalance
}

func(wallet *ETHWallet) GetBalanceInEth(client *ethclient.Client) *big.Float{
	return wallet.WeiToEth(wallet.GetBalanceInWei(client))
}

func(wallet *ETHWallet) GetPendingBalanceInWei(client *ethclient.Client) *big.Int {

	pendingBalanceChannel := make(chan big.Int)
	var wg sync.WaitGroup

	for i := 0; i < int(wallet.index) ;i++ {
		wg.Add(1)
		go func(index int) {
			account := wallet.GetAccount(wallet.makeDerivationPathFromIndex(uint64(index)))
			pendingBalance, err := client.PendingBalanceAt(context.Background(), account.Address)
			if err != nil {
				pendingBalanceChannel <- *big.NewInt(int64(0))
			} else {
				pendingBalanceChannel <- *pendingBalance
			}
			wg.Done()
		}(i)
	}

	go func(){
		wg.Wait()
		close(pendingBalanceChannel)
	}()

	totalPendingWeiBalance := big.NewInt(int64(0))
	for balance := range pendingBalanceChannel {
		totalPendingWeiBalance.Add(totalPendingWeiBalance, &balance)
	}
	return totalPendingWeiBalance
}

func(wallet *ETHWallet) GetPendingBalanceInEth(client *ethclient.Client) *big.Float{
	return wallet.WeiToEth(wallet.GetPendingBalanceInWei(client))
}
