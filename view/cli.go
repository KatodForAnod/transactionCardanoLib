package view

import (
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"transactionCardanoLib/cardanocli"
	"transactionCardanoLib/config"
)

type Frontend struct {
	conf config.Config
}

const (
	buildTransaction = 1
	signTransaction  = 2
	exitCommand      = 3
)

var (
	startMsg = fmt.Sprintf(
		"%d. Build transaction\n"+
			"%d. Sign transaction\n"+
			"%d. Exit\n",
		buildTransaction, signTransaction, exitCommand)
)

func (f Frontend) Start(conf config.Config) error {
	f.conf = conf
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
		if err := f.buildTransaction(); err != nil {
			log.Println(err)
			return err
		}
	case signTransaction:
		fmt.Println("input id")
		var id string
		var obj config.TokenStruct

		err := cardanocli.TransactionSign(id, obj)
		if err != nil {
			log.Println(err)
			return err
		}
	default:
		fmt.Println("unsupported command")
	}

	return nil
}

func (f Frontend) buildTransaction() error {
	var fee, txHash, txIx, output, tokenName1, tokenName2, tokenAmount string

	fmt.Println("write fee")
	if _, err := fmt.Scan(&fee); err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("write txHash")
	if _, err := fmt.Scan(&txHash); err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("write txIx")
	if _, err := fmt.Scan(&txIx); err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("write output")
	if _, err := fmt.Scan(&output); err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("write tokenName1")
	if _, err := fmt.Scan(&tokenName1); err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("write tokenName2")
	if _, err := fmt.Scan(&tokenName2); err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("write tokenAmount")
	if _, err := fmt.Scan(&tokenAmount); err != nil {
		log.Println(err)
		return err
	}

	tokenName1 = hex.EncodeToString([]byte(tokenName1))
	tokenName2 = hex.EncodeToString([]byte(tokenName2))

	err := cardanocli.TransactionBuild(fee, txHash, txIx, f.conf.Token.PaymentAddress, output,
		strconv.FormatInt(f.conf.Token.TokenAmount, 10), // tokenAmount ???
		tokenName1, tokenName2, f.conf.Token.PolicyID, f.conf.Token.PolicySigningFilePath)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
