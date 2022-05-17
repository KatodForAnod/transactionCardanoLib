package cardanocli

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"transactionCardanoLib/config"
	"transactionCardanoLib/files"
)

type NftPolicy struct {
	Policy
	slot string
}

func (commPolicy *NftPolicy) Init(f files.Files, conf config.Config) {
	commPolicy.Policy.Init(f, conf)
	commPolicy.slot = "0" //TODO
}

func (commPolicy *NftPolicy) generatePolicyScript() error {
	if err := os.MkdirAll("./"+commPolicy.f.GetPolicyDirName(), os.ModePerm); err != nil {
		log.Println(err)
		return err
	}

	policyScript, err := os.Create(commPolicy.f.GetPolicyScriptFile())
	if err != nil {
		log.Println(err)
		return err
	}
	defer policyScript.Close()

	policyScript.WriteString("{\n")
	policyScript.WriteString("  \"type\": \"all\",")
	policyScript.WriteString("  \"scripts\":")
	policyScript.WriteString("  [")
	policyScript.WriteString("   {")
	policyScript.WriteString("  \"type\": \"before\",\n")
	policyScript.WriteString("  \"slot\":" + commPolicy.slot + ",\n")
	policyScript.WriteString("},")
	policyScript.WriteString("   {")
	policyScript.WriteString("  \"type\": \"sig\",\n")
	policyScript.WriteString("  \"keyHash\": \"")

	var buf bytes.Buffer
	cmd := exec.Command("cardano-cli", "address", "key-hash",
		"--payment-verification-key-file", commPolicy.f.GetPolicyVerificationkeyFile())
	cmd.Stdout = &buf
	if err = cmd.Start(); err != nil {
		panic(err)
	}
	cmd.Wait()

	keyHash := strings.ReplaceAll(buf.String(), "\n", "")
	policyScript.WriteString(keyHash + "\"\n")
	policyScript.WriteString("}")
	policyScript.WriteString("]")
	policyScript.WriteString("}")
	return nil
}

func (commPolicy *NftPolicy) generatePolicyID() error {
	policyIdFile, err := os.Create(commPolicy.f.GetPolicyIDFile())
	if err != nil {
		log.Println(err)
		return err
	}
	defer policyIdFile.Close()

	var buf bytes.Buffer
	cmd := exec.Command("cardano-cli", "transaction", "policyid",
		"--script-file", "./"+commPolicy.f.GetPolicyScriptFile())
	cmd.Stdout = &buf
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	cmd.Wait()

	body := strings.ReplaceAll(buf.String(), "\n", "")
	policyIdFile.WriteString(body)

	commPolicy.policyID = body

	return nil
}

func (commPolicy *NftPolicy) GeneratePolicyFiles() error {
	err := commPolicy.generateProtocolParameters()
	if err != nil {
		log.Println(err)
		return err
	}

	err = commPolicy.generatePolicyScript()
	if err != nil {
		log.Println(err)
		return err
	}

	err = commPolicy.generatePolicyID()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
