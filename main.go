package main

import (
	"encoding/hex"
	"io/ioutil"
	"log"
	"transactionCardanoLib/cardanocli"
	"transactionCardanoLib/config"
	"transactionCardanoLib/files"
	"transactionCardanoLib/view"
)

func main() {
	log.SetFlags(log.Lshortfile)

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	for i, token := range conf.Token {
		conf.Token[i].TokenName = hex.EncodeToString([]byte(token.TokenName))
	}

	f := files.Files{}
	f.Init(conf)

	var policy cardanocli.Policy
	if !conf.UsingExistingPolicy {
		policy.Init(f, conf)
		err := policy.GeneratePolicyFiles()
		if err != nil {
			log.Println(err)
			return
		}
	}

	policyIdBytes, err := ioutil.ReadFile(f.GetPolicyIDFile())
	if err != nil {
		log.Println(err)
		return
	}

	baseParams := cardanocli.BaseTransactionParams{
		PaymentAddr: conf.PaymentAddress,
		PolicyID:    string(policyIdBytes),
		ID:          conf.ID,
	}

	createTokens := cardanocli.CreateTokens{}
	createTokens.Init(baseParams, cardanocli.TransactionParams{}, f)

	sendTokens := cardanocli.SendTokens{}
	sendTokens.Init(baseParams, cardanocli.TransactionParams{}, f)

	front := view.Frontend{}
	front.SetConfAndCardanoLib(conf, createTokens, sendTokens, policy)

	front.Start()
}
