package internal

import (
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

var PORT string
var DB_PORT string
var DB_USERNAME string
var DB_PASSWORD string
var DB_NAME string

var defaultValues = map[string]string{"PORT": "8080", "DB_PORT": "6543", "DB_USERNAME": "user", "DB_PASSWORD": "pass", "DB_NAME": "db"}

func Port() string {
	return PORT
}

func Db_port() string {
	return DB_PORT
}

func Db_username() string {
	return DB_USERNAME
}

func Db_password() string {
	return DB_PASSWORD
}

func Db_name() string {
	return DB_NAME
}

type conf struct {
	Port        string `yaml:"port"`
	Db_port     string `yaml:"db_port"`
	Db_username string `yaml:"db_username"`
	Db_password string `yaml:"db_password"`
	Db_name     string `yaml:"db_name"`
}

func (c *conf) GetFiledValue(fieldName string) string {
	switch fieldName {
	case "port":
		return c.Port
	case "db_port":
		return c.Db_port
	case "db_username":
		return c.Db_username
	case "db_password":
		return c.Db_password
	case "db_name":
		return c.Db_name
	default:
		return ""
	}
}

func init() {
	var condidates = make(map[string]string)
	condidates["PORT"] = os.Getenv("MICROSERVICE_PORT")
	condidates["DB_PORT"] = os.Getenv("DB_CONNECTION_PORT")
	condidates["DB_USERNAME"] = os.Getenv("DB_CONNECTION_USERNAME")
	condidates["DB_PASSWORD"] = os.Getenv("DB_CONNECTION_PASSWORD")
	condidates["DB_NAME"] = os.Getenv("DB_CONNECTION_NAME")

	if configYamlFile, err := os.ReadFile("configs/config.yaml"); err == nil {
		var configYaml conf

		if unmarshErr := yaml.Unmarshal(configYamlFile, &configYaml); unmarshErr == nil {
			for key, value := range condidates {
				if value == "" {
					condidates[key] = configYaml.GetFiledValue(strings.ToLower(key))
				}
			}
		}
	}

	for key, value := range condidates {
		if value == "" {
			condidates[key] = defaultValues[key]
		}
	}

	PORT = condidates["PORT"]
	DB_PORT = condidates["DB_PORT"]
	DB_USERNAME = condidates["DB_USERNAME"]
	DB_PASSWORD = condidates["DB_PASSWORD"]
	DB_NAME = condidates["DB_NAME"]
}
