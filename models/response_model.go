package models

// Response body data model
type ResponseModel struct {
	Status  string       `json:"status"`
	Payload PayloadModel `json:"payload"`
}

// Reponse body payload data model
type PayloadModel struct {
	Payload string `json:"payload"`
}
