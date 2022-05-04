package files

import (
	"transactionCardanoLib/config"
)

/*var (
	PolicyDirName         = "policy"
	RawTransactionFile    = "matx.raw"    // Raw transaction to mint token
	SignedTransactionFile = "matx.signed" // Signed transaction to mint token

	RawTransactionSendTokenFile    = "rec_matx.raw"
	SignedTransactionSendTokenFile = "rec_matx.signed"

	MetadataAttrFile = "metadata.json" // Metadata to specify NFT attributes
	PaymentAddrFile  = "payment.addr"  // Address to send/receive

	PaymentSignKeyFile        = "payment.skey"                   // Payment signing key
	PaymentVerifyKeyFile      = "payment.vkey"                   // Payment verification key
	PolicyScriptFile          = PolicyDirName + "/policy.script" // Script to generate the policyID
	PolicySigningKeyFile      = PolicyDirName + "/policy.skey"   // Policy signing key
	PolicyVerificationkeyFile = PolicyDirName + "/policy.vkey"   // Policy verification key
	PolicyIDFile              = PolicyDirName + "/policyID"      // File which holds the policy ID
	ProtocolParametersFile    = "protocol.json"                  // Protocol parameters
)*/

type Files struct {
	policyDirName         string
	rawTransactionFile    string
	signedTransactionFile string

	rawTransactionSendTokenFile    string
	signedTransactionSendTokenFile string

	metadataAttrFile string
	//PaymentAddrFile  string

	paymentSignKeyFile        string // Payment signing key
	paymentVerifyKeyFile      string // Payment verification key
	policyScriptFile          string // Script to generate the policyID
	policySigningKeyFile      string // Policy signing key
	policyVerificationkeyFile string // Policy verification key
	protocolParametersFile    string // Protocol parameters

	policyIDFile string
}

func (f *Files) Init(config config.Config) {
	f.policyDirName = "policy"
	f.rawTransactionFile = "matx.raw"       // Raw transaction to mint token
	f.signedTransactionFile = "matx.signed" // Signed transaction to mint token

	f.rawTransactionSendTokenFile = "rec_matx.raw"
	f.signedTransactionSendTokenFile = "rec_matx.signed"

	//f.PaymentAddrFile  = "payment.addr"  // Address to send/receive

	f.paymentSignKeyFile = config.PaymentSKeyFilePath
	f.paymentVerifyKeyFile = config.PaymentVKeyFilePath

	if config.UsingExistingPolicy {
		f.policyScriptFile = config.PolicyScriptFilePath
		f.policySigningKeyFile = config.PolicySigningFilePath
		f.policyVerificationkeyFile = config.PolicyVerificationFilePath
		f.policyIDFile = config.PolicyID // from config should load file

		f.metadataAttrFile = config.MetadataAttrFile
	} else {
		f.policyScriptFile = f.policyDirName + "/policy.script"        // Script to generate the policyID
		f.policySigningKeyFile = f.policyDirName + "/policy.skey"      // Policy signing key
		f.policyVerificationkeyFile = f.policyDirName + "/policy.vkey" // Policy verification key
		f.policyIDFile = f.policyDirName + "/policyID"

		f.metadataAttrFile = "metadata.json" // Metadata to specify NFT attributes
	}

	f.protocolParametersFile = "protocol.json" // Protocol parameters
}

func (p *Files) GetPolicyDirName() string {
	return p.policyDirName
}
func (p *Files) GetRawTransactionFile() string {
	return p.rawTransactionFile
}
func (p *Files) GetSignedTransactionFile() string {
	return p.signedTransactionFile
}
func (p *Files) GetRawTransactionSendTokenFile() string {
	return p.rawTransactionSendTokenFile
}
func (p *Files) GetSignedTransactionSendTokenFile() string {
	return p.signedTransactionSendTokenFile
}
func (p *Files) GetMetadataAttrFile() string {
	return p.metadataAttrFile
}
func (p *Files) GetPaymentSignKeyFile() string {
	return p.paymentSignKeyFile
}
func (p *Files) GetPaymentVerifyKeyFile() string {
	return p.paymentVerifyKeyFile
}
func (p *Files) GetPolicyScriptFile() string {
	return p.policyScriptFile
}
func (p *Files) GetPolicySigningKeyFile() string {
	return p.policySigningKeyFile
}
func (p *Files) GetPolicyVerificationkeyFile() string {
	return p.policyVerificationkeyFile
}
func (p *Files) GetPolicyIDFile() string {
	return p.policyIDFile
}

func (p *Files) GetProtocolParametersFile() string {
	return p.protocolParametersFile
}
