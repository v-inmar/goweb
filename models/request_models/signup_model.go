package request_models

type SignupModel struct {
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
	EmailAddress string `json:"email address"`
	Password string `json:"password"`
}