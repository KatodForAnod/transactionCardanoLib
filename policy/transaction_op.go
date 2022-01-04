package policy

import (
	"fmt"
	"log"
	"os/exec"
	"transactionCardanoLib/dto"
)

const transactionSignTmpl = "cardano-cli transaction sign " +
	"--signing-key-file payment.skey  " +
	"--signing-key-file %s " +
	"--testnet-magic %s " +
	"--tx-body-file matx.raw  " +
	"--out-file matx.signed"

func TransactionSign(id string, tokenStruct dto.TokenStruct) error {
	comm := fmt.Sprintf(transactionSignTmpl, tokenStruct.PolicySigningFilePath, id)

	err := exec.Command(comm).Run()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

const transactionBuildTmpl = "cardano-cli transaction build-raw " +
	"--fee %s " +
	"--tx-in $%s#$%s " +
	"--tx-out $%s+$%s+\"$%s $%s + $%s $%s\" " +
	"--mint=\"$%s $%s + $%s $%s\" " +
	"--minting-script-file %s " +
	"--out-file matx.raw"

func TransactionBuild(fee, txHash, txIx, address, output, tokenAmount,
	tokenName1, tokenName2, policyScriptFilePath string) error {
	comm := fmt.Sprintf(transactionBuildTmpl, fee, txHash, txIx, address, output, tokenAmount, tokenName1,
		tokenAmount, tokenName2, tokenAmount, tokenName1, tokenAmount, tokenName2, policyScriptFilePath)

	err := exec.Command(comm).Run()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

const transactionPreBuildTmpl = "cardano-cli query utxo --address %s --testnet-magic %s"

func TransactionPreBuild(address, id string) (cliOutput string, err error) {
	comm := fmt.Sprintf(transactionPreBuildTmpl, address, id)

	out, err := exec.Command(comm).Output()
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(out), nil
}
