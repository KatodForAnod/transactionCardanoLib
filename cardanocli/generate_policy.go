package cardanocli

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"transactionCardanoLib/config"
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
	ID          string
}

/*func (c *CardanoLib) GeneratePaymentFiles() (err error) {
	err = exec.Command("cardano-cli", "address", "key-gen",
		"--verification-key-file", PaymentVerifyKeyFile,
		"--signing-key-file", PaymentSignKeyFile).Run()
	if err != nil {
		log.Println(err)
		return err
	}

	err = exec.Command("cardano-cli", "address", "build", "--payment-verification-key-file",
		PaymentVerifyKeyFile, "--out-file", PaymentAddrFile, "--testnet-magic", c.TransactionParams.ID).Run()
	if err != nil {
		log.Println(err)
		return err
	}

	addr, err := os.ReadFile(PaymentAddrFile)
	if err != nil {
		log.Println(err)
		return err
	}

	c.TransactionParams.PaymentAddr = string(addr)

	return nil
}*/

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

	body := strings.ReplaceAll(buf.String(), "\n", "")
	policyIdFile.WriteString(body)

	c.TransactionParams.PolicyID = body

	return nil
}

func (c *CardanoLib) GenerateProtocol() error {
	err := exec.Command("cardano-cli", "query", "protocol-parameters",
		"--testnet-magic", c.TransactionParams.ID, "--out-file", ProtocolParametersFile).Run()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *CardanoLib) UseExistPolicy(conf config.Config) error {
	PolicyScriptFile = conf.PolicyScriptFilePath
	PolicySigningKeyFile = conf.PolicySigningFilePath
	PolicyVerificationkeyFile = conf.PolicyVerificationFilePath

	c.TransactionParams.PolicyID = conf.PolicyID

	return nil
}

func (c *CardanoLib) GeneratePolicyFiles() error {
	err := exec.Command("cardano-cli", "address", "key-gen",
		"--verification-key-file", PaymentVerifyKeyFile,
		"--signing-key-file", PaymentSignKeyFile).Run()
	if err != nil {
		log.Println(err)
		return err
	}

	err = c.GenerateProtocol()
	if err != nil {
		log.Println(err)
		return err
	}
	err = c.GeneratePolicy()
	if err != nil {
		log.Println(err)
		return err
	}
	err = c.GeneratePolicyID()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
