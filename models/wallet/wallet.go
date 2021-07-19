package wallet

import (
	// "fmt"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"

	// "github.com/ethereum/go-ethereum/core/types"
	"github.com/alvinwong12/cold-wallet/utils"
	"github.com/ethereum/go-ethereum/accounts"

	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

const (
	PURPOSE = "44'"
)

type HDWallet = hdwallet.Wallet

type ColdWallet struct {
	*HDWallet
	OwnerName string `json:"OwnerName"`
	Purpose string `json:"Purpose"`
	CoinType string `json:"CoinType"`
	EncryptedMnemonic string `json:"EncryptedMnemonic"`
	Index uint64 `json:"Index"`
}


func NewColdWallet(mnemonic string, ownerName string, coinType string) *ColdWallet {
	hdWallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}
	coldWallet := ColdWallet{
		HDWallet: hdWallet,
		OwnerName: ownerName,
		Purpose: PURPOSE,
		CoinType: coinType,
		EncryptedMnemonic: mnemonic,
		Index: uint64(0),
	}
	return &coldWallet
}



// func(coldWallet *ColdWallet) SetOwnerName(ownerName string){
// 	coldWallet.ownerName = ownerName
// }

// func(coldWallet *ColdWallet) GetOwnerName() string{
// 	return coldWallet.ownerName
// }

func(coldWallet *ColdWallet) MakeDerivationPathFromIndex(index uint64) string {
	return "m/" + coldWallet.Purpose + "/" + coldWallet.CoinType + "/" + "0'/0/" + strconv.FormatUint(index, 10)
}

func(coldWallet *ColdWallet) GetDerivationPath() string {
	return coldWallet.MakeDerivationPathFromIndex(coldWallet.Index)
}

func(coldWallet *ColdWallet) MakeNewAccount() *accounts.Account {
	coldWallet.Index += 1
	account := coldWallet.GetAccount(coldWallet.Index);
	return account
}

func(coldWallet *ColdWallet) GetAccount(index uint64) *accounts.Account {
	derivationPath := coldWallet.MakeDerivationPathFromIndex(index)
	path := hdwallet.MustParseDerivationPath(derivationPath)
	account, err := coldWallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}
	return &account
}

func(coldWallet *ColdWallet) GetAllAccounts() []*accounts.Account {
	accountsChannel := make(chan *accounts.Account)
	var wg sync.WaitGroup

	for i := 0; i <= int(coldWallet.Index) ;i++ {
		wg.Add(1)
		go func(index int){
			account := coldWallet.GetAccount(uint64(index))
			accountsChannel <- account
			wg.Done()
		}(i)
	}
	go func(){
		wg.Wait()
		close(accountsChannel)
	}()

	allAccounts := make([]*accounts.Account, coldWallet.Index+1)
	for acc := range accountsChannel {
		allAccounts = append(allAccounts, acc)
	}
	return allAccounts
}

// func(coldWallet *ColdWallet) GetIndex() uint64 {
// 	return coldWallet.index
// }

// func(coldWallet *ColdWallet) SetIndex(index uint64) {
// 	coldWallet.index = index
// }

func(coldWallet *ColdWallet) ToJSON() string {
	coldWalletJson, err := json.MarshalIndent(coldWallet, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(coldWalletJson)
}

func LoadWalletFromFile(file string) *ColdWallet {
	jsonData := utils.ImportFromFile(file)
	var coldWallet ColdWallet
	err := json.Unmarshal(jsonData, &coldWallet)
	if err != nil {
		log.Fatal(err)
	}
	return &coldWallet
}

func ExportWalletToFile(coldWallet *ColdWallet, file string) {
	utils.ExportToFile(coldWallet.ToJSON(), file)
	fmt.Printf("Wallet saved!\n")
}
