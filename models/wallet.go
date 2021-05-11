package models

import (
	"fmt"
	"log"
	// "github.com/ethereum/go-ethereum/core/types"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

type HDWallet = hdwallet.Wallet

type ColdWallet struct {
	*HDWallet
	ownerName string
	purpose string
	coinType string
	encryptedMnemonic string
}

type ETHWallet struct {
	*ColdWallet
}

func NewETHWallet(ownerName string, mnemonic string) *ETHWallet {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}
	ethWallet := ETHWallet{ColdWallet: &ColdWallet{HDWallet: wallet, ownerName: ownerName, purpose: "44'", coinType: "60'"}}
	return &ethWallet
}

func(wallet *ColdWallet) Hello() {
	fmt.Println("Hello")
	fmt.Println(wallet.ownerName)
}

func(wallet *ColdWallet) SetName(ownerName string){
	wallet.ownerName = ownerName
}

func(wallet *ColdWallet) GetName() string{
	return wallet.ownerName
}

