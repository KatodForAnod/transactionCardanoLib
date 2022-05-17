package cardanocli

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"transactionCardanoLib/files"
)

type PolicyContract interface {
	Init(f files.Files, id string)
	GeneratePolicyFiles() error
}

type Policy struct {
	f  files.Files
	id string

	policyID string
}

func (c *Policy) Init(f files.Files, id string) {
	c.f = f
	c.id = id
}

func (c *Policy) generatePolicyScript() (err error) {
	if err = os.MkdirAll("./"+c.f.GetPolicyDirName(), os.ModePerm); err != nil {
		log.Println(err)
		return err
	}

	err = exec.Command("cardano-cli", "address", "key-gen", "--verification-key-file",
		c.f.GetPolicyVerificationkeyFile(), "--signing-key-file", c.f.GetPolicySigningKeyFile()).Run()
	if err != nil {
		log.Println(err)
		return err
	}

	policyScript, err := os.Create(c.f.GetPolicyScriptFile())
	if err != nil {
		log.Println(err)
		return err
	}
	defer policyScript.Close()

	policyScript.WriteString("{\n")
	policyScript.WriteString("  \"keyHash\": \"")

	var buf bytes.Buffer
	cmd := exec.Command("cardano-cli", "address", "key-hash",
		"--payment-verification-key-file", c.f.GetPolicyVerificationkeyFile())
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

func (c *Policy) generatePolicyID() error {
	policyIdFile, err := os.Create(c.f.GetPolicyIDFile())
	if err != nil {
		log.Println(err)
		return err
	}
	defer policyIdFile.Close()

	var buf bytes.Buffer
	cmd := exec.Command("cardano-cli", "transaction", "policyid",
		"--script-file", "./"+c.f.GetPolicyScriptFile())
	cmd.Stdout = &buf
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	cmd.Wait()

	body := strings.ReplaceAll(buf.String(), "\n", "")
	policyIdFile.WriteString(body)

	c.policyID = body

	return nil
}

func (c *Policy) generateProtocolParameters() error {
	err := exec.Command("cardano-cli", "query", "protocol-parameters",
		"--testnet-magic", c.id, "--out-file", c.f.GetProtocolParametersFile()).Run()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *Policy) GeneratePolicyFiles() error {
	err := c.generateProtocolParameters()
	if err != nil {
		log.Println(err)
		return err
	}

	err = c.generatePolicyScript()
	if err != nil {
		log.Println(err)
		return err
	}

	err = c.generatePolicyID()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
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
