package services

import (
	"errors"
	"fmt"

	"github.com/alvinwong12/cold-wallet/models/wallet"
	"github.com/alvinwong12/cold-wallet/utils"
)

type LoadWalletService struct {

}

func (service *LoadWalletService) Run(serviceConfig *ServiceConfig) (interface{}, error){
	if !utils.FileExists(serviceConfig.WalletFilePath) {
		return nil, errors.New(fmt.Sprintf("Wallet file -> %s ,cannot be read.\n", serviceConfig.WalletFilePath))
	}
	return wallet.LoadWalletFromFile(serviceConfig.WalletFilePath, serviceConfig.CoinType, serviceConfig.password, serviceConfig.encrypted), nil
}
