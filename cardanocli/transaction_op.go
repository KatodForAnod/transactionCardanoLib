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

func (c *CardanoLib) CardanoQueryUtxo(id string) (cliOutPut string, err error) {
	addr, err := os.ReadFile(c.FilePaths.PaymentAddrFile)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var buf bytes.Buffer
	cmd := exec.Command("cardano-cli", "query", "utxo",
		"--address", string(addr), "--testnet-magic", id)
	cmd.Stdout = &buf
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	cmd.Wait()

	return buf.String(), nil
}

func TransactionSign(id string, token config.TokenStruct) error {
	//comm := fmt.Sprintf(transactionSignTmpl, token.PolicySigningFilePath, id)

	err := exec.Command("cardano-cli", "transaction", "sign", "--signing-key-file",
		PaymentSignKeyFile, "--signing-key-file", token.PolicySigningFilePath,
		"--testnet-magic", id, "--tx-body-file", RawTransactionFile, "--out-file", SignedTransactionFile).Run()
	if err != nil {
		log.Println(err)
		return err
	}
	// get error msg
	return nil
}

// TransactionBuild - tokenName1 and tokenName2 must be in base16
func (c *CardanoLib) TransactionBuild(tokenName []string) error {
	if len(tokenName) < 1 {
		return errors.New("")
	}

	addr, err := os.ReadFile(c.FilePaths.PaymentAddrFile)
	if err != nil {
		log.Println(err)
		return err
	}

	policyId, err := os.ReadFile(c.FilePaths.PolicyIDFile)
	if err != nil {
		log.Println(err)
		return err
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
		return err
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return nil
}

func (c *CardanoLib) CalculateFee(id string) (string, error) {
	var buf bytes.Buffer
	cmd := exec.Command("cardano-cli", "transaction", "calculate-min-fee",
		"--tx-body-file", c.FilePaths.RawTransactionFile, "--tx-in-count", "1",
		"--tx-out-count", "1", "--witness-count", "2", "--testnet-magic", id,
		"--protocol-params-file", c.FilePaths.ProtocolParametersFile)
	cmd.Stdout = &buf
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	cmd.Wait()

	arr := strings.Split(buf.String(), " ")
	if len(arr) < 2 {
		return "", errors.New("split error")
	}

	return arr[0], nil
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

func (c *CardanoLib) TransactionSign(id string) error {
	err := exec.Command("cardano-cli", "transaction", "sign",
		"--signing-key-file", c.FilePaths.PolicySigningKeyFile,
		"--signing-key-file", c.FilePaths.PolicySigningKeyFile,
		"--testnet-magic", id, "--tx-body-file", c.FilePaths.RawTransactionFile,
		"--out-file", c.FilePaths.SignedTransactionFile).Run()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *CardanoLib) TransactionSubmit(id string) error {
	err := exec.Command("cardano-cli", "transaction", "submit",
		"--tx-file", c.FilePaths.SignedTransactionFile,
		"--testnet-magic", id).Run()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
