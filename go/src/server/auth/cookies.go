package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"strings"
)

const minKeyLen = 16

var secretKey []byte

func getSecretKey() []byte {
	if secretKey != nil {
		return secretKey
	}
	secretKey = []byte(os.Getenv("SECRET"))
	if len(secretKey) < minKeyLen {
		panic("SECRET absent or too short!")
	}
	return secretKey
}

func SignCookie(cookie *http.Cookie) {
	sigBytes := hmac.New(sha256.New, getSecretKey()).Sum([]byte(cookie.Value))
	sig := base64.StdEncoding.EncodeToString(sigBytes)
	cookie.Value += "." + sig
}

func DecodeCookie(cookie *http.Cookie) (string, error) {
	split := strings.SplitN(cookie.Value, ".", 2)
	if len(split) != 2 {
		return "", errors.New("Not a valid signed cookie")
	}
	value := split[0]
	sigBase64 := split[1]

	sig, err := base64.StdEncoding.DecodeString(sigBase64)
	if err != nil {
		return value, err
	}

	sig2 := hmac.New(sha256.New, getSecretKey()).Sum([]byte(value))

	if !hmac.Equal(sig, sig2) {
		return value, errors.New("Signature mismatch")
	}

	return value, nil
}

func AuthCookie(username string) http.Cookie {
	cookie := http.Cookie{
		Name:   "auth",
		Value:  username,
		Secure: false,
	}
	SignCookie(&cookie)
	return cookie
}

func SetAuthCookie(w http.ResponseWriter, username string) {
	cookie := AuthCookie(username)
	http.SetCookie(w, &cookie)
}

func CheckAuthCookie(req *http.Request) (string, error) {
	cookie, err := req.Cookie("auth")
	if err != nil {
		return "", err
	}

	username, err := DecodeCookie(cookie)
	return username, err
}
