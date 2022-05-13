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
	"transactionCardanoLib/files"
)

type CreateTokens struct {
	base          BaseTransactionParams
	processParams TransactionParams
	f             files.Files
}

func (c *CreateTokens) Init(base BaseTransactionParams,
	processParams TransactionParams,
	f files.Files) {
	c.base = base
	c.f = f
	c.processParams = processParams
}

func (c *CreateTokens) SetBaseParams(base BaseTransactionParams) {
	c.base = base
}

func (c *CreateTokens) SetProcessParams(processParams TransactionParams) {
	c.processParams = processParams
}

func (c *CreateTokens) SetFileParams(f files.Files) {
	c.f = f
}

func (c *CreateTokens) CardanoQueryUtxo() (cliOutPut string, errorOutput []string, err error) {
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

// TransactionBuild - tokenName1 and tokenName2 must be in base16
func (c *CreateTokens) TransactionBuild(oldTokens []config.Token, newTokens []config.Token) (errorOutput []string, err error) {
	if len(newTokens) < 1 || len(oldTokens) < 1 {
		return errorOutput, errors.New("")
	}

	cmd := exec.Command("cardano-cli", "query", "protocol-parameters",
		"--testnet-magic", c.base.ID,
		"--out-file", c.f.GetProtocolParametersFile())
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return errorOutput, err
	}

	txOut := fmt.Sprintf("%s+%s+", c.base.PaymentAddr, c.processParams.Output)
	var oldMint string
	for i := 0; i < len(oldTokens); i++ {
		oldMint += fmt.Sprintf("%s %s.%s + ",
			oldTokens[i].TokenAmount, c.base.PolicyID, oldTokens[i].TokenName)
	}
	txOut += oldMint

	mint := fmt.Sprintf("%s %s.%s", newTokens[0].TokenAmount, c.base.PolicyID, newTokens[0].TokenName)
	for i := 1; i < len(newTokens); i++ {
		mint += fmt.Sprintf(" + %s %s.%s",
			newTokens[i].TokenAmount, c.base.PolicyID, newTokens[i].TokenName)
	}
	txOut += mint

	cmd = exec.Command("cardano-cli", "transaction", "build-raw",
		"--fee", c.processParams.Fee, "--tx-in",
		c.processParams.TxHash+"#"+c.processParams.Txix, "--tx-out", txOut, "--mint="+mint,
		"--minting-script-file", c.f.GetPolicyScriptFile(),
		"--out-file", c.f.GetRawTransactionFile())
	stderr, _ = cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return errorOutput, err
	}
	fmt.Println(cmd.String())
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		errorOutput = append(errorOutput, scanner.Text())
	}

	if len(errorOutput) > 0 {
		return errorOutput, fmt.Errorf("TransactionBuild error")
	}

	return errorOutput, nil
}

func (c *CreateTokens) CalculateFee() (errorOutput []string, err error) {
	var buf bytes.Buffer
	cmd := exec.Command("cardano-cli", "transaction", "calculate-min-fee",
		"--tx-body-file", c.f.GetRawTransactionFile(), "--tx-in-count", "1",
		"--tx-out-count", "1", "--witness-count", "2", "--testnet-magic", c.base.ID,
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

func (c *CreateTokens) CalculateOutPut() error {
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

	output := funds - fee
	c.processParams.Output = strconv.Itoa(int(output))
	return nil
}

func (c *CreateTokens) TransactionSign() (errorOutput []string, err error) {
	cmd := exec.Command("cardano-cli", "transaction", "sign",
		"--signing-key-file", c.f.GetPaymentSignKeyFile(),
		"--signing-key-file", c.f.GetPolicySigningKeyFile(),
		"--testnet-magic", c.base.ID, "--tx-body-file", c.f.GetRawTransactionFile(),
		"--out-file", c.f.GetSignedTransactionFile())
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

func (c *CreateTokens) TransactionSubmit() (errorOutput []string, err error) {
	cmd := exec.Command("cardano-cli", "transaction", "submit",
		"--tx-file", c.f.GetSignedTransactionFile(),
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
