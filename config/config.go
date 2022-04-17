package config

import (
	"github.com/bykovme/goconfig"
	"log"
)

type Config struct {
	ID                         string  `json:"id"`
	Token                      []Token `json:"token"`
	PaymentAddress             string  `json:"payment_address"`
	UsingExistingPolicy        bool    `json:"using_existing_policy"`
	PolicyID                   string  `json:"policy_id"`
	PolicyScriptFilePath       string  `json:"policy_script_file_path"`
	PolicySigningFilePath      string  `json:"policy_signing_file_path"`
	PolicyVerificationFilePath string  `json:"policy_verification_file_path"`
}

type Token struct {
	TokenName   string `json:"token_name"`
	TokenAmount int64  `json:"token_amount"`
}

const cConfigPath = "conf.config"

func LoadConfig() (loadedConfig Config, err error) {
	log.Println("Start loading config...")
	usrHomePath, err := goconfig.GetUserHomePath()
	if err != nil {
		log.Println(err)
		return loadedConfig, err
	}

	err = goconfig.LoadConfig(cConfigPath, &loadedConfig)
	if err == nil {
		return loadedConfig, nil
	}

	log.Println("Config", usrHomePath+cConfigPath, "not found")
	return loadedConfig, err
}
