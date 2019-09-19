package helpers

import (
	"go-es/constants"
	"net/http"

	"github.com/tnistest/internal/models"
)

var Message = map[int]models.APIResponse{
	constants.APIGeneralSuccess: models.APIResponse{Message: "OK", HttpCode: http.StatusOK},
	constants.APIErrorUnknown:   models.APIResponse{Message: "Terjadi kesalahan pada kudo.", HttpCode: http.StatusInternalServerError},
}

func GetAPIResponse(apiStatusCode int, data interface{}) (resp models.APIResponse) {
	resp.Code = apiStatusCode
	if data == nil {
		resp = Message[apiStatusCode]
		resp.Code = apiStatusCode
	} else {
		resp.Message = data
		resp.HttpCode = http.StatusOK
	}
	return
}
