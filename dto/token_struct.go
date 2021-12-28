package dto

import "os"

type TokenStruct struct {
	TokenName              string   `json:"token_name"`
	TokenAmount            int64    `json:"token_amount"`
	PaymentAddress         string   `json:"payment_address"`
	UsingExistingPolicy    bool     `json:"using_existing_policy"`
	PolicyScriptFile       *os.File `json:"policy_script_file"`
	PolicySigningFile      *os.File `json:"policy_signing_file"`
	PolicyVerificationFile *os.File `json:"policy_verification_file"`
}
