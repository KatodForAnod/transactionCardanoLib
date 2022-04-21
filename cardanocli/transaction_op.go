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

func (c *CardanoLib) CardanoQueryUtxo() (cliOutPut string, errorOutput []string, err error) {
	var buf bytes.Buffer
	cmd := exec.Command("cardano-cli", "query", "utxo",
		"--address", c.TransactionParams.PaymentAddr, "--testnet-magic", c.TransactionParams.ID)
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

// TransactionBuild - tokenName1 and tokenName2 must be in base16
func (c *CardanoLib) TransactionBuild(tokens []config.Token) (errorOutput []string, err error) {
	if len(tokens) < 1 {
		return errorOutput, errors.New("")
	}

	txOut := fmt.Sprintf("%s+%s+", c.TransactionParams.PaymentAddr, c.TransactionParams.Output)
	mint := fmt.Sprintf("%s %s.%s", tokens[0].TokenAmount, c.TransactionParams.PolicyID, tokens[0].TokenName)
	for i := 1; i < len(tokens); i++ {
		mint += fmt.Sprintf(" + %s %s.%s",
			tokens[i].TokenAmount, c.TransactionParams.PolicyID, tokens[i].TokenName)
	}
	txOut += mint

	cmd := exec.Command("cardano-cli", "transaction", "build-raw",
		"--fee", c.TransactionParams.Fee, "--tx-in",
		c.TransactionParams.TxHash+"#"+c.TransactionParams.Txix, "--tx-out", txOut, "--mint="+mint,
		"--minting-script-file", PolicyScriptFile,
		"--out-file", RawTransactionFile)
	stderr, _ := cmd.StderrPipe()

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

func (c *CardanoLib) CalculateFee() (fee string, errorOutput []string, err error) {
	var buf bytes.Buffer
	cmd := exec.Command("cardano-cli", "transaction", "calculate-min-fee",
		"--tx-body-file", RawTransactionFile, "--tx-in-count", "1",
		"--tx-out-count", "1", "--witness-count", "2", "--testnet-magic", c.TransactionParams.ID,
		"--protocol-params-file", ProtocolParametersFile)
	cmd.Stdout = &buf
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return "", errorOutput, err
	}

	cmd.Wait()

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		errorOutput = append(errorOutput, scanner.Text())
	}

	if len(errorOutput) > 0 {
		return "", errorOutput, fmt.Errorf("CalculateFee error")
	}

	arr := strings.Split(buf.String(), " ")
	if len(arr) < 2 {
		return "", errorOutput, errors.New("split error")
	}

	return arr[0], errorOutput, err
}

func (c *CardanoLib) CalculateOutPut() (string, error) {
	funds, err := strconv.ParseInt(c.TransactionParams.Funds, 10, 64)
	if err != nil {
		log.Println(err)
		return "", err
	}

	fee, err := strconv.ParseInt(c.TransactionParams.Fee, 10, 64)
	if err != nil {
		log.Println(err)
		return "", err
	}

	output := funds - fee
	return strconv.Itoa(int(output)), nil
}

func (c *CardanoLib) TransactionSign() (errorOutput []string, err error) {
	cmd := exec.Command("cardano-cli", "transaction", "sign",
		"--signing-key-file", PaymentSignKeyFile,
		"--signing-key-file", PolicySigningKeyFile,
		"--testnet-magic", c.TransactionParams.ID, "--tx-body-file", RawTransactionFile,
		"--out-file", SignedTransactionFile)
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

func (c *CardanoLib) TransactionSubmit() (errorOutput []string, err error) {
	cmd := exec.Command("cardano-cli", "transaction", "submit",
		"--tx-file", SignedTransactionFile,
		"--testnet-magic", c.TransactionParams.ID)
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

func (c *CardanoLib) TransactionBuildSendingToken(tokens []config.Token,
	sendToken config.Token) (errorOutput []string, err error) {
	for i, token := range tokens {
		if token.TokenName == sendToken.TokenName {
			tokenAll, err := strconv.ParseInt(token.TokenAmount, 10, 64)
			if err != nil {
				log.Println(err)
				return errorOutput, err
			}

			tokenSendAmount, err := strconv.ParseInt(sendToken.TokenAmount, 10, 64)
			if err != nil {
				log.Println(err)
				return errorOutput, err
			}

			amountLeft := strconv.Itoa(int(tokenAll - tokenSendAmount))
			tokens[i].TokenAmount = amountLeft
			break
		} else if i == len(tokens) {
			return errorOutput, errors.New("token not found")
		}
	}

	txOut := fmt.Sprintf("%s+%s+", c.TransactionParams.Receiver, c.TransactionParams.ReceiverOutput)
	txOut += fmt.Sprintf("%s %s.%s", sendToken.TokenAmount, c.TransactionParams.PolicyID, sendToken.TokenName)

	txOut2 := fmt.Sprintf("%s+%s", c.TransactionParams.PaymentAddr,
		c.TransactionParams.Output)
	txOut2 += fmt.Sprintf("+%s %s.%s",
		tokens[0].TokenAmount, c.TransactionParams.PolicyID, tokens[0].TokenName)

	for i := 1; i < len(tokens); i++ {
		txOut2 += fmt.Sprintf("+ %s %s.%s",
			tokens[i].TokenAmount, c.TransactionParams.PolicyID, tokens[i].TokenName)
	}

	cmd := exec.Command("cardano-cli", "transaction", "build-raw",
		"--fee", c.TransactionParams.Fee,
		"--tx-in", c.TransactionParams.TxHash+"#"+c.TransactionParams.Txix,
		"--tx-out", txOut,
		"--tx-out", txOut2,
		"--out-file", RawTransactionSendTokenFile)
	stderr, _ := cmd.StderrPipe()

	fmt.Println(cmd.String())
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

func (c *CardanoLib) TransactionSignSendingToken() (errorOutput []string, err error) {
	cmd := exec.Command("cardano-cli", "transaction", "sign",
		"--signing-key-file", PaymentSignKeyFile,
		"--testnet-magic", c.TransactionParams.ID,
		"--tx-body-file", RawTransactionSendTokenFile,
		"--out-file", SignedTransactionSendTokenFile)
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

func (c *CardanoLib) TransactionSendTokenSubmit() (errorOutput []string, err error) {
	cmd := exec.Command("cardano-cli", "transaction", "submit",
		"--tx-file", SignedTransactionSendTokenFile,
		"--testnet-magic", c.TransactionParams.ID)
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

func (c *CardanoLib) CalculateFeeSendingToken() (fee string, errorOutput []string, err error) {
	var buf bytes.Buffer
	cmd := exec.Command("cardano-cli", "transaction", "calculate-min-fee",
		"--tx-body-file", RawTransactionFile, "--tx-in-count", "1",
		"--tx-out-count", "2", "--witness-count", "1", "--testnet-magic", c.TransactionParams.ID,
		"--protocol-params-file", ProtocolParametersFile)
	cmd.Stdout = &buf
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return "", errorOutput, err
	}

	cmd.Wait()

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		errorOutput = append(errorOutput, scanner.Text())
	}

	if len(errorOutput) > 0 {
		return "", errorOutput, fmt.Errorf("CalculateFee error")
	}

	arr := strings.Split(buf.String(), " ")
	if len(arr) < 2 {
		return "", errorOutput, errors.New("split error")
	}

	return arr[0], errorOutput, err
}

func (c *CardanoLib) CalculateOutPutSendingToken() (string, error) {
	funds, err := strconv.ParseInt(c.TransactionParams.Funds, 10, 64)
	if err != nil {
		log.Println(err)
		return "", err
	}

	fee, err := strconv.ParseInt(c.TransactionParams.Fee, 10, 64)
	if err != nil {
		log.Println(err)
		return "", err
	}

	receiverOutput, err := strconv.ParseInt(c.TransactionParams.ReceiverOutput, 10, 64)
	if err != nil {
		log.Println(err)
		return "", err
	}

	output := funds - fee - receiverOutput
	return strconv.Itoa(int(output)), nil
}
