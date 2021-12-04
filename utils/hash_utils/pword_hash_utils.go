package hash_utils

import "golang.org/x/crypto/bcrypt"


func PasswordHash(password string) (string, error){
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil{
		return "", err
	}
	return string(hashedBytes), nil
}

func PasswordCheck(hashed string, password string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
