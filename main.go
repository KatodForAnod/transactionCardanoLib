package main

import (
	"transactionCardanoLib/config"
	"transactionCardanoLib/policy"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic(1)
	}

	if !conf.Token.UsingExistingPolicy {
		addr, err := policy.GeneratePaymentAddr(conf.Token.ID)
		if err != nil {
			panic(2)
		}

		conf.Token.PaymentAddress = addr
	}

}
