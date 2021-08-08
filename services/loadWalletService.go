package services

import (
	"errors"
	"fmt"

	"github.com/alvinwong12/cold-wallet/models/coin"
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

func chooseCoinType() (coin.CoinType, error){
	fmt.Println("Select a coin type from the menu:")
	for i, coinType := range coin.GetSupportedCoinTypes(){
		fmt.Printf("%s %d\n", coinType, i)
	}
	var choice coin.CoinType
	fmt.Scanln(&choice)
	if choice.CheckSupportCompatability() {
		return choice, nil
	} else {
		return coin.NOT_A_COIN, errors.New("Invalid choice")
	}
}

func getWalletFile() string {
	fmt.Println("Please enter the file path of your wallet:")
	var file string
	fmt.Scanln(&file)
	return file
}
