package models

type APIResponse struct {
	Code     int         `json:"code"`
	Message  interface{} `json:"message"`
	HttpCode int         `json:"httpCode"`
}
