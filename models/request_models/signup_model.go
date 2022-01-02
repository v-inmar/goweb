package request_models

import (
	"encoding/json"
	"errors"
	"strings"
)

type SignupModel struct {
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
	EmailAddress string `json:"email address"`
	Password string `json:"password"`
}

func (m *SignupModel) HasValidValues()error{
	if len(m.Firstname) < 1 || len(m.Firstname) > 16{
		return errors.New("firstname must be between 1 and 16 characters long")
	}

	if len(m.Lastname) < 1 || len(m.Lastname) > 16{
		return errors.New("lastname must be between 1 and 16 characters long")
	}

	if len(m.EmailAddress) < 5 || len(m.EmailAddress) > 128{
		return errors.New("email address must be between 5 and 128 characters long")
	}

	if !strings.Contains(m.EmailAddress, "@"){
		return errors.New("email address must be a valid email address")
	}

	if !strings.Contains(m.EmailAddress, "."){
		return errors.New("email address must be a valid email address")
	}

	if len(m.Password) < 8{
		return errors.New("password must be atleast 8 characters long")
	}

	return nil
}

func (m *SignupModel) ToStringNoPassword() (string, error){
	m.Password = ""
	modelByteArr, err := json.Marshal(m)
	if err != nil{
		return "", err
	}
	return string(modelByteArr), nil
}