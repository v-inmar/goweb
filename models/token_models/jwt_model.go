package token_models

import "encoding/json"

type JWTModel struct {
	AccessToken string `json:"X-Access-Token"`
	RefreshToken string `json:"X-Refresh-Token"`
}

func (m *JWTModel) Stringify()(string, error){
	modelByteArr, err := json.Marshal(m)
	if err != nil{
		return "", err
	}
	return string(modelByteArr), nil
}