package main

import (
	"io/ioutil"
	"log"
	"transactionCardanoLib/cardanocli"
	"transactionCardanoLib/config"
	"transactionCardanoLib/view"
)

func main() {
	log.SetFlags(log.Lshortfile)
	conf, err := config.LoadConfig()
	if err != nil {
		panic(1)
	}

	token := conf.Token
	if !token.UsingExistingPolicy {
		_, _, paymentAddrFileName, err := cardanocli.GeneratePaymentAddr(conf.Token.ID)
		if err != nil {
			log.Println(err)
			panic(2)
		}

		fileContent, err := ioutil.ReadFile(paymentAddrFileName)
		if err != nil {
			log.Println(err)
			panic(3)
		}
		conf.Token.PaymentAddress = string(fileContent)

		conf.Token.PolicyVerificationFilePath,
			conf.Token.PolicySigningFilePath,
			conf.Token.PolicyScriptFilePath,
			err = cardanocli.GeneratePolicy()
		if err != nil {
			log.Println(err)
			panic(4)
		}

		policyIDFilePath, err := cardanocli.GeneratePolicyID()
		if err != nil {
			log.Println(err)
			panic(5)
		}

		fileContent, err = ioutil.ReadFile(policyIDFilePath)
		if err != nil {
			log.Println(err)
			panic(6)
		}
		conf.Token.PolicyID = string(fileContent)
	}

	front := view.Frontend{}
	front.Start(conf)
}
