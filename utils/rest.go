package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RandySteven/onboard-be/entities/payloads/responses"
)

func ContentType(w http.ResponseWriter, contentType string) {
	w.Header().Set("Content-Type", contentType)
}

func BindJSON(req *http.Request, request interface{}) error {
	return json.NewDecoder(req.Body).Decode(request)
}

func ResponseHandler(w http.ResponseWriter, responseCode int, message string, dataKey *string, responseData any, err error) {
	ContentType(w, "application/json")
	w.WriteHeader(responseCode)
	responseMap := make(map[string]any)
	if dataKey != nil && responseData != nil {
		responseMap[*dataKey] = responseData
	}
	response := responses.NewResponse(message, responseMap, err)
	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		log.Println("failed to encode response", encodeErr)
	}
}
