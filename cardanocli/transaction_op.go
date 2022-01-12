package cardanocli

import (
	"fmt"
	"log"
	"os/exec"
	"transactionCardanoLib/config"
)

const transactionSignTmpl = "cardano-cli transaction sign " +
	"--signing-key-file payment.skey  " +
	"--signing-key-file %s " +
	"--testnet-magic %s " +
	"--tx-body-file " + RawTransactionFile + " " +
	"--out-file " + SignedTransactionFile

func TransactionSign(id string, token config.TokenStruct) error {
	//comm := fmt.Sprintf(transactionSignTmpl, token.PolicySigningFilePath, id)

	err := exec.Command("cardano-cli", "transaction", "sign", "--signing-key-file",
		PaymentSignKeyFile, "-signing-key-file", token.PolicySigningFilePath,
		"--testnet-magic", id, "--tx-body-file", RawTransactionFile, "--out-file", SignedTransactionFile).Run()
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
	"--out-file " + RawTransactionFile

// TransactionBuild - tokenName1 and tokenName2 must be in base16
func TransactionBuild(fee, txHash, txIx, address, output, tokenAmount,
	tokenName1, tokenName2, policyScriptFilePath string) error {
	/*comm := fmt.Sprintf(transactionBuildTmpl, fee, txHash, txIx, address, output, tokenAmount, tokenName1,
	tokenAmount, tokenName2, tokenAmount, tokenName1, tokenAmount, tokenName2, policyScriptFilePath)*/
	txOut := fmt.Sprintf("$%s+$%s+\"$%s $%s + $%s $%s\"", address, output, tokenAmount, tokenName1,
		tokenAmount, tokenName2)
	mint := fmt.Sprintf("$%s $%s + $%s $%s", tokenAmount, tokenName1, tokenAmount, tokenName2)

	err := exec.Command("cardano-cli", "transaction", "build-raw",
		"--fee", fee, "--tx-in", txHash, "--tx-out", txOut, "--mint=", mint,
		"--minting-script-file", policyScriptFilePath,
		"--out-file", RawTransactionFile).Run()
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

/*func base16Encode(input string)(string, error) {

}*/
