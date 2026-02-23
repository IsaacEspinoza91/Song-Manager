package handler

import (
	"encoding/json"
	"net/http"
)

// estructura estandar para que el frontend maneje errores
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"` // omitempty hace que no aparezca en el JSON si está vacío
}

// WriteError es un helper para enviar respuestas de error en formato JSON
func WriteJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errResp := ErrorResponse{
		Error: message,
	}

	json.NewEncoder(w).Encode(errResp)
}

// WriteJSON es un helper para enviar respuestas exitosas
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Si data es nil (ej. en un DELETE), no codificamos nada
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
