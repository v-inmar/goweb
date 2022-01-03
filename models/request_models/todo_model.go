package request_models

import (
	"encoding/json"
	"errors"
)

type TodoModel struct {
	Title string `json:"title"`
	Body string `json:"body"`
}

func (m *TodoModel) HasValidValues()error{
	if len(m.Title) < 1 || len(m.Title) > 128 {
		return errors.New("title must be between 1 and 128 characters long")
	}
	return nil
}

func (m *TodoModel) Stringify()(string, error){
	modelByteArr, err := json.Marshal(m)
	if err != nil{
		return "", err
	}
	return string(modelByteArr), nil
}