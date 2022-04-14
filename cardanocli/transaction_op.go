package cardanocli

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"transactionCardanoLib/config"
)

func (c *CardanoLib) InitCardanoQueryUtxo(id string) (cliOutPut string, err error) {
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

const transactionSignTmpl = "cardano-cli transaction sign " +
	"--signing-key-file payment.skey  " +
	"--signing-key-file %s " +
	"--testnet-magic %s " +
	"--tx-body-file " + RawTransactionFile + " " +
	"--out-file " + SignedTransactionFile

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

	txOut := fmt.Sprintf("%s+%s", string(addr), c.TransactionParams.Output)
	mint := fmt.Sprintf("+\"%s %s.%s", c.TransactionParams.TokenAmount, string(policyId), tokenName[0])
	for i := 1; i < len(tokenName); i++ {
		mint += fmt.Sprintf(" + %s %s.%s",
			c.TransactionParams.TokenAmount, string(policyId), tokenName[0])
	}
	mint += "\""
	txOut += mint

	cmd := exec.Command("cardano-cli", "transaction", "build-raw",
		"--Fee", c.TransactionParams.Fee, "--tx-in",
		c.TransactionParams.TxHash+"#"+c.TransactionParams.Txix, "--tx-out", txOut, "--mint=", mint,
		"--minting-script-file", c.FilePaths.PolicyScriptFile,
		"--out-file", c.FilePaths.RawTransactionFile)
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return nil
}

const transactionPreBuildTmpl = "cardano-cli query utxo --address %s --testnet-magic %s"

func TransactionPreBuild(address, id string) (cliOutput string, err error) {
	comm := fmt.Sprintf(transactionPreBuildTmpl, address, id)

	out, err := exec.Command(comm).Output()
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(out), nil
}
