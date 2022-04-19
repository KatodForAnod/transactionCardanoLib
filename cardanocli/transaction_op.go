package cardanocli

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"transactionCardanoLib/config"
)

func (c *CardanoLib) CardanoQueryUtxo(id string) (cliOutPut string, errorOutput []string, err error) {
	addr, err := os.ReadFile(c.FilePaths.PaymentAddrFile)
	if err != nil {
		log.Println(err)
		return "", errorOutput, err
	}

	var buf bytes.Buffer
	cmd := exec.Command("cardano-cli", "query", "utxo",
		"--address", string(addr), "--testnet-magic", id)
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

	addr, err := os.ReadFile(c.FilePaths.PaymentAddrFile)
	if err != nil {
		log.Println(err)
		return errorOutput, err
	}

	policyId, err := os.ReadFile(c.FilePaths.PolicyIDFile)
	if err != nil {
		log.Println(err)
		return errorOutput, err
	}

	txOut := fmt.Sprintf("%s+%s+", string(addr), c.TransactionParams.Output)
	mint := fmt.Sprintf("%s %s.%s", tokens[0].TokenAmount, string(policyId), tokens[0].TokenName)
	for i := 1; i < len(tokens); i++ {
		mint += fmt.Sprintf(" + %s %s.%s",
			tokens[i].TokenAmount, string(policyId), tokens[i].TokenName)
	}
	txOut += mint

	cmd := exec.Command("cardano-cli", "transaction", "build-raw",
		"--fee", c.TransactionParams.Fee, "--tx-in",
		c.TransactionParams.TxHash+"#"+c.TransactionParams.Txix, "--tx-out", txOut, "--mint="+mint,
		"--minting-script-file", c.FilePaths.PolicyScriptFile,
		"--out-file", c.FilePaths.RawTransactionFile)
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

func (c *CardanoLib) CalculateFee(id string) (fee string, errorOutput []string, err error) {
	var buf bytes.Buffer
	cmd := exec.Command("cardano-cli", "transaction", "calculate-min-fee",
		"--tx-body-file", c.FilePaths.RawTransactionFile, "--tx-in-count", "1",
		"--tx-out-count", "1", "--witness-count", "2", "--testnet-magic", id,
		"--protocol-params-file", c.FilePaths.ProtocolParametersFile)
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

func (c *CardanoLib) TransactionSign(id string) (errorOutput []string, err error) {
	cmd := exec.Command("cardano-cli", "transaction", "sign",
		"--signing-key-file", c.FilePaths.PaymentSignKeyFile,
		"--signing-key-file", c.FilePaths.PolicySigningKeyFile,
		"--testnet-magic", id, "--tx-body-file", c.FilePaths.RawTransactionFile,
		"--out-file", c.FilePaths.SignedTransactionFile)
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

func (c *CardanoLib) TransactionSubmit(id string) (errorOutput []string, err error) {
	cmd := exec.Command("cardano-cli", "transaction", "submit",
		"--tx-file", c.FilePaths.SignedTransactionFile,
		"--testnet-magic", id)
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

/*func (c *CardanoLib) TransactionBuildSendingToken(receiverAddr,
	receiverOutput string) (errorOutput []string, err error) {
	policyId, err := os.ReadFile(c.FilePaths.PolicyIDFile)
	if err != nil {
		log.Println(err)
		return errorOutput, err
	}

	txOut := fmt.Sprintf("%s+%s+", string(addr), c.TransactionParams.Output)
	mint := fmt.Sprintf("%s %s.%s", c.TransactionParams.TokenAmount, string(policyId), tokenName[0])
	for i := 1; i < len(tokenName); i++ {
		mint += fmt.Sprintf(" + %s %s.%s",
			c.TransactionParams.TokenAmount, string(policyId), tokenName[i])
	}
	txOut += mint

	cmd := exec.Command("cardano-cli", "transaction", "build-raw",
		"--fee", c.TransactionParams.Fee, "--tx-in",
		c.TransactionParams.TxHash+"#"+c.TransactionParams.Txix, "--tx-out", txOut, "--mint="+mint,
		"--minting-script-file", c.FilePaths.PolicyScriptFile,
		"--out-file", c.FilePaths.RawTransactionFile)
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
}*/
