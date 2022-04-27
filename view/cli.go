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
	sendTokens   cardanocli.SendTokens
}

const (
	buildTransaction    = 1
	signTransaction     = 2
	exitCommand         = 10
	showCardanoUtxo     = 3
	submitTransaction   = 4
	generatePolicyFiles = 5
	createTokens        = 6
	sendTokens          = 7
)

var (
	startMsg = fmt.Sprintf("%d. Create tokens\n"+
		"%d. Send tokens\n"+
		"%d. Exit\n",
		createTokens, sendTokens, exitCommand)
	transactionOpMsg = fmt.Sprintf(
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
	createTokens cardanocli.CreateTokens,
	sendTokens cardanocli.SendTokens) {
	f.conf = conf
	f.createTokens = createTokens
	f.sendTokens = sendTokens
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
			}
		}
	}
}

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

		var processParams cardanocli.TransactionParams
		processParams.Fee = "0"
		processParams.Output = "0"

		fmt.Println("input txHash")
		fmt.Scan(&processParams.TxHash)
		fmt.Println("input txIx")
		fmt.Scan(&processParams.Txix)
		fmt.Println("input amount")
		fmt.Scan(&processParams.Funds)
		fmt.Println("input receiver")
		fmt.Scan(&processParams.Receiver)
		fmt.Println("input receiverOutput")
		fmt.Scan(&processParams.ReceiverOutput)

		f.sendTokens.SetProcessParams(processParams)

		var amount int
		fmt.Println("how many tokens do u have?")
		fmt.Scan(&amount)

		var tokens []config.Token
		for i := 0; i < amount; i++ {
			var token config.Token
			fmt.Println("input name of token")
			fmt.Scan(&token.TokenName)
			fmt.Println("input amount of token")
			fmt.Scan(&token.TokenAmount)
			tokens = append(tokens, token)
		}

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
