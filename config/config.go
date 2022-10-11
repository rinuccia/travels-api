package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Port string `yaml:"port"`
	DB   struct {
		Username string `yaml:"username"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		DBName   string `yaml:"db_name"`
	} `yaml:"db"`
}

var instance *Config

func GetConfig() (*Config, error) {
	instance = &Config{}
	err := cleanenv.ReadConfig("./config/config.yml", instance)
	if err != nil {
		return nil, err
	}
	return instance, nil
}
