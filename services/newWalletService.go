package services

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alvinwong12/cold-wallet/models/coin"
	"github.com/alvinwong12/cold-wallet/models/wallet"
)

type NewWalletService struct {

}

func (service *NewWalletService) Run(serviceConfig *ServiceConfig) (interface{}, error){
	ownerName := getOwnerName()
	mnemonic := getMnemonic()
	switch serviceConfig.CoinType {
		case coin.ETHEUREM:
			newWallet := wallet.NewETHWallet(mnemonic, ownerName)
			newWallet.ExportWalletToFileEncrypted(serviceConfig.WalletFilePath, serviceConfig.password)
			fmt.Printf("Wallet: %s\n", newWallet.ToJSON())
			return newWallet, nil
		default:
			newWallet, err := wallet.NewColdWallet(mnemonic , ownerName , serviceConfig.CoinType)
			newWallet.ExportWalletToFileEncrypted(serviceConfig.WalletFilePath, serviceConfig.password)
			if err != nil {
				return nil, err
			}
			fmt.Printf("Wallet: %s\n", newWallet.ToJSON())
			return newWallet, nil
	}
}

func getOwnerName() string {
	fmt.Println("Please enter the owner name of this wallet:")
	var name string
	fmt.Scanln(&name)
	return name
}

func getMnemonic() string{
	fmt.Println("Please enter the mnemonic used for the new wallet (eg. \"tag volcano eight thank tide danger coast health above argue embrace heavy\") :")
	var mnemonic string
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	mnemonic = scanner.Text()
	return strings.TrimSpace(mnemonic)
}
