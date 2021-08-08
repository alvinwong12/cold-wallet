package services

import (
	"errors"
	"fmt"

	"github.com/alvinwong12/cold-wallet/models/coinType"
	"github.com/alvinwong12/cold-wallet/models/wallet"
)

type LoadWalletService struct {

}

func (service *LoadWalletService) Run() (interface{}, error){
	chosenCoin, err := chooseCoinType()
	if err != nil {
		return nil, err
	}
	file := getWalletFile()
	return wallet.LoadWalletFromFile(file, chosenCoin), nil
}

func chooseCoinType() (coinType.CoinType, error){
	fmt.Println("Select a coin type from the menu:")
	for i, coin := range coinType.GetSupportedCoinTypes(){
		fmt.Printf("%s %d\n", coin, i)
	}
	var choice coinType.CoinType
	fmt.Scanln(&choice)
	if choice.CheckSupportCompatability() {
		return choice, nil
	} else {
		return coinType.NOT_A_COIN, errors.New("Invalid choice")
	}
}

func getWalletFile() string {
	fmt.Println("Please enter the file path of your wallet:")
	var file string
	fmt.Scanln(&file)
	return file
}
