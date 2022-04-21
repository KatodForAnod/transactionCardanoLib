package main

import (
	"encoding/hex"
	"log"
	"transactionCardanoLib/cardanocli"
	"transactionCardanoLib/config"
	"transactionCardanoLib/view"
)

func main() {
	log.SetFlags(log.Lshortfile)

	cardanoLib := cardanocli.CardanoLib{
		TransactionParams: cardanocli.TransactionParams{},
	}

	front := view.Frontend{}
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	for i, token := range conf.Token {
		conf.Token[i].TokenName = hex.EncodeToString([]byte(token.TokenName))
	}

	/*conf := config.Config{
		ID:                         "1097911063",
		Token:                      []config.Token{
			{
				TokenName:   hex.EncodeToString([]byte("exampleToken")),
				TokenAmount: "10001",
			},
			{
				TokenName:   hex.EncodeToString([]byte("exampleToken2")),
				TokenAmount: "10001",
			},
		},
		PaymentAddress:             "",
		PaymentVKeyFilePath: cardanocli.PaymentVerifyKeyFile,
		PaymentSKeyFilePath: cardanocli.PaymentSignKeyFile,
		UsingExistingPolicy:        true,
		PolicyID:                   "",
		PolicyScriptFilePath:       cardanocli.PolicyScriptFile,
		PolicySigningFilePath:      cardanocli.PolicySigningKeyFile,
		PolicyVerificationFilePath: cardanocli.PolicyVerificationkeyFile,
	}*/
	front.SetConfAndCardanoLib(conf, cardanoLib)

	front.Start()
}
