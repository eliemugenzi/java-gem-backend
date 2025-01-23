package utils

import "golang.org/x/crypto/bcrypt"


func HashPassword(password [] byte) string {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)

	if err != nil {
		panic(err)
	}
	return string(hash)
}

func ComparePassword(hashedPassword string, password []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), password)
	return err == nil
	}