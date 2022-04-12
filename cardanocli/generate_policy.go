package cardanocli

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
)

type CardanoLib struct {
	PaymentVerifyKeyFile      string
	PaymentSignKeyFile        string
	PaymentAddrFile           string
	PolicyVerificationkeyFile string
	PolicySigningKeyFile      string
	PolicyDirName             string
	PolicyScriptFile          string
	PolicyIDFile              string
	ProtocolParametersFile    string
}

func (c *CardanoLib) GeneratePaymentFiles(id string) (err error) {
	err = exec.Command("cardano-cli", "address", "key-gen",
		"--verification-key-file", c.PaymentVerifyKeyFile, "--signing-key-file", c.PaymentSignKeyFile).Run()
	if err != nil {
		log.Println(err)
		return err
	}

	err = exec.Command("cardano-cli", "address", "build", "--payment-verification-key-file",
		c.PaymentVerifyKeyFile, "--out-file", c.PaymentAddrFile, "--testnet-magic", id).Run()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *CardanoLib) GeneratePolicy() (err error) {
	if err = os.MkdirAll("./"+c.PolicyDirName, os.ModePerm); err != nil {
		log.Println(err)
		return err
	}

	err = exec.Command("cardano-cli", "address", "key-gen", "--verification-key-file",
		c.PolicyVerificationkeyFile, "--signing-key-file", c.PolicySigningKeyFile).Run()
	if err != nil {
		log.Println(err)
		return err
	}

	policyScript, err := os.Create(c.PolicyScriptFile)
	if err != nil {
		log.Println(err)
		return err
	}
	defer policyScript.Close()

	policyScript.WriteString("{\n")
	policyScript.WriteString("  \"keyHash\": \"")

	var buf bytes.Buffer
	cmd := exec.Command("cardano-cli", "address", "key-hash",
		"--payment-verification-key-file", c.PolicyVerificationkeyFile)
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
	policyIdFile, err := os.Create(c.PolicyIDFile)
	if err != nil {
		log.Println(err)
		return err
	}
	defer policyIdFile.Close()

	var buf bytes.Buffer
	cmd := exec.Command("cardano-cli", "transaction", "policyid",
		"--script-file", "./"+c.PolicyScriptFile)
	cmd.Stdout = &buf
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	cmd.Wait()

	policyIdFile.WriteString(strings.ReplaceAll(buf.String(), "\n", ""))

	return nil
}

func (c *CardanoLib) GenerateProtocol(id string) error {
	err := exec.Command("cardano-cli", "query", "key-gen", "protocol-parameters",
		"--testnet-magic", id, "--out-file", c.ProtocolParametersFile).Run()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
