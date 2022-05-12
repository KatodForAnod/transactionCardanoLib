package view

import (
	"errors"
	"fmt"
	"log"
	"transactionCardanoLib/cardanocli"
)

func (f *Frontend) switcherCreateNft(command int) error {
	if !f.conf.UsingExistingPolicy {
		err := errors.New("does not work without policy")
		log.Println(err)
		return err
	}

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

		processParams, _, err := cardanocli.Parse(cliOut)
		if err != nil {
			log.Println(err)
			return err
		}

		fmt.Println("input output")
		fmt.Scan(&processParams.Output)
		fmt.Println("input slotnumber")
		fmt.Scan(&processParams.SlotNumber)

		f.sendTokens.SetProcessParams(processParams)

		errOutput, err = f.createTokens.TransactionBuild(f.conf.Token)
		if err != nil {
			for _, s := range errOutput {
				fmt.Println(s)
			}
			log.Println(err)
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
