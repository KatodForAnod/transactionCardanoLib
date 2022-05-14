package cardanocli

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"os/exec"
	"strings"
	"transactionCardanoLib/config"
	"transactionCardanoLib/files"
)

type TransactionContract interface {
	TransactionBuild(tokens []config.Token,
		sendToken []config.Token) (errorOutput []string, err error)
	CalculateFee() (errorOutput []string, err error)
	CalculateOutPut() error
	TransactionSign() (errorOutput []string, err error)
	TransactionSubmit() (errorOutput []string, err error)
	CardanoQueryUtxo() (cliOutPut string, errorOutput []string, err error)
	ParseUtxo(cliOutput string) ([]TransactionParams, [][]config.Token, error) // change output vars
	SetBaseParams(base BaseTransactionParams)
	SetProcessParams(processParams TransactionParams)
	SetFileParams(f files.Files)
}

type SuperTransactionClass struct {
	base          BaseTransactionParams
	processParams TransactionParams
	f             files.Files
}

func (c *SuperTransactionClass) Init(base BaseTransactionParams,
	processParams TransactionParams,
	f files.Files) {
	c.base = base
	c.f = f
	c.processParams = processParams
}

func (c *SuperTransactionClass) SetBaseParams(base BaseTransactionParams) {
	c.base = base
}

func (c *SuperTransactionClass) SetProcessParams(processParams TransactionParams) {
	c.processParams = processParams
}

func (c *SuperTransactionClass) SetFileParams(f files.Files) {
	c.f = f
}

func (c *SuperTransactionClass) CardanoQueryUtxo() (cliOutPut string, errorOutput []string, err error) {
	var buf bytes.Buffer
	cmd := exec.Command("cardano-cli", "query", "utxo",
		"--address", c.base.PaymentAddr, "--testnet-magic", c.base.ID)
	cmd.Stdout = &buf
	stderr, _ := cmd.StderrPipe()

	if err = cmd.Start(); err != nil {
		log.Println(err)
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			errorOutput = append(errorOutput, scanner.Text())
		}
		return "", errorOutput, err
	}

	cmd.Wait()

	return buf.String(), errorOutput, nil
}

func (c *SuperTransactionClass) TransactionBuild(tokens []config.Token,
	sendToken []config.Token) (errorOutput []string, err error) {
	return errorOutput, err
}

func (c *SuperTransactionClass) TransactionSign() (errorOutput []string, err error) {
	return errorOutput, err
}

func (c *SuperTransactionClass) TransactionSubmit() (errorOutput []string, err error) {
	return errorOutput, err
}

func (c *SuperTransactionClass) CalculateFee() (errorOutput []string, err error) {
	return errorOutput, err
}

func (c *SuperTransactionClass) CalculateOutPut() error {
	return nil
}

func (c *SuperTransactionClass) ParseUtxo(cliOutput string) ([]TransactionParams, [][]config.Token, error) {
	cliOutputArr := strings.Split(cliOutput, "\n")
	if len(cliOutputArr) < 3 {
		return []TransactionParams{}, nil, errors.New("error split")
	}

	var paramsAllTxHash []TransactionParams
	var tokensAllTxHash [][]config.Token

	for j := 2; j < len(cliOutputArr); j++ {
		str := cliOutputArr[j]
		vars := strings.Fields(str)

		txHash := vars[0]
		txIx := vars[1]
		amountLovelace := vars[2]

		var tokens []config.Token
		for i := 4; i < len(vars)-2; i += 3 {
			policyAndTokenName := strings.Split(vars[i+2], ".")
			tokenName := policyAndTokenName[1]

			token := config.Token{
				TokenName:   tokenName,
				TokenAmount: vars[i+1],
			}

			tokens = append(tokens, token)
		}

		tokensAllTxHash = append(tokensAllTxHash, tokens)
		paramsAllTxHash = append(paramsAllTxHash, TransactionParams{
			TxHash: txHash,
			Txix:   txIx,
			Funds:  amountLovelace,
		})
	}

	return paramsAllTxHash, tokensAllTxHash, nil
}
