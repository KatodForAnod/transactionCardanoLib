package view

import (
	"fmt"
	"log"
	"transactionCardanoLib/cardanocli"
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

		processParams, tokens, err := cardanocli.Parse(cliOut)
		if err != nil {
			log.Println(err)
			return err
		}

		processParams.Fee = "0"
		processParams.Output = "0"

		fmt.Println("input receiver")
		fmt.Scan(&processParams.Receiver)
		fmt.Println("input receiverOutput")
		fmt.Scan(&processParams.ReceiverOutput)

		f.sendTokens.SetProcessParams(processParams)

		var sendToken config.Token
		fmt.Println("input name of token to send")
		fmt.Scan(&sendToken.TokenName)
		fmt.Println("input amount of token to send")
		fmt.Scan(&sendToken.TokenAmount)

		errOutput, err = f.sendTokens.TransactionBuild(tokens, sendToken)
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

		errOutput, err = f.sendTokens.TransactionBuild(tokens, sendToken)
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
