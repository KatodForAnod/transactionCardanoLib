package main

import (
	"io/ioutil"
	"log"
	"transactionCardanoLib/config"
	"transactionCardanoLib/policy"
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
		_, _, paymentAddrFileName, err := policy.GeneratePaymentAddr(conf.Token.ID)
		if err != nil {
			log.Println(err)
			panic(2)
		}

		fileContent, err := ioutil.ReadFile(paymentAddrFileName)
		if err != nil {
			log.Println(err)
			panic(3)
		}
		token.PaymentAddress = string(fileContent)

		token.PolicyVerificationFilePath,
			token.PolicySigningFilePath,
			token.PolicyScriptFilePath,
			err = policy.GeneratePolicy()
		if err != nil {
			log.Println(err)
			panic(4)
		}
	}

	front := view.Frontend{}
	front.Start(conf)
}
