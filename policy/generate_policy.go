package policy

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

const paymentAddrFileName = "payment.addr"

const keyGenCardano = "cardano-cli address key-gen " +
	"--verification-key-file payment.vkey " +
	"--signing-key-file payment.skey"

const addressBuildCardano = "cardano-cli address build " +
	"--payment-verification-key-file payment.vkey " +
	"--out-file %s --testnet-magic %s"

func GeneratePaymentAddr(id string) (string, error) {
	err := exec.Command(keyGenCardano).Run()
	if err != nil {
		log.Println(err)
		return "", err
	}

	err = exec.Command(fmt.Sprintf(addressBuildCardano, paymentAddrFileName, id)).Run()
	if err != nil {
		log.Println(err)
		return "", err
	}

	fileContent, err := ioutil.ReadFile(paymentAddrFileName)
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
		"--verification-key-file %s/%s " +
		"--signing-key-file %s/%s"

	keyHashGen = "cardano-cli address key-hash " +
		"--payment-verification-key-file %s/%s"

	paymentVerificationFilename = "payment.vkey"
	paymentSigningFilename      = "payment.skey"
	policyScriptFilename        = "policy.script"
	policyDirName               = "policy"
)

func GeneratePolicy() (signingKeyFilePath, verificationKeyFilePath, policyScriptFilePath string, err error) {
	if err := os.Mkdir(policyDirName, 0755); err != nil {
		log.Println(err)
		return
	}
	err = exec.Command(fmt.Sprintf(keyGenPolicy, policyDirName, paymentVerificationFilename,
		policyDirName, paymentSigningFilename)).Run()
	if err != nil {
		log.Println(err)
		return
	}

	signingKeyFilePath = policyDirName + "/" + paymentSigningFilename
	verificationKeyFilePath = policyDirName + "/" + paymentVerificationFilename
	policyScriptFilePath = policyDirName + "/" + policyScriptFilename

	err = exec.Command(fmt.Sprintf(keyHashGen, policyDirName, paymentVerificationFilename)).Run()
	if err != nil {
		log.Println(err)
		return
	}

	policyScript, err := os.Create(policyScriptFilePath)
	if err != nil {
		log.Println(err)
		return
	}

	content, err := ioutil.ReadFile(verificationKeyFilePath)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = policyScript.WriteString(fmt.Sprintf(scriptContent, string(content)))
	if err != nil {
		log.Println(err)
		return
	}

	return signingKeyFilePath, verificationKeyFilePath, policyScriptFilePath, nil
}

const policyIdGen = "cardano-cli transaction policyid" +
	" --script-file ./policy/policy.script >> policy/policyID"

func GeneratePolicyID() (*os.File, error) {
	err := exec.Command("cardano-cli transaction policyid" +
		" --script-file ./policy/policy.script >> policy/policyID").Run()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	file, err := os.Open("policy/policyID")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return file, nil
}
