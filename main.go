package main

import (
	"io/ioutil"
	"transactionCardanoLib/config"
	"transactionCardanoLib/policy"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic(1)
	}

	token := conf.Token
	if !token.UsingExistingPolicy {
		_, _, paymentAddrFileName, err := policy.GeneratePaymentAddr(conf.Token.ID)
		if err != nil {
			panic(2)
		}

		fileContent, err := ioutil.ReadFile(paymentAddrFileName)
		if err != nil {
			panic(3)
		}

		token.PaymentAddress = string(fileContent)
	}

}
