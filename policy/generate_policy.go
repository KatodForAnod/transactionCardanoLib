package policy

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

const keyGenCardano = "cardano-cli address key-gen " +
	"--verification-key-file " + PaymentVerifyKeyFile + " " +
	"--signing-key-file " + PaymentSignKeyFile

const addressBuildCardano = "cardano-cli address build " +
	"--payment-verification-key-file " + PaymentVerifyKeyFile + " " +
	"--out-file %s --testnet-magic %s"

func GeneratePaymentAddr(id string) (string, error) {
	err := exec.Command("cardano-cli", "address", "key-gen", "--verification-key-file",
		PaymentVerifyKeyFile, "--signing-key-file", PaymentSignKeyFile).Run()
	if err != nil {
		log.Println(err)
		return "", err
	}

	err = exec.Command("cardano-cli", "address", "build", "--payment-verification-key-file",
		PaymentVerifyKeyFile, "--out-file", PaymentAddrFile, "--testnet-magic", id).Run()
	if err != nil {
		log.Println(err)
		return "", err
	}

	fileContent, err := ioutil.ReadFile(PaymentAddrFile)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(fileContent), nil
}

const (
	scriptContent = "{\n" +
		"\"keyHash\": \"%s\"," +
		"\"type\": \"sig\"," +
		"}"

	keyGenPolicy = "cardano-cli address key-gen " +
		"--verification-key-file %s" +
		"--signing-key-file %s"

	keyHashGen = "cardano-cli address key-hash " +
		"--payment-verification-key-file %s"
)

func GeneratePolicy() (err error) {
	if err := os.Mkdir(PolicyDirName, 0755); err != nil {
		log.Println(err)
		return err
	}
	err = exec.Command(fmt.Sprintf(keyGenPolicy, PaymentVerifyKeyFile,
		PaymentSignKeyFile)).Run()
	if err != nil {
		log.Println(err)
		return
	}

	err = exec.Command(fmt.Sprintf(keyHashGen, PaymentVerifyKeyFile)).Run()
	if err != nil {
		log.Println(err)
		return
	}

	policyScript, err := os.Create(PolicyScriptFile)
	if err != nil {
		log.Println(err)
		return
	}

	content, err := ioutil.ReadFile(PolicyVerificationkeyFile)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = policyScript.WriteString(fmt.Sprintf(scriptContent, string(content)))
	if err != nil {
		log.Println(err)
		return
	}

	return nil
}

const policyIdGen = "cardano-cli transaction policyid" +
	" --script-file ./" + PolicyScriptFile + ">> " + PolicyIDFile

func GeneratePolicyID() (*os.File, error) {
	err := exec.Command("cardano-cli transaction policyid" +
		" --script-file ./policy/policy.script >> policy/policyID").Run()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	file, err := os.Open(PolicyIDFile)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return file, nil
}
