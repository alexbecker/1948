package auth

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const minKeyLen = 16

var (
	db        *sql.DB
	salt      []byte
	secretKey []byte
)

func Init(dbName, saltString, secretKeyString string) {
	var err error
	db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}

	salt = []byte(saltString)
	secretKey = []byte(secretKeyString)
	if len(secretKey) < minKeyLen {
		panic("Secret key absent or too short!")
	}
}
