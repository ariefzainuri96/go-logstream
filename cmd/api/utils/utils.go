package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/response"	
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	// 1. Marshal the data (Fail Fast on marshalling error)
	respBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("ERROR: Failed to marshal response data: %v", err)
		// Fallback to plain text 500 error if marshalling fails
		RespondError(w, http.StatusInternalServerError, "Internal Server Error: Failed to serialize response.")		
		return
	}

	// 2. Set headers and write status/body
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(respBytes)
}

// respondError writes a standardized BaseResponse for errors.
// This is the key method to eliminate your boilerplate.
func RespondError(w http.ResponseWriter, status int, message string) {
	// Construct the standardized error response body
	errorResp := response.BaseResponse{
		Status:  int64(status),
		Message: message,
	}
	// Use writeJSON internally
	WriteJSON(w, status, errorResp)
}