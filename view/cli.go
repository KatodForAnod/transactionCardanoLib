package view

import (
	"fmt"
	"log"
	"strconv"
	"transactionCardanoLib/config"
	"transactionCardanoLib/policy"
)

type Frontend struct {
	conf config.Config
}

const (
	buildTransaction = 1
	signTransaction  = 2
	exitCommand      = 3
	startMsg         = "1. Build transaction\n" +
		"2. Sign transaction\n" +
		"3. Exit\n"
)

func (f Frontend) Start(conf config.Config) error {
	fmt.Print(startMsg)

	for {
		var choiceCommand int
		if _, err := fmt.Scan(&choiceCommand); err != nil {
			log.Println(err)
			return err
		}

		if choiceCommand == exitCommand {
			return nil
		}

		if err := f.switcher(choiceCommand); err != nil {
			log.Println(err)
			return err
		}
	}
}

func (f Frontend) switcher(command int) error {
	switch command {
	case buildTransaction:
		fmt.Println("write fee, txHash, txIx, output, tokenAmount, tokenName1, tokenName2")
		var fee, thHash, txIx, output, tokenName1, tokenName2 string
		err := policy.TransactionBuild(fee, thHash, txIx, f.conf.Token.PaymentAddress, output,
			strconv.FormatInt(f.conf.Token.TokenAmount, 10),
			tokenName1, tokenName2, f.conf.Token.PolicySigningFilePath)

		if err != nil {
			log.Println(err)
			return err
		}
	case signTransaction:
		fmt.Println("input id")
		var id string
		var obj config.TokenStruct

		err := policy.TransactionSign(id, obj)
		if err != nil {
			log.Println(err)
			return err
		}
	default:
		fmt.Println("unsupported command")
	}

	return nil
}
