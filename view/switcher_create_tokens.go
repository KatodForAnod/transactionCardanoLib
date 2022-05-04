package view

import (
	"fmt"
	"log"
	"transactionCardanoLib/cardanocli"
)

func (f *Frontend) switcherCreateTokens(command int) error {
	switch command {
	case buildTransaction:
		cliOut, errOutput, err := f.createTokens.CardanoQueryUtxo()
		if err != nil {
			log.Println(err)
			for _, s := range errOutput {
				fmt.Println(s)
			}
			return err
		}
		fmt.Println(cliOut)

		var processParams cardanocli.TransactionParams
		processParams.Fee = "300000"
		processParams.Output = "0"

		fmt.Println("input txHash")
		fmt.Scan(&processParams.TxHash)
		fmt.Println("input txIx")
		fmt.Scan(&processParams.Txix)
		fmt.Println("input amount")
		fmt.Scan(&processParams.Funds)

		f.createTokens.SetProcessParams(processParams)

		errOutput, err = f.createTokens.TransactionBuild(f.conf.Token)
		if err != nil {
			for _, s := range errOutput {
				fmt.Println(s)
			}

			log.Println(err)
			return err
		}

		errOutput, err = f.createTokens.CalculateFee()
		if err != nil {
			log.Println(err)
			for _, s := range errOutput {
				fmt.Println(s)
			}
			return err
		}

		err = f.createTokens.CalculateOutPut()
		if err != nil {
			log.Println(err)
			return err
		}

		errOutput, err = f.createTokens.TransactionBuild(f.conf.Token)
		if err != nil {
			log.Println(err)
			for _, s := range errOutput {
				fmt.Println(s)
			}
			return err
		}
	case signTransaction:
		errOutput, err := f.createTokens.TransactionSign()
		if err != nil {
			log.Println(err)
			for _, s := range errOutput {
				fmt.Println(s)
			}
			return err
		}
	case submitTransaction:
		errOutput, err := f.createTokens.TransactionSubmit()
		if err != nil {
			log.Println(err)
			for _, s := range errOutput {
				fmt.Println(s)
			}
			return err
		}
	case showCardanoUtxo:
		cliOut, errOutput, err := f.createTokens.CardanoQueryUtxo()
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
