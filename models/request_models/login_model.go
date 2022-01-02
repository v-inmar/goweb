package request_models

type LoginModel struct {
	EmailAddress string `json:"email address"`
	Password string `json:"password"`
}