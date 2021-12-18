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
	gen, err := uuid.NewRandom()
	if err != nil{
		return err
	}
	val := gen.String()
	if val == ""{
		return errors.New("empty string for GenerateForUPID")
	}

	subVal := strings.Replace(val,"-", "", -1)[0:8]
	s.Value = subVal
	return nil
}
