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
	buildTransaction    = 1
	signTransaction     = 2
	exitCommand         = 10
	showCardanoUtxo     = 3
	submitTransaction   = 4
	generatePolicyFiles = 5
)

var (
	startMsg = fmt.Sprintf(
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
	cardanoLib cardanocli.CardanoLib) {
	f.conf = conf
	f.cardanoLib = cardanoLib

	f.cardanoLib.TransactionParams.ID = conf.ID                      // default init var
	f.cardanoLib.TransactionParams.PaymentAddr = conf.PaymentAddress // default init var
	cardanocli.PaymentSignKeyFile = f.conf.PaymentSKeyFilePath       // default init var
	cardanocli.PaymentVerifyKeyFile = f.conf.PaymentVKeyFilePath     // default init var

	if f.conf.UsingExistingPolicy {
		f.cardanoLib.UseExistPolicy(conf)
	} else {
		f.cardanoLib.GeneratePolicyFiles()
	}
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
		cliOut, errOutput, err := f.cardanoLib.CardanoQueryUtxo()
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

		errOutput, err = f.cardanoLib.TransactionBuild(f.conf.Token)
		if err != nil {
			for _, s := range errOutput {
				fmt.Println(s)
			}

			log.Println(err)
			return err
		}

		fee, errOutput, err := f.cardanoLib.CalculateFee()
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

		errOutput, err = f.cardanoLib.TransactionBuild(f.conf.Token)
		if err != nil {
			log.Println(err)
			for _, s := range errOutput {
				fmt.Println(s)
			}
			return err
		}
	case signTransaction:
		errOutput, err := f.cardanoLib.TransactionSign()
		if err != nil {
			log.Println(err)
			for _, s := range errOutput {
				fmt.Println(s)
			}
			return err
		}
	case submitTransaction:
		errOutput, err := f.cardanoLib.TransactionSubmit()
		if err != nil {
			log.Println(err)
			for _, s := range errOutput {
				fmt.Println(s)
			}
			return err
		}
	case showCardanoUtxo:
		cliOut, errOutput, err := f.cardanoLib.CardanoQueryUtxo()
		if err != nil {
			log.Println(err)
			for _, s := range errOutput {
				fmt.Println(s)
			}
			return err
		}
		fmt.Println(cliOut)
	default:
		fmt.Println("unsupported command")
	}

	return nil
}
