package view

import (
	"errors"
	"fmt"
	"log"
	"transactionCardanoLib/config"
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

		processParams, tokens, err := f.sendTokens.ParseUtxo(cliOut)
		if err != nil {
			log.Println(err)
			return err
		} else if len(processParams) == 0 {
			return errors.New("params not found")
		}

		processParams[0].Fee = "300000"
		processParams[0].Output = "0"
		f.createTokens.SetProcessParams(processParams[0])

		var tokenTmp []config.Token
		if len(tokens) > 0 {
			tokenTmp = tokens[0]
		}

		errOutput, err = f.createTokens.TransactionBuild(tokenTmp, f.conf.Token)
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

		errOutput, err = f.createTokens.TransactionBuild(tokens[0], f.conf.Token)
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
