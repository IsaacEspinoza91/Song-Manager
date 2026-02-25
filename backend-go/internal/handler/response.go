package handler

import (
	"encoding/json"
	"net/http"
)

// estructura estandar para que el frontend maneje errores
type APIError struct {
	// Código HTTP (ej. 400, 404, 500)
	Status int `json:"status"`
	// Mensaje amigable para el usuario
	Message string `json:"message"`
	// Details contiene información técnica o validaciones específicas (opcional)
	Details interface{} `json:"details,omitempty"`
}

type InfoResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// WriteError es un helper para estandarizar la respuesta de errores
func WriteError(w http.ResponseWriter, status int, message string, details interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errResp := APIError{
		Status:  status,
		Message: message,
		Details: details,
	}

	json.NewEncoder(w).Encode(errResp)
}

// WriteJSON es un helper para enviar respuestas exitosas
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Helper para respuestas con codigo y mensaje
func WriteMessageJSON(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := InfoResponse{
		Status:  status,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

// WriteNoContent se usa específicamente para respuestas 204 (como un DELETE exitoso)
// Solo devolver codigo 204 No Content, el cuerpo debe estar vacio
func WriteNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
