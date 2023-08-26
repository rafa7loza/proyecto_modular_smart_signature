package auth

import (
	"os"
	"log"

	"testing"
	"deviloza.com.mx/misc"
)

func TestDBConn(t *testing.T) {
	baseDir := os.Getenv("MICRO_SERVICES_MODS")
	secretsPath := baseDir + "/secrets.json"
	s, err := misc.NewSecrets(secretsPath)
	if err != nil {
		t.Fatalf(`Failed to parse secrets file from %s\n(err)%v`, 
			secretsPath, 
			err,
		)
	}

	log.Println(s)
	dsn := s.GenerateDSN()
	log.Println(dsn)
	db := Conn(dsn)
	log.Println(db)
}