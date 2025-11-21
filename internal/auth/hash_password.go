package auth

import (
	"errors"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	pwdHash, errHash := argon2id.CreateHash(password, argon2id.DefaultParams)
	if errHash != nil {
		return "", errors.New("error creating hash: " + errHash.Error())
	}
	return pwdHash, nil
}
