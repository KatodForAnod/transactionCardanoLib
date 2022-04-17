package view

import (
	"fmt"
	"log"
	"transactionCardanoLib/cardanocli"
	"transactionCardanoLib/config"
)

type Frontend struct {
	conf       config.Config
	cardanoLib cardanocli.CardanoLib
}

const (
	buildTransaction  = 1
	signTransaction   = 2
	exitCommand       = 10
	showCardanoUtxo   = 3
	submitTransaction = 4
)

var (
	startMsg = fmt.Sprintf(
		"%d. Build transaction\n"+
			"%d. Sign transaction\n"+
			"%d. Show cardano utxo\n"+
			"%d. Submit transaction\n"+
			"%d. Exit\n",
		buildTransaction, signTransaction,
		showCardanoUtxo, submitTransaction, exitCommand)
)

func (f *Frontend) SetConfAndCardanoLib(conf config.Config,
	cardanoLib cardanocli.CardanoLib) {
	f.conf = conf
	f.cardanoLib = cardanoLib
}

func (f *Frontend) Start() error {
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

func (f *Frontend) switcher(command int) error {
	switch command {
	case buildTransaction:
		cliOut, errOutput, err := f.cardanoLib.CardanoQueryUtxo(f.conf.ID)
		if err != nil {
			log.Println(err)
			for _, s := range errOutput {
				fmt.Println(s)
			}
			return err
		}
		fmt.Println(cliOut)

		f.cardanoLib.TransactionParams.Fee = "300000"
		f.cardanoLib.TransactionParams.Output = "0"

		fmt.Println("input txHash")
		fmt.Scan(&f.cardanoLib.TransactionParams.TxHash)
		fmt.Println("input txIx")
		fmt.Scan(&f.cardanoLib.TransactionParams.Txix)
		fmt.Println("input amount")
		fmt.Scan(&f.cardanoLib.TransactionParams.Funds)
		fmt.Println("input token amount")
		fmt.Scan(&f.cardanoLib.TransactionParams.TokenAmount)

		var tokenNames []string
		for _, token := range f.conf.Token {
			tokenNames = append(tokenNames, token.TokenName)
		}
		f.cardanoLib.TransactionBuild(tokenNames)

		fee, errOutput, err := f.cardanoLib.CalculateFee(f.conf.ID)
		if err != nil {
			log.Println(err)
			for _, s := range errOutput {
				fmt.Println(s)
			}
			return err
		}
		f.cardanoLib.TransactionParams.Fee = fee

		output, err := f.cardanoLib.CalculateOutPut()
		if err != nil {
			log.Println(err)
			return err
		}
		f.cardanoLib.TransactionParams.Output = output

		errOutput, err = f.cardanoLib.TransactionBuild(tokenNames)
		if err != nil {
			log.Println(err)
			for _, s := range errOutput {
				fmt.Println(s)
			}
			return err
		}
	case signTransaction:
		errOutput, err := f.cardanoLib.TransactionSign(f.conf.ID)
		if err != nil {
			log.Println(err)
			for _, s := range errOutput {
				fmt.Println(s)
			}
			return err
		}
	case submitTransaction:
		errOutput, err := f.cardanoLib.TransactionSubmit(f.conf.ID)
		if err != nil {
			log.Println(err)
			for _, s := range errOutput {
				fmt.Println(s)
			}
			return err
		}
	default:
		fmt.Println("unsupported command")
	}

	return nil
}
