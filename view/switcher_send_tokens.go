package view

import (
	"errors"
	"fmt"
	"log"
	"transactionCardanoLib/config"
)

func (f *Frontend) switcherSendTokens(command int) error {
	switch command {
	case buildTransaction:
		cliOut, errOutput, err := f.sendTokens.CardanoQueryUtxo()
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
		} else if len(processParams) == 0 || len(tokens) == 0 {
			return errors.New("params not found")
		}

		processParams[0].Fee = "0"
		processParams[0].Output = "0"

		fmt.Println("input receiver")
		fmt.Scan(&processParams[0].Receiver)
		fmt.Println("input receiverOutput")
		fmt.Scan(&processParams[0].ReceiverOutput)

		f.sendTokens.SetProcessParams(processParams[0])

		var sendToken config.Token
		fmt.Println("input name of token to send")
		fmt.Scan(&sendToken.TokenName)
		fmt.Println("input amount of token to send")
		fmt.Scan(&sendToken.TokenAmount)

		errOutput, err = f.sendTokens.TransactionBuild(tokens[0], []config.Token{sendToken})
		if err != nil {
			for _, s := range errOutput {
				fmt.Println(s)
			}

			log.Println(err)
			return err
		}

		errOutput, err = f.sendTokens.CalculateFee()
		if err != nil {
			log.Println(err)
			for _, s := range errOutput {
				fmt.Println(s)
			}
			return err
		}

		err = f.sendTokens.CalculateOutPut()
		if err != nil {
			log.Println(err)
			return err
		}

		errOutput, err = f.sendTokens.TransactionBuild(tokens[0], []config.Token{sendToken})
		if err != nil {
			log.Println(err)
			for _, s := range errOutput {
				fmt.Println(s)
			}
			return err
		}
	case signTransaction:
		errOutput, err := f.sendTokens.TransactionSign()
		if err != nil {
			log.Println(err)
			for _, s := range errOutput {
				fmt.Println(s)
			}
			return err
		}
	case submitTransaction:
		errOutput, err := f.sendTokens.TransactionSubmit()
		if err != nil {
			log.Println(err)
			for _, s := range errOutput {
				fmt.Println(s)
			}
			return err
		}
	case showCardanoUtxo:
		cliOut, errOutput, err := f.sendTokens.CardanoQueryUtxo()
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
