package wallet

import (
	// "fmt"
	"encoding/json"
	"log"
	"strconv"

	// "github.com/ethereum/go-ethereum/core/types"
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

func(coldWallet *ColdWallet) GetNewAccount() *accounts.Account {
	account := coldWallet.GetAccount(coldWallet.GetDerivationPath());
	coldWallet.Index += 1
	return account
}

func(coldWallet *ColdWallet) GetAccount(derivationPath string) *accounts.Account {
	path := hdwallet.MustParseDerivationPath(derivationPath)
	account, err := coldWallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}
	return &account
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

