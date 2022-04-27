package cardanocli

import "transactionCardanoLib/config"

type TokenCreate interface {
	TransactionBuild(tokens []config.Token) (errorOutput []string, err error)
	CalculateFee() (fee string, errorOutput []string, err error)
	CalculateOutPut() (string, error)
	TransactionSign() (errorOutput []string, err error)
	TransactionSubmit() (errorOutput []string, err error)
}

type TokenSend interface {
	TransactionBuildSendingToken(tokens []config.Token,
		sendToken config.Token) (errorOutput []string, err error)
	TransactionSignSendingToken() (errorOutput []string, err error)
	TransactionSendTokenSubmit() (errorOutput []string, err error)
	CalculateFeeSendingToken() (fee string, errorOutput []string, err error)
	CalculateOutPutSendingToken() (string, error)
}
