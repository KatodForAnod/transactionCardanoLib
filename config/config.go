package config

import (
	"github.com/bykovme/goconfig"
	"log"
)

type Config struct {
	ID                  string  `json:"id"` //1
	Token               []Token `json:"token"`
	PaymentAddress      string  `json:"payment_address"`         //2
	PaymentSKeyFilePath string  `json:"payment_s_key_file_path"` //3
	PaymentVKeyFilePath string  `json:"payment_v_key_file_path"` //4

	SlotNft  string `json:"slot_nft"`  //5
	IpfsHash string `json:"ipfs_hash"` //6

	UsingExistingPolicy        bool   `json:"using_existing_policy"`
	PolicyIDFile               string `json:"policy_id"`
	PolicyScriptFilePath       string `json:"policy_script_file_path"`
	PolicySigningFilePath      string `json:"policy_signing_file_path"`
	PolicyVerificationFilePath string `json:"policy_verification_file_path"`

	MetadataAttrFile string `json:"metadata_attr_file"`
}

type Token struct {
	TokenName   string `json:"token_name"`
	TokenAmount string `json:"token_amount"`
}

const cConfigPath = "conf.config"

func LoadConfig() (loadedConfig Config, err error) {
	log.Println("Start loading config...")
	/*usrHomePath, err := goconfig.GetUserHomePath()
	if err != nil {
		log.Println(err)
		return loadedConfig, err
	}*/

	err = goconfig.LoadConfig(cConfigPath, &loadedConfig)
	return loadedConfig, err
}
