package cardanocli

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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

	return PaymentVerifyKeyFile, PaymentSignKeyFile, PaymentAddrFile, nil
}

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

	policyScript, err := os.Create(PolicyScriptFile)
	if err != nil {
		log.Println(err)
		return
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
	policyScript.WriteString("}\n")

	return PolicyVerificationkeyFile, PolicySigningKeyFile, PolicyScriptFile, nil
}

const policyIdGen = "cardano-cli transaction policyid" +
	" --script-file ./" + PolicyScriptFile + ">> " + PolicyIDFile

func GeneratePolicyID() (*os.File, error) {
	cmd := exec.Command("cardano-cli", "transaction", "policyid",
		"--script-file", "./policy/policy.script")
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	file, err := os.Open(PolicyIDFile)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return file, nil
}
