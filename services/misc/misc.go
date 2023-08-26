package misc

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type DBCreds struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Schema   string `json:"schema"`
	Port     string `json:"port"`
	Address  string `json:"ipAddress"`
}

type Secrets struct {
	ApiDB          DBCreds `json:"db"`
	SendGridAPIKey string  `json:"SENDGRID_API_KEY"`
}

func ReadJsonFile(filePath string, obj *Secrets) error {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	bytes, err := io.ReadAll(jsonFile)

	if err != nil {
		return err
	}

	json.Unmarshal(bytes, &obj)

	return nil
}

func NewSecrets(secretsPath string) (*Secrets, error) {
	var secrets Secrets
	err := ReadJsonFile(secretsPath, &secrets)
	if err != nil {
		return nil, err
	}
	return &secrets, nil
}

func (secrets *Secrets) GenerateDSN() string {
	DEFAULT_OPTS := "charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?%s",
		secrets.ApiDB.User,
		secrets.ApiDB.Password,
		secrets.ApiDB.Address,
		secrets.ApiDB.Port,
		secrets.ApiDB.Schema,
		DEFAULT_OPTS,
	)
	return dsn
}
