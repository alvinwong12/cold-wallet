package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alvinwong12/cold-wallet/services"
)

func main(){

	if (len(os.Args) < 2){
		usage()
		os.Exit(1)
	}
	command := os.Args[1]

	var service services.Service
	switch strings.ToLower(command) {
		case "new-wallet":
			service = &services.NewWalletService{}
		case "new-account":
			service = &services.NewAccountService{}
		case "get-account":
			service = &services.GetAccountService{}
		case "get-balance":
			service = &services.CheckBalanceService{}
		default:
			usage()
			os.Exit(1)
	}
	serviceConfig, err := services.Init()
	if err != nil {
		log.Fatal(err)
	}
	_, err = service.Run(serviceConfig)
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Println("Commands available:")
	fmt.Println("")
	fmt.Println("	new-wallet			Create a new Wallet.")
	fmt.Println("	new-account			Create a new account.")
	fmt.Println("	get-account			Get an account address.")
	fmt.Println("	get-balance			Check wallet balance. (Caution! Requires internet access)")
	fmt.Println("")
}
