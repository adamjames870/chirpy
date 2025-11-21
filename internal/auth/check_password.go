package auth

import (
	"errors"

	"github.com/alexedwards/argon2id"
)

func CheckPasswordHash(password, hash string) (bool, error) {
	match, errMatch := argon2id.ComparePasswordAndHash(password, hash)
	if errMatch != nil {
		return false, errors.New("error creating hash: " + errMatch.Error())
	}
	return match, nil
}
