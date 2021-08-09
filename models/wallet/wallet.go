package wallet

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/alvinwong12/cold-wallet/models/coin"
	"github.com/alvinwong12/cold-wallet/utils"
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
	CoinType coin.CoinType `json:"CoinType"`
	Mnemonic string `json:"Mnemonic"`
	Index uint64 `json:"Index"`
}

func makeNewHDWallet(mnemonic string) *HDWallet {
	hdWallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}
	return hdWallet
}

func NewColdWallet(mnemonic string, ownerName string, coinType coin.CoinType) (*ColdWallet, error) {
	if !coinType.CheckSupportCompatability() {
		return  nil, &coin.UnsupportedCoinError{Message: coinType.String() + " is not supported!"}
	}
	coldWallet := ColdWallet{
		HDWallet: makeNewHDWallet(mnemonic),
		OwnerName: ownerName,
		Purpose: PURPOSE,
		CoinType: coinType,
		Mnemonic: mnemonic,
		Index: uint64(0),
	}
	return &coldWallet, nil
}

func(coldWallet *ColdWallet) WithIndex(index uint64) *ColdWallet {
	coldWallet.Index = index
	return coldWallet
}

func(coldWallet *ColdWallet) MakeDerivationPathFromIndex(index uint64) string {
	return "m/" + coldWallet.Purpose + "/" + coldWallet.CoinType.Repr() + "/" + "0'/0/" + strconv.FormatUint(index, 10)
}

func(coldWallet *ColdWallet) MakeNewAccount() *accounts.Account {
	coldWallet.Index += 1
	account := coldWallet.GetAccount(coldWallet.Index);
	return account
}

func(coldWallet *ColdWallet) GetAccount(index uint64) *accounts.Account {
	derivationPath := coldWallet.MakeDerivationPathFromIndex(index)
	path := hdwallet.MustParseDerivationPath(derivationPath)
	account, err := coldWallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}
	return &account
}

func(coldWallet *ColdWallet) GetAllAccounts() []*accounts.Account {
	accountsChannel := make(chan *accounts.Account)
	var wg sync.WaitGroup

	for i := 0; i <= int(coldWallet.Index) ;i++ {
		wg.Add(1)
		go func(index int){
			account := coldWallet.GetAccount(uint64(index))
			accountsChannel <- account
			wg.Done()
		}(i)
	}
	go func(){
		wg.Wait()
		close(accountsChannel)
	}()

	allAccounts := make([]*accounts.Account, coldWallet.Index+1)
	var cur int=0
	for acc := range accountsChannel {
		allAccounts[cur] = acc
		cur++
	}
	return allAccounts
}

func(coldWallet *ColdWallet) ToJSON() string {
	coldWalletJson, err := json.MarshalIndent(coldWallet, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(coldWalletJson)
}

func(coldWallet ColdWallet) ExportWalletToFile(file string) {
	utils.ExportToFile(coldWallet.ToJSON(), file)
	fmt.Printf("Wallet saved!\n")
}

func(coldWallet ColdWallet) ExportWalletToFileEncrypted(file string, key string) {
	encryptedWallet := ColdWallet{
		OwnerName: coldWallet.OwnerName,
		Purpose: coldWallet.Purpose,
		CoinType: coldWallet.CoinType,
		Mnemonic: coldWallet.Mnemonic,
		Index: coldWallet.Index,
	}
	encryptedWallet.Mnemonic = hex.EncodeToString(utils.Encrypt([]byte(encryptedWallet.Mnemonic), key))
	encryptedWallet.ExportWalletToFile(file)
}

func LoadWalletFromFile(file string, coinType coin.CoinType, password string, encrypted bool) interface{} {
	switch coinType {
		case coin.ETHEUREM:
			return loadETHWalletFromFile(file, password, encrypted)
		default:
			return loadColdWalletFromFile(file, password, encrypted)
	}
}

func loadColdWalletFromFile(file string, password string, encrypted bool) *ColdWallet {
	coldWallet := ColdWallet{}
	jsonData := utils.ImportFromFile(file)
	err := json.Unmarshal(jsonData, &coldWallet)
	if err != nil {
		log.Fatal(err)
	}
	if encrypted {
		mnemonicBytes, err := hex.DecodeString(coldWallet.Mnemonic)
		if err != nil {
			log.Fatal(err)
		}
		coldWallet.Mnemonic = string(utils.Decrypt(mnemonicBytes, password))
	}
	coldWallet.HDWallet = makeNewHDWallet(coldWallet.Mnemonic)
	return &coldWallet
}

func(coldWallet *ColdWallet) GetLatestAccount() *accounts.Account {
	return coldWallet.GetAccount(coldWallet.Index)
}
