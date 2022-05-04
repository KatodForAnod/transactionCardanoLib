package cardanocli

import (
	"errors"
	"strings"
	"transactionCardanoLib/config"
)

func Parse(cliOutput string) (TransactionParams, []config.Token, error) {
	cliOutputArr := strings.Split(cliOutput, "\n")
	if len(cliOutputArr) < 3 {
		return TransactionParams{}, nil, errors.New("error split")
	}

	str := cliOutputArr[2]
	vars := strings.Fields(str)

	txHash := vars[0]
	txIx := vars[1]
	amountLovelace := vars[2]

	var tokens []config.Token
	for i := 4; i < len(vars); i += 3 {
		policyAndTokenName := strings.Split(vars[i+2], ".")
		tokenName := policyAndTokenName[1]

		token := config.Token{
			TokenName:   tokenName,
			TokenAmount: vars[i+1],
		}

		tokens = append(tokens, token)
	}

	return TransactionParams{
			TxHash: txHash,
			Txix:   txIx,
			Funds:  amountLovelace,
		},
		tokens,
		nil
}
