package dto

type TokenStruct struct {
	TokenName                  string `json:"token_name"`
	TokenAmount                int64  `json:"token_amount"`
	PaymentAddress             string `json:"payment_address"`
	UsingExistingPolicy        bool   `json:"using_existing_policy"`
	PolicyScriptFilePath       string `json:"policy_script_file_path"`
	PolicySigningFilePath      string `json:"policy_signing_file_path"`
	PolicyVerificationFilePath string `json:"policy_verification_file_path"`
}
