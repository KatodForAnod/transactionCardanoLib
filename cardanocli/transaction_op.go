package cardanocli

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"transactionCardanoLib/config"
)

func (c *CardanoLib) InitCardanoQueryUtxo(id string) (cliOutPut string, err error) {
	addr, err := os.ReadFile(c.PaymentAddrFile)
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

const transactionBuildTmpl = "cardano-cli transaction build-raw " +
	"--fee %s " +
	"--tx-in $%s#$%s " +
	"--tx-out $%s+$%s+\"$%s $%s + $%s $%s\" " +
	"--mint=\"$%s $%s + $%s $%s\" " +
	"--minting-script-file %s " +
	"--out-file " + RawTransactionFile

// TransactionBuild - tokenName1 and tokenName2 must be in base16
func TransactionBuild(fee, txHash, txIx, address, output, tokenAmount,
	tokenName1, tokenName2, policyId, policyScriptFilePath string) error {
	txOut := fmt.Sprintf("%s+%s+%s %s.%s + %s %s.%s", address, output, tokenAmount,
		policyId, tokenName1, tokenAmount, policyId, tokenName2)
	//txOut = strings.ReplaceAll(txOut, "\\", "")

	mint := fmt.Sprintf("%s %s.%s + %s %s.%s", tokenAmount, policyId,
		tokenName1, tokenAmount, policyId, tokenName2)
	mint = strings.ReplaceAll(mint, "\\", "")

	cmd := exec.Command("cardano-cli", "transaction", "build-raw",
		"--fee", fee, "--tx-in", txHash+"#"+txIx, "--tx-out", txOut, "--mint", mint,
		"--minting-script-file", policyScriptFilePath,
		"--out-file", RawTransactionFile)
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
