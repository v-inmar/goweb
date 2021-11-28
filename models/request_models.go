package models

type RequestBodyModel struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type RequestSignUpBodyModel struct {
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
	EmailAddress string `json:"email address"`
	Password string `json:"password"`
}

type RequestLoginBodyModel struct {
	EmailAddress string `json:"email address"`
	Password string `json:"password"`
}