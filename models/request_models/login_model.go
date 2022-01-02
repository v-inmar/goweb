package request_models

import (
	"encoding/json"
	"errors"
	"strings"
)

type LoginModel struct {
	EmailAddress string `json:"email address"`
	Password string `json:"password"`
}

func (m *LoginModel) HasValidValues()error{
	if len(m.EmailAddress) < 1{
		return errors.New("email address is empty")
	}

	if len(m.Password) < 1{
		return errors.New("password is empty")
	}

	if !strings.Contains(m.EmailAddress, "@") || !strings.Contains(m.EmailAddress, "."){
		return errors.New("invalid email address format")
	}

	return nil
}

func (m *LoginModel) ToStringNoPassword() (string, error){

	// a copy is made so the actual model's value is untouched
	// and this copy can then be used to marshal to a string
	modelCopy := make(map[string]string)
	modelCopy["email address"] = m.EmailAddress
	modelCopy["password"] = ""

	modelByteArr, err := json.Marshal(modelCopy)
	if err != nil{
		return "", err
	}
	return string(modelByteArr), nil
}