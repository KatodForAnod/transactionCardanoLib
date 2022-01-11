package policy

const (
	PolicyDirName             = "policy"
	RawTransactionFile        = "matx.raw"                      // Raw transaction to mint token
	SignedTransactionFile     = "matx.signed"                   // Signed transaction to mint token
	MetadataAttrFile          = "metadata.json"                 // Metadata to specify NFT attributes
	PaymentAddrFile           = "payment.addr"                  // Address to send/receive
	PaymentSignKeyFile        = "payment.skey"                  // Payment signing key
	PaymentVerifyKeyFile      = "payment.vkey"                  // Payment verification key
	PolicyScriptFile          = PolicyDirName + "/policy.scipt" // Script to generate the policyID
	PolicySigningKeyFile      = PolicyDirName + "/policy.skey"  // Policy signing key
	PolicyVerificationkeyFile = PolicyDirName + "/policy.vkey"  // Policy verification key
	PolicyIDFile              = PolicyDirName + "/policyID"     // File which holds the policy ID
	ProtocolParametersFile    = "protocol.json"                 // Protocol parameters
)
