package cardanocli

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
)

type CardanoLib struct {
	TransactionParams TransactionParams
}

type TransactionParams struct {
	TxHash      string
	Txix        string
	Funds       string
	Fee         string
	Output      string
	PaymentAddr string
	PolicyID    string
}

func (c *CardanoLib) GeneratePaymentFiles(id string) (err error) {
	err = exec.Command("cardano-cli", "address", "key-gen",
		"--verification-key-file", PaymentVerifyKeyFile,
		"--signing-key-file", PaymentSignKeyFile).Run()
	if err != nil {
		log.Println(err)
		return err
	}

	err = exec.Command("cardano-cli", "address", "build", "--payment-verification-key-file",
		PaymentVerifyKeyFile, "--out-file", PaymentAddrFile, "--testnet-magic", id).Run()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *CardanoLib) GeneratePolicy() (err error) {
	if err = os.MkdirAll("./"+PolicyDirName, os.ModePerm); err != nil {
		log.Println(err)
		return err
	}

	err = exec.Command("cardano-cli", "address", "key-gen", "--verification-key-file",
		PolicyVerificationkeyFile, "--signing-key-file", PolicySigningKeyFile).Run()
	if err != nil {
		log.Println(err)
		return err
	}

	policyScript, err := os.Create(PolicyScriptFile)
	if err != nil {
		log.Println(err)
		return err
	}
	defer policyScript.Close()

	policyScript.WriteString("{\n")
	policyScript.WriteString("  \"keyHash\": \"")

	var buf bytes.Buffer
	cmd := exec.Command("cardano-cli", "address", "key-hash",
		"--payment-verification-key-file", PolicyVerificationkeyFile)
	cmd.Stdout = &buf
	if err = cmd.Start(); err != nil {
		panic(err)
	}
	cmd.Wait()

	keyHash := strings.ReplaceAll(buf.String(), "\n", "")
	policyScript.WriteString(keyHash + "\",\n")
	policyScript.WriteString("  \"type\": \"sig\"\n")
	policyScript.WriteString("}")

	return nil
}

func (c *CardanoLib) GeneratePolicyID() error {
	policyIdFile, err := os.Create(PolicyIDFile)
	if err != nil {
		log.Println(err)
		return err
	}
	defer policyIdFile.Close()

	var buf bytes.Buffer
	cmd := exec.Command("cardano-cli", "transaction", "policyid",
		"--script-file", "./"+PolicyScriptFile)
	cmd.Stdout = &buf
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	cmd.Wait()

	policyIdFile.WriteString(strings.ReplaceAll(buf.String(), "\n", ""))

	return nil
}

func (c *CardanoLib) GenerateProtocol(id string) error {
	err := exec.Command("cardano-cli", "query", "protocol-parameters",
		"--testnet-magic", id, "--out-file", ProtocolParametersFile).Run()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
