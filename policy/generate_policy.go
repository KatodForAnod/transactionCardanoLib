package policy

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func GeneratePaymentAddr(id string) (*os.File, error) {
	err := exec.Command("cardano-cli address key-gen " +
		"--verification-key-file payment.vkey --signing-key-file payment.skey").Run()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = exec.Command("cardano-cli address build "+
		"--payment-verification-key-file payment.vkey --out-file payment.addr --testnet-magic ", id).Run()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	file, err := os.Open("payment.addr")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return file, nil
}

const scriptContent = "{\n" +
	"\"keyHash\": \"%s\"," +
	"\"type\": \"sig\"," +
	"}"

func GeneratePolicy() (signingKey, verificationKey, policyScript *os.File, err error) {
	if err := os.Mkdir("policy", 0755); err != nil {
		log.Println(err)
		return
	}
	err = exec.Command("cardano-cli address key-gen " +
		"--verification-key-file policy/payment.vkey --signing-key-file policy/payment.skey").Run()
	if err != nil {
		log.Println(err)
		return
	}

	verificationKey, err = os.Open("policy/payment.vkey")
	if err != nil {
		log.Println(err)
		return
	}

	signingKey, err = os.Open("policy/payment.skey")
	if err != nil {
		log.Println(err)
		return
	}

	policyScript, err = os.Create("policy/policy.script")
	if err != nil {
		log.Println(err)
		return
	}

	err = exec.Command("cardano-cli address key-hash " +
		"--payment-verification-key-file policy/policy.vkey").Run()
	if err != nil {
		log.Println(err)
		return
	}

	content, err := ioutil.ReadFile("policy/policy.vkey")
	if err != nil {
		log.Println(err)
		return
	}

	_, err = policyScript.WriteString(fmt.Sprintf(scriptContent, string(content)))
	if err != nil {
		log.Println(err)
		return
	}

	return
}

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
