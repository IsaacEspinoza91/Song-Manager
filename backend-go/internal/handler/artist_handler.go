package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/IsaacEspinoza91/Song-Manager/internal/domain"
)

type ArtistHandler struct {
	service domain.ArtistService
}

func NewArtistHandler(service domain.ArtistService) *ArtistHandler {
	return &ArtistHandler{service: service}
}

// CREATE (POST /artists)
func (h *ArtistHandler) Create(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar el JSON entrante
	var input domain.ArtistInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// 2. Llamar a la capa de servicio
	// Contexto (r.Context()) viaja desde aquí hasta la base de datos
	artist, err := h.service.Create(r.Context(), &input)
	if err != nil {
		// Si el servicio devuelve un error (ej. validación fallida), responder 400
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 3. Responder con éxito (201 Created) y el objeto creado en JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(artist)
}

// GET ALL (GET /artists)
func (h *ArtistHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	artists, err := h.service.GetAll(r.Context())
	if err != nil {
		// Si falla la base de datos, es un error interno del servidor (500)
		http.Error(w, "Error interno obteniendo artistas", http.StatusInternalServerError)
		return
	}

	// Responder con éxito (200 OK)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(artists)
}

// GET ID (GET /artists/{id})
func (h *ArtistHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "El ID de la URL debe ser un número válido", http.StatusBadRequest)
		return
	}

	artist, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if err.Error() == "artista no encontrado" {
			http.Error(w, err.Error(), http.StatusNotFound) // 404
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(artist)
}

// UPDATE (PUT /artists/{id})
func (h *ArtistHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Extraer ID, definido en router {id}
	idString := r.PathValue("id")

	// Convertir string a int64 (base 10, 64 bits)
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "El ID de la URL debe ser un número válido", http.StatusBadRequest)
		return
	}

	// Decodificar el JSON entrante al DTO (ArtistInput)
	var input domain.ArtistInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
		return
	}

	// Pasar el ID y los datos a la capa de Servicio
	artist, err := h.service.Update(r.Context(), id, &input)
	if err != nil {
		// Caso error fue porque no se encontró el artista
		if err.Error() == "artista no encontrado" {
			http.Error(w, err.Error(), http.StatusNotFound) // 404
			return
		}
		// Cualquier otro error de validación o base de datos
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK
	json.NewEncoder(w).Encode(artist)
}

// DELETE (DELETE /artists/{id})
func (h *ArtistHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Extraer y convertir el ID
	idString := r.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "El ID de la URL debe ser un número válido", http.StatusBadRequest)
		return
	}

	// Llamar al servicio. Soft Delete
	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if err.Error() == "artista no encontrado" {
			http.Error(w, err.Error(), http.StatusNotFound) // 404
			return
		}
		http.Error(w, "Error al eliminar el artista", http.StatusInternalServerError)
		return
	}

	// Responder al frontend indicando exito, sin contenido en el body
	// 204 No Content es el código HTTP estándar para un DELETE exitoso
	w.WriteHeader(http.StatusNoContent)
}
