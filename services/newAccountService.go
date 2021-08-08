package services

import (
	"fmt"

	"github.com/alvinwong12/cold-wallet/models/coin"
	"github.com/alvinwong12/cold-wallet/models/wallet"
)
type NewAccountService struct {

}

func (service *NewAccountService) Run(serviceConfig *ServiceConfig) (interface{}, error) {
	loadWalletService :=  LoadWalletService{}
	loadedWallet, err := loadWalletService.Run(serviceConfig)
	if err != nil {
		return nil, err
	}

	switch serviceConfig.CoinType {
		case coin.ETHEUREM:
			return loadedWallet.(*wallet.ETHWallet).MakeNewAccount(), nil
		default:
			return nil, &coin.UnsupportedCoinError{Message: fmt.Sprintf("NewAccountService: %s is currently unsupported by this service", serviceConfig.CoinType.String())} 
	}
}
