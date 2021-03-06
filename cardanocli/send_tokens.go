package cardanocli

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"transactionCardanoLib/config"
)

type SendTokens struct {
	SuperTransactionClass
}

func (c *SendTokens) TransactionBuild(tokens []config.Token,
	sendToken []config.Token) (errorOutput []string, err error) {
	var copyTokens []config.Token
	for _, token := range tokens {
		copyTokens = append(copyTokens, token)
	}

	cmd := exec.Command("cardano-cli", "query", "protocol-parameters",
		"--testnet-magic", c.base.ID,
		"--out-file", c.f.GetProtocolParametersFile())
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return errorOutput, err
	}

	for i, token := range copyTokens {
		if token.TokenName == sendToken[0].TokenName {
			tokenAll, err := strconv.ParseInt(token.TokenAmount, 10, 64)
			if err != nil {
				log.Println(err)
				return errorOutput, err
			}

			tokenSendAmount, err := strconv.ParseInt(sendToken[0].TokenAmount, 10, 64)
			if err != nil {
				log.Println(err)
				return errorOutput, err
			}

			amountLeft := strconv.Itoa(int(tokenAll - tokenSendAmount))
			copyTokens[i].TokenAmount = amountLeft
			break
		} else if i == len(copyTokens) {
			return errorOutput, errors.New("token not found")
		}
	}

	txOut := fmt.Sprintf("%s+%s+", c.processParams.Receiver, c.processParams.ReceiverOutput)
	txOut += fmt.Sprintf("%s %s.%s", sendToken[0].TokenAmount, c.base.PolicyID, sendToken[0].TokenName)

	txOut2 := fmt.Sprintf("%s+%s", c.base.PaymentAddr,
		c.processParams.Output)
	txOut2 += fmt.Sprintf("+%s %s.%s",
		copyTokens[0].TokenAmount, c.base.PolicyID, copyTokens[0].TokenName)

	for i := 1; i < len(copyTokens); i++ {
		txOut2 += fmt.Sprintf("+ %s %s.%s",
			copyTokens[i].TokenAmount, c.base.PolicyID, copyTokens[i].TokenName)
	}

	cmd = exec.Command("cardano-cli", "transaction", "build-raw",
		"--fee", c.processParams.Fee,
		"--tx-in", c.processParams.TxHash+"#"+c.processParams.Txix,
		"--tx-out", txOut,
		"--tx-out", txOut2,
		"--out-file", c.f.GetRawTransactionSendTokenFile())
	stderr, _ = cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return errorOutput, err
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		errorOutput = append(errorOutput, scanner.Text())
	}

	if len(errorOutput) > 0 {
		return errorOutput, fmt.Errorf("TransactionBuild error")
	}

	return errorOutput, nil
}

func (c *SendTokens) TransactionSign() (errorOutput []string, err error) {
	cmd := exec.Command("cardano-cli", "transaction", "sign",
		"--signing-key-file", c.f.GetPaymentSignKeyFile(),
		"--testnet-magic", c.base.ID,
		"--tx-body-file", c.f.GetRawTransactionSendTokenFile(),
		"--out-file", c.f.GetSignedTransactionSendTokenFile())
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return errorOutput, err
	}

	cmd.Wait()

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		errorOutput = append(errorOutput, scanner.Text())
	}

	if len(errorOutput) > 0 {
		return errorOutput, fmt.Errorf("TransactionSign error")
	}

	return errorOutput, nil
}

func (c *SendTokens) TransactionSubmit() (errorOutput []string, err error) {
	cmd := exec.Command("cardano-cli", "transaction", "submit",
		"--tx-file", c.f.GetSignedTransactionSendTokenFile(),
		"--testnet-magic", c.base.ID)
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return errorOutput, err
	}

	cmd.Wait()

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		errorOutput = append(errorOutput, scanner.Text())
	}

	if len(errorOutput) > 0 {
		return []string{}, fmt.Errorf("TransactionSubmit error")
	}

	return []string{}, nil
}

func (c *SendTokens) CalculateFee() (errorOutput []string, err error) {
	var buf bytes.Buffer
	cmd := exec.Command("cardano-cli", "transaction", "calculate-min-fee",
		"--tx-body-file", c.f.GetRawTransactionFile(), "--tx-in-count", "1",
		"--tx-out-count", "2", "--witness-count", "1", "--testnet-magic", c.base.ID,
		"--protocol-params-file", c.f.GetProtocolParametersFile())
	cmd.Stdout = &buf
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return errorOutput, err
	}

	cmd.Wait()

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		errorOutput = append(errorOutput, scanner.Text())
	}

	if len(errorOutput) > 0 {
		return errorOutput, fmt.Errorf("CalculateFee error")
	}

	arr := strings.Split(buf.String(), " ")
	if len(arr) < 2 {
		return errorOutput, errors.New("split error")
	}

	c.processParams.Fee = arr[0]
	return errorOutput, err
}

func (c *SendTokens) CalculateOutPut() error {
	funds, err := strconv.ParseInt(c.processParams.Funds, 10, 64)
	if err != nil {
		log.Println(err)
		return err
	}

	fee, err := strconv.ParseInt(c.processParams.Fee, 10, 64)
	if err != nil {
		log.Println(err)
		return err
	}

	receiverOutput, err := strconv.ParseInt(c.processParams.ReceiverOutput, 10, 64)
	if err != nil {
		log.Println(err)
		return err
	}

	output := funds - fee - receiverOutput
	c.processParams.Output = strconv.Itoa(int(output))
	return nil
}
