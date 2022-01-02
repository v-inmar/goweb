package response_models

import "net/http"

type StatusModel struct {
	Code int64 `json:"code"`
	Status string `json:"status"`
	Message string `json:"message"`
	Payload string `json:"payload"`
}

/*************/
/* 200 - 299 */
/*************/

// 200 OK
func (m *StatusModel) OK(message, payload string) {
	m.Code = http.StatusOK
	m.Status = "OK"
	m.Message = message
	m.Payload = payload
}

// 201 Created
func (m *StatusModel) Created(message, payload string) {
	m.Code = http.StatusCreated
	m.Status = "Created"
	m.Message = message
	m.Payload = payload
}

/*************/
/* 400 - 499 */
/*************/

// 400 - Bad Request
func (m *StatusModel) BadRequest(message, payload string){
	m.Code = http.StatusBadRequest
	m.Status = "Bad Request"
	m.Message = message
	m.Payload = payload
}

// 401 - Unauthorized
func (m *StatusModel) Unauthorized(message, payload string){
	m.Code = http.StatusUnauthorized
	m.Status = "Unauthorized"
	m.Message = message
	m.Payload = payload
}

// 404 - Not Found
func (m *StatusModel) NotFound(message, payload string){
	m.Code = http.StatusNotFound
	m.Status = "Not Found"
	m.Message = message
	m.Payload = payload
}

// 409 - Conflict
func (m *StatusModel) Conflict(message, payload string){
	m.Code = http.StatusConflict
	m.Status = "Conflict"
	m.Message = message
	m.Payload = payload
}

/*************/
/* 500 - 599 */
/*************/

// 500 - Server Error
func (m *StatusModel) ServerError(message, payload string){
	m.Code = http.StatusInternalServerError
	m.Status = "Server Error"
	m.Message = message
	m.Payload = payload
}