package services

import (

	"fmt"
	"github.com/alvinwong12/cold-wallet/models/coin"
	"github.com/alvinwong12/cold-wallet/models/wallet"
)

type CheckBalanceService struct {

}



func (service *CheckBalanceService) Run(serviceConfig *ServiceConfig) (interface{}, error){
	loadWalletService :=  LoadWalletService{}
	loadedWallet, err := loadWalletService.Run(serviceConfig)
	if err != nil {
		return nil, err
	}

	switch serviceConfig.CoinType {
		case coin.ETHEUREM:
			return ethBalanceService(loadedWallet.(*wallet.ETHWallet), serviceConfig)
		default:
			return nil, &coin.UnsupportedCoinError{Message: fmt.Sprintf("CheckBalanceService: %s is currently unsupported by this service", serviceConfig.CoinType.String())} 
	}
}

func ethBalanceService(ethWallet *wallet.ETHWallet, serviceConfig *ServiceConfig) (interface{}, error) {
	if !checkEtheureumNetworkStatus(serviceConfig.ETHNetworkClient) {
		return nil, &ServiceUnavailableError{Message: "CheckBalanceService: Cannot connect to the Etheurem network."}
	}
	return ethWallet.GetBalanceInEth(serviceConfig.ETHNetworkClient), nil
}
