package models

// TODO: Refactor to separate files

type ResponseBodyModel struct {
	PID   string `json:"pid"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}

type ResponseBodyStatusModel struct {
	Code int `json:"code"`
	Status string `json:"status"`
	Message string `json:"message"`
}

func (m *ResponseBodyStatusModel) OK(message string){
	m.Code = 200
	m.Status = "OK"
	m.Message = message
}

type ResponseBodyErrorMessageModel struct {
	Code int `json:"code"`
	Status string `json:"status"`
	Message string `json:"message"`
}

func (m *ResponseBodyErrorMessageModel) InternalServerError(message string) {
	m.Code = 500
	m.Status = "Server Error"
	m.Message = message
}

func (m *ResponseBodyErrorMessageModel) BadRequest(message string) {
	m.Code = 400
	m.Status = "Bad Request"
	m.Message = message
}

func (m *ResponseBodyErrorMessageModel) Conflict(message string) {
	m.Code = 409
	m.Status = "Conflict"
	m.Message = message
}

type ResponseBodyJWTModel struct {
	Code int `json:"code"`
	Status string `json:"status"`
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}