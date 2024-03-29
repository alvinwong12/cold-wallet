package services

import (
	"context"
	"fmt"

	"github.com/alvinwong12/cold-wallet/models/coin"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	LOCAL_NET = "http://127.0.0.1:7545"
	MAIN_NET = "https://mainnet.infura.io"
)

type Service interface {
	Run(serviceConfig *ServiceConfig) (interface{}, error)
}

// config used to store all information services require
type ServiceConfig struct {
	CoinType coin.CoinType
	WalletFilePath string
	ETHNetworkClient *ethclient.Client
	password string
	encrypted bool
}

type ServiceUnavailableError struct {
	Message string
}

func (e *ServiceUnavailableError) Error() string {
	return e.Message
}

func Init() (*ServiceConfig, error) {
	file := getWalletFile()
	password := getPassword()
	chosenCoin := chooseCoinType()
	
	// serviceConfig
	serviceConfig := ServiceConfig{CoinType: chosenCoin, WalletFilePath: file}
	switch chosenCoin {
		case coin.ETHEUREM:
			client, err := ethclient.Dial(LOCAL_NET)
			if err != nil {
				return nil, err
			}
			serviceConfig.ETHNetworkClient = client
		default:
			break
	}

	serviceConfig.password = password
	serviceConfig.encrypted = true
	return &serviceConfig, nil
}

func chooseCoinType() coin.CoinType{
	fmt.Println("Select a coin type from the menu:")
	for i, coinType := range coin.GetSupportedCoinTypes(){
		fmt.Printf("%s %d\n", coinType, i)
	}
	var choice coin.CoinType
	fmt.Scanln(&choice)
	for choice.CheckSupportCompatability() == false {
		fmt.Println("Invalid choice. please choose again")
		fmt.Scanln(&choice)
	}
	return choice
}

func getWalletFile() string {
	fmt.Println("Please enter the file path of your wallet:")
	var file string
	fmt.Scanln(&file)
	return file
}

func checkEtheureumNetworkStatus(client *ethclient.Client) bool {
	if client == nil {
		return false
	}
	_, err := client.NetworkID(context.Background())
	if err != nil {
		return false
	}
	return true
}

func getPassword() string {
	fmt.Println("Please enter the password for encrypting your wallet (16 digits):")
	var password string
	fmt.Scanln(&password)
	for len(password) != 16 {
		fmt.Println("Please enter a 16 digit password:")
		fmt.Scanln(&password)
	}
	return password
}
