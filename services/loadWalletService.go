package services

import (
	"github.com/alvinwong12/cold-wallet/models/wallet"
)

type LoadWalletService struct {

}

func (service *LoadWalletService) Run(serviceConfig *ServiceConfig) (interface{}, error){
	return wallet.LoadWalletFromFile(serviceConfig.WalletFilePath, serviceConfig.CoinType, serviceConfig.password, serviceConfig.encrypted), nil
}
