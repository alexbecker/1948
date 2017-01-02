package auth

import (
	"github.com/DATA-DOG/go-sqlmock"
)

func InitMock(saltString, secretKeyString string) sqlmock.Sqlmock {
	var mock sqlmock.Sqlmock
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		panic(err)
	}

	salt = []byte(saltString)
	secretKey = []byte(secretKeyString)
	return mock
}
