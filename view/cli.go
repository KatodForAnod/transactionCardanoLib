package view

import (
	"fmt"
	"log"
	"transactionCardanoLib/cardanocli"
	"transactionCardanoLib/config"
)

type Frontend struct {
	conf         config.Config
	createTokens cardanocli.CreateTokens
	sendTokens   cardanocli.SendTokens
}

const (
	buildTransaction    = 1
	signTransaction     = 2
	exitCommand         = 10
	showCardanoUtxo     = 3
	submitTransaction   = 4
	generatePolicyFiles = 5
	createTokens        = 6
	sendTokens          = 7
)

var (
	startMsg = fmt.Sprintf("%d. Create tokens\n"+
		"%d. Send tokens\n"+
		"%d. Exit\n",
		createTokens, sendTokens, exitCommand)
	transactionOpMsg = fmt.Sprintf(
		"%d. Build transaction\n"+
			"%d. Sign transaction\n"+
			"%d. Show cardano utxo\n"+
			"%d. Submit transaction\n"+
			"%d. Exit\n",
		buildTransaction, signTransaction,
		showCardanoUtxo, submitTransaction,
		exitCommand)
)

func (f *Frontend) SetConfAndCardanoLib(conf config.Config,
	createTokens cardanocli.CreateTokens,
	sendTokens cardanocli.SendTokens) {
	f.conf = conf
	f.createTokens = createTokens
	f.sendTokens = sendTokens
}

func (f *Frontend) Start() error {
	for {
		fmt.Print(startMsg)
		var choiceCommand int
		if _, err := fmt.Scan(&choiceCommand); err != nil {
			log.Println(err)
			return err
		}

		if choiceCommand == exitCommand {
			return nil
		}

		for {
			fmt.Print(transactionOpMsg)
			var choiceCommandTransaction int
			if _, err := fmt.Scan(&choiceCommandTransaction); err != nil {
				log.Println(err)
				return err
			}
			if choiceCommandTransaction == exitCommand {
				break
			}
			switch choiceCommand {
			case createTokens:
				if err := f.switcherCreateTokens(choiceCommandTransaction); err != nil {
					log.Println(err)
					return err
				}
			case sendTokens:
				if err := f.switcherSendTokens(choiceCommandTransaction); err != nil {
					log.Println(err)
					return err
				}
			}
		}
	}
}
