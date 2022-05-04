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
	creatPolicy  cardanocli.Policy
	sendTokens   cardanocli.SendTokens
}

const (
	buildTransaction  = 1
	signTransaction   = 2
	exitCommand       = 10
	showCardanoUtxo   = 3
	submitTransaction = 4
	createNft         = 9
	createTokens      = 6
	sendTokens        = 7
	createPolicy      = 8
)

var (
	startMsg = fmt.Sprintf("%d. Create tokens\n"+
		"%d. Send tokens\n"+
		"%d. Create policy\n"+
		"%d. Create nft\n"+
		"%d. Exit\n",
		createTokens, sendTokens, createPolicy,
		createNft, exitCommand)

	transactionOpMsg = fmt.Sprintf("%d. Build transaction\n"+
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
	sendTokens cardanocli.SendTokens,
	createPolicy cardanocli.Policy) {
	f.conf = conf
	f.createTokens = createTokens
	f.sendTokens = sendTokens
	f.creatPolicy = createPolicy
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
			case createPolicy:
				if err := f.creatPolicy.GeneratePolicyFiles(); err != nil {
					log.Println(err)
					return err
				}
			case createNft:
				if err := f.switcherCreateNft(choiceCommandTransaction); err != nil {
					log.Println(err)
					return err
				}
			}
		}
	}
}
