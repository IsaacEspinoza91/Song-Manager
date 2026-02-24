package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/IsaacEspinoza91/Song-Manager/internal/domain"
)

type SongHandler struct {
	service domain.SongService
}

func NewSongHandler(service domain.SongService) *SongHandler {
	return &SongHandler{service: service}
}

// CREATE con o sin artistas (POST /songs)
func (h *SongHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Decodificar JSON entrante
	var input domain.SongInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteError(w, http.StatusBadRequest, "Formato JSON inválido", err.Error())
		return
	}

	song, err := h.service.Create(r.Context(), &input)
	if err != nil {
		var valErrs domain.ValidationError
		if errors.As(err, &valErrs) {
			WriteError(w, http.StatusBadRequest, "Datos de entrada inválidos", valErrs)
			return
		}

		// if caso id artista no existe code 400 bad request

		WriteError(w, http.StatusInternalServerError, "No se pudo crear la cancion", err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, song) // 201
}

// GET ID (GET /songs/{id})
func (h *SongHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 54)
	if err != nil {
		log.Printf("[ERROR en Handler] %v\n", err)
		WriteError(w, http.StatusBadRequest, "El ID de la URL debe ser un número entero válido mayor a 0", nil)
		return
	}
	if id <= 0 {
		WriteError(w, http.StatusBadRequest, "El ID debe ser mayor a 0", nil)
		return
	}

	song, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrSongNotFound) {
			WriteError(w, http.StatusNotFound, err.Error(), nil) // 404
			return
		}
		log.Printf("[ERROR INTERNO en Handler] %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Error al buscar la cancion", nil) // 500
		return
	}

	WriteJSON(w, http.StatusOK, song) // 200
}

// GET ALL (GET /songs/all)
func (h *SongHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	songs, err := h.service.GetAll(r.Context())
	if err != nil {
		log.Printf("[ERROR INTERNO en Handler] %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Error interno obteniendo canciones", nil)
	}
	// Slice vacio sino hay canciones
	if songs == nil {
		songs = []domain.Song{}
	}

	WriteJSON(w, http.StatusOK, songs) // 200
}

// GET ALL PAG (GET /songs?page=1&limit=10&artist_id=1&artist_name=shakira&name=sordo)
func (h *SongHandler) GetAllPaginated(w http.ResponseWriter, r *http.Request) {
	// Extraer query params
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10
	}
	pagination := domain.PaginationParams{
		Page:  page,
		Limit: limit,
	}

	// Extraemos el ID numérico si viene en la query
	var artistaID int64
	if artistIDStr := r.URL.Query().Get("artist_id"); artistIDStr != "" {
		artistaID, _ = strconv.ParseInt(artistIDStr, 10, 64)
	}
	filter := domain.SongFilter{
		Title:      r.URL.Query().Get("title"),
		ArtistID:   artistaID,
		ArtistName: r.URL.Query().Get("artist_name"),
	}

	paginatedData, err := h.service.GetAllPaginated(r.Context(), filter, pagination)
	if err != nil {
		log.Printf("[ERROR INTERNO en Handler] %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Error obteniendo la lista de canciones", nil)
		return
	}

	WriteJSON(w, http.StatusOK, paginatedData) // 200
}

// UPDATE (PUT /artist/{id})
func (h *SongHandler) Update(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 54)
	if err != nil {
		log.Printf("[ERROR en Handler] %v\n", err)
		WriteError(w, http.StatusBadRequest, "El ID de la URL debe ser un número entero válido mayor a 0", nil)
		return
	}
	if id <= 0 {
		WriteError(w, http.StatusBadRequest, "El ID debe ser mayor a 0", nil)
		return
	}

	var input domain.SongInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteError(w, http.StatusBadRequest, "Formato JSON inválido", err.Error())
		return
	}

	song, err := h.service.Update(r.Context(), id, &input)
	if err != nil {
		var valErrs domain.ValidationError
		if errors.As(err, &valErrs) {
			WriteError(w, http.StatusBadRequest, "Datos de actualización inválidos", valErrs)
			return
		}

		if errors.Is(err, domain.ErrSongNotFound) {
			WriteError(w, http.StatusNotFound, err.Error(), nil) // 404
			return
		}

		// if caso id artista no existe code 400 bad request

		log.Printf("[ERROR INTERNO en Handler] %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Error actualizando la cancion", nil) // 500
		return
	}

	WriteJSON(w, http.StatusOK, song)
}

// DELETE (DELETE /artist/{id})
func (h *SongHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 54)
	if err != nil {
		log.Printf("[ERROR en Handler] %v\n", err)
		WriteError(w, http.StatusBadRequest, "El ID de la URL debe ser un número entero válido mayor a 0", nil)
		return
	}
	if id <= 0 {
		WriteError(w, http.StatusBadRequest, "El ID debe ser mayor a 0", nil)
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrSongNotFound){
			WriteError(w, http.StatusNotFound, err.Error(), nil) // 404
			return
		}
		log.Printf("[ERROR INTERNO en Handler] %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Error al eliminar la cancion", nil) // 500
		return
	}

	WriteNoContent(w)
}
