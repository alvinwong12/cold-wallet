package ethWallet

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"math/big"
	"sync"

	"github.com/alvinwong12/cold-wallet/models/wallet"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	ETHEUREM = "60'"
)

type ETHWallet struct {
	*wallet.ColdWallet `json:"Wallet"`
	Wei *big.Float `json:"WEI"`
	Gas_limit_simple_tx uint64 `json:"GasLimitSimpleTx"`
}

func NewETHWallet(mnemonic string, ownerName string) *ETHWallet {
	ethWallet := ETHWallet{
		ColdWallet: wallet.NewColdWallet(mnemonic, ownerName, ETHEUREM),
		Wei: big.NewFloat(math.Pow10(18)),
		Gas_limit_simple_tx: uint64(21000),
	}
	return &ethWallet
}

func(etheuremWallet *ETHWallet) EthToWei(ethAmount *big.Float) *big.Int {
	amount, _ := new(big.Float).Mul(ethAmount, etheuremWallet.Wei).Int(new(big.Int))
	return amount
}

func(etheuremWallet *ETHWallet) WeiToEth(weiAmount *big.Int) *big.Float {
	weiInBigFloat := new(big.Float)
	weiInBigFloat.SetString(weiAmount.String())
	amount := new(big.Float).Quo(weiInBigFloat, etheuremWallet.Wei)
	return amount
}

func(etheuremWallet *ETHWallet) GetBalanceInWei(client *ethclient.Client) *big.Int{
	balanceChannel := make(chan big.Int)
	var wg sync.WaitGroup

	for i := 0; i < int(etheuremWallet.ColdWallet.Index) ;i++ {
		wg.Add(1)
		go func(index int) {
			account := etheuremWallet.GetAccount(etheuremWallet.ColdWallet.MakeDerivationPathFromIndex(uint64(index)))
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

func(etheuremWallet *ETHWallet) GetBalanceInEth(client *ethclient.Client) *big.Float{
	return etheuremWallet.WeiToEth(etheuremWallet.GetBalanceInWei(client))
}

func(etheuremWallet *ETHWallet) GetPendingBalanceInWei(client *ethclient.Client) *big.Int {

	pendingBalanceChannel := make(chan big.Int)
	var wg sync.WaitGroup

	for i := 0; i < int(etheuremWallet.ColdWallet.Index) ;i++ {
		wg.Add(1)
		go func(index int) {
			account := etheuremWallet.GetAccount(etheuremWallet.ColdWallet.MakeDerivationPathFromIndex(uint64(index)))
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

func(etheuremWallet *ETHWallet) GetPendingBalanceInEth(client *ethclient.Client) *big.Float{
	return etheuremWallet.WeiToEth(etheuremWallet.GetPendingBalanceInWei(client))
}

func(etheuremWallet *ETHWallet) ToJSON() string {
	etheuremWalletJson, err := json.MarshalIndent(etheuremWallet, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(etheuremWalletJson)
}
