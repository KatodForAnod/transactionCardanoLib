package cardanocli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func GeneratePaymentAddr(id string) (verifyFile, signFile,
	paymentAddrFile string, err error) {
	err = exec.Command("cardano-cli", "address", "key-gen",
		"--verification-key-file", PaymentVerifyKeyFile, "--signing-key-file", PaymentSignKeyFile).Run()
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}

	err = exec.Command("cardano-cli", "address", "build", "--payment-verification-key-file",
		PaymentVerifyKeyFile, "--out-file", PaymentAddrFile, "--testnet-magic", id).Run()
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}

	/*fileContent, err := ioutil.ReadFile(PaymentAddrFile)
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}*/

	return PaymentVerifyKeyFile, PaymentSignKeyFile, PaymentAddrFile, nil
}

const (
	scriptContent = "{\n" +
		"\"keyHash\": \"%s\"," +
		"\"type\": \"sig\"" +
		"}"
)

func GeneratePolicy() (verifyFile, signFile, scriptFile string, err error) {
	if err = os.MkdirAll("./"+PolicyDirName, os.ModePerm); err != nil {
		log.Println(err)
		return
	}
	err = exec.Command("cardano-cli", "address", "key-gen", "--verification-key-file",
		PolicyVerificationkeyFile, "--signing-key-file", PolicySigningKeyFile).Run()
	if err != nil {
		log.Println(err)
		return
	}

	err = exec.Command("cardano-cli", "address", "key-hash",
		"--payment-verification-key-file", PolicyVerificationkeyFile).Run()
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

	return PolicyVerificationkeyFile, PolicySigningKeyFile, PolicyScriptFile, nil
}

const policyIdGen = "cardano-cli transaction policyid" +
	" --script-file ./" + PolicyScriptFile + ">> " + PolicyIDFile

func GeneratePolicyID() (*os.File, error) {
	err := exec.Command("cardano-cli transaction policyid" +
		" --script-file ./policy/policy.script >> policy/policyID").Run() // Not checked
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
