package random_utils

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type RandomString struct {
	Value string
}


// Creates a rnadom string to be used as upid value
func (s* RandomString) GenerateForUPID() error{
	val, err := generateString()
	if err != nil{
		return err
	}
	subVal := strings.Replace(val,"-", "", -1)[0:8]
	s.Value = subVal
	return nil
}

// Creates a random strings to be used as auth value
func (s* RandomString) GenerateAuth() error{
	val1, err := generateString()
	if err != nil{
		return err
	}

	val2, err := generateString()
	if err != nil{
		return err
	}

	s.Value = (val1+val2)[:64]
	return nil
}

// helps generating random string
func generateString()(string, error){
	gen, err := uuid.NewRandom()
	if err != nil{
		return "", err
	}
	val := gen.String()
	if val == ""{
		return "", errors.New("empty string for Generate")
	}
	return val, nil
}
