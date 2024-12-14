package config

type JWT struct {
	SigningKey  string `yaml:"signing_key"`
	TokenExpiry string `yaml:"token_expiry"`
}
