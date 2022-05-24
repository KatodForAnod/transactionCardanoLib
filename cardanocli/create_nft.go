package cardanocli

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"transactionCardanoLib/config"
	"transactionCardanoLib/files"
)

type CreateNFT struct {
	SuperTransactionClass
}

func (c *CreateNFT) Init(base BaseTransactionParams,
	processParams TransactionParams,
	f files.Files) {
	c.base = base
	c.f = f
	c.processParams = processParams
}

func (c *CreateNFT) TransactionBuild(tokens []config.Token) (errorOutput []string, err error) {
	txOut := fmt.Sprintf("%s+%s+", c.base.PaymentAddr, c.processParams.Output)
	mint := fmt.Sprintf("%s %s.%s", tokens[0].TokenAmount, c.base.PolicyID, tokens[0].TokenName)
	for i := 1; i < len(tokens); i++ {
		mint += fmt.Sprintf(" + %s %s.%s",
			tokens[i].TokenAmount, c.base.PolicyID, tokens[i].TokenName)
	}
	txOut += mint

	cmd := exec.Command("cardano-cli", "transaction", "build",
		"--mainnet",
		"--alonzo-era",
		"--tx-in", c.processParams.TxHash+"#"+c.processParams.Txix,
		"--change-address", c.base.PaymentAddr,
		"--tx-out", txOut,
		"--mint="+mint,
		"--minting-script-file", c.f.GetPolicyScriptFile(),
		"--metadata-json-file", c.f.GetMetadataAttrFile(),
		"--invalid-hereafter", c.processParams.SlotNumber,
		"--witness-override", "2",
		"--out-file", c.f.GetRawTransactionFile())
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return errorOutput, err
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		errorOutput = append(errorOutput, scanner.Text())
	}

	if len(errorOutput) > 0 {
		return errorOutput, fmt.Errorf("TransactionBuild error")
	}

	return errorOutput, nil
}

func (c *CreateNFT) TransactionSign() (errorOutput []string, err error) {
	cmd := exec.Command("cardano-cli", "transaction", "sign",
		"--signing-key-file", c.f.GetPaymentSignKeyFile(),
		"--signing-key-file", c.f.GetPolicySigningKeyFile(),
		"--mainnet",
		"--tx-body-file", c.f.GetRawTransactionFile(),
		"--out-file", c.f.GetSignedTransactionFile())
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return errorOutput, err
	}

	cmd.Wait()

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		errorOutput = append(errorOutput, scanner.Text())
	}

	if len(errorOutput) > 0 {
		return errorOutput, fmt.Errorf("TransactionSign error")
	}

	return errorOutput, nil
}

func (c *CreateNFT) TransactionSubmit() (errorOutput []string, err error) {
	cmd := exec.Command("cardano-cli", "transaction", "submit",
		"--tx-file", c.f.GetSignedTransactionFile(),
		"--mainnet")
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return errorOutput, err
	}

	cmd.Wait()

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		errorOutput = append(errorOutput, scanner.Text())
	}

	if len(errorOutput) > 0 {
		return []string{}, fmt.Errorf("TransactionSubmit error")
	}

	return []string{}, nil
}

func (c *CreateNFT) CreateMetadata() error {
	metaData, err := os.Create(c.f.GetMetadataAttrFile()) //TODO add new files for nft policy
	if err != nil {
		log.Println(err)
		return err
	}
	defer metaData.Close()

	policyIdBytes, err := ioutil.ReadFile(c.f.GetPolicyIDFile())
	if err != nil {
		log.Println(err)
		return err
	}

	metaData.WriteString("{\n")
	metaData.WriteString("  \"721\": {")
	metaData.WriteString("  \"" + string(policyIdBytes) + "\": {")
	metaData.WriteString("  \"image\":" + "[\"https://ipfs.io/ipfs/\", \"" + c.processParams.Nft.ImageIPFSHash + "\"],")
	metaData.WriteString("  \"mediaType\":\"" + c.processParams.Nft.MediaType + "\",")
	metaData.WriteString("  \"description\":\"" + c.processParams.Nft.Description + "\"")
	metaData.WriteString("  }")
	metaData.WriteString("  },")
	metaData.WriteString("  \"version\":\"1.0\"")
	metaData.WriteString("  }")
	metaData.WriteString("  }")

	return nil
}
