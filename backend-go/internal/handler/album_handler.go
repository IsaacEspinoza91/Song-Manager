package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/IsaacEspinoza91/Song-Manager/internal/domain"
)

type AlbumHandler struct {
	service domain.AlbumService
}

func NewAlbumHandler(service domain.AlbumService) *AlbumHandler {
	return &AlbumHandler{service: service}
}

// Create con artistas pero con tracks opcionales (POST /albums)
func (h *AlbumHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.AlbumInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteError(w, http.StatusBadRequest, "Formato JSON inválido", err.Error())
		return
	}

	album, err := h.service.Create(r.Context(), &input)
	if err != nil {
		var valErrs domain.ValidationError
		if errors.As(err, &valErrs) {
			WriteError(w, http.StatusBadRequest, "Datos de entrada inválidos", valErrs)
			return
		}
		// Casos id artista o id cancion no existe

		WriteError(w, http.StatusInternalServerError, "No se pudo crear la cancion", err.Error()) // 500
		return
	}

	WriteJSON(w, http.StatusCreated, album) // 201
}

// GET ID (GET /albums/{id})
func (h *AlbumHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		log.Printf("[ERROR en Handler] %v\n", err)
		WriteError(w, http.StatusBadRequest, "El ID de la URL debe ser un número entero válido mayor a 0", nil)
		return
	}
	if id <= 0 {
		WriteError(w, http.StatusBadRequest, "El ID debe ser mayor a 0", nil)
		return
	}

	album, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrAlbumNotFound) {
			WriteError(w, http.StatusNotFound, err.Error(), nil) // 404
			return
		}
		log.Printf("[ERROR INTERNO en Handler] %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Error al buscar la cancion", nil) // 500
		return
	}

	WriteJSON(w, http.StatusOK, album) // 200
}

// GET ALL PAG (GET /albums?page=1&limit=10&artist_id=1)
func (h *AlbumHandler) GetAllPaginated(w http.ResponseWriter, r *http.Request) {
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

	// Extraer artist_id si es viene en query params
	filter := domain.AlbumFilter{
		Title:      r.URL.Query().Get("title"),
		Type:       r.URL.Query().Get("type"),
		ArtistName: r.URL.Query().Get("artist_name"),
	}
	if artistIDStr := r.URL.Query().Get("artist_id"); artistIDStr != "" {
		filter.ArtistID, _ = strconv.ParseInt(artistIDStr, 10, 64)
	}

	paginatedData, err := h.service.GetAllPaginated(r.Context(), filter, pagination)
	if err != nil {
		log.Printf("[ERROR INTERNO en Handler] %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Error obteniendo la lista de albums", nil)
		return
	}

	WriteJSON(w, http.StatusOK, paginatedData) // 200
}

// GET ALL BY Artist ID (GET /albums/artist/{artist_id})
func (h *AlbumHandler) GetAlbumsByArtistID(w http.ResponseWriter, r *http.Request) {
	idArtistID := r.PathValue("artist_id")
	idArtist, err := strconv.ParseInt(idArtistID, 10, 64)
	if err != nil {
		log.Printf("[ERROR en Handler] %v\n", err)
		WriteError(w, http.StatusBadRequest, "El ID de la URL debe ser un número entero válido mayor a 0", nil)
		return
	}
	if idArtist <= 0 {
		WriteError(w, http.StatusBadRequest, "El ID de artista debe ser mayor a 0", nil)
		return
	}

	albums, err := h.service.GetAlbumsByArtistID(r.Context(), idArtist)
	if err != nil {
		log.Printf("[ERROR INTERNO en Handler] %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Error interno obteniendo albums", nil)
	}

	// Vacio si no hay albums
	if albums == nil {
		albums = []domain.Album{}
	}

	WriteJSON(w, http.StatusOK, albums) // 200
}

// UPDATE (PUT /albums/{id})
func (h *AlbumHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "El ID de la URL debe ser un número entero válido", nil)
		return
	}
	if id <= 0 {
		WriteError(w, http.StatusBadRequest, "El ID debe ser mayor a 0", nil)
		return
	}

	var input domain.AlbumInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteError(w, http.StatusBadRequest, "Formato JSON inválido", err.Error())
		return
	}

	album, err := h.service.Update(r.Context(), id, &input)
	if err != nil {
		var valErrs domain.ValidationError
		if errors.As(err, &valErrs) {
			WriteError(w, http.StatusBadRequest, "Datos de actualización inválidos", valErrs)
			return
		}

		if errors.Is(err, domain.ErrAlbumNotFound) {
			WriteError(w, http.StatusNotFound, err.Error(), nil) // 404
			return
		}

		// IF caso id artista no existe 400 bad request

		log.Printf("[ERROR INTERNO] PUT /albums/%d: %v\n", id, err)
		WriteError(w, http.StatusInternalServerError, "Error actualizando el álbum", nil) // 500
		return
	}

	WriteJSON(w, http.StatusOK, album)
}

// DELETE (DELETE /albums/{id})
func (h *AlbumHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "El ID de la URL debe ser un número entero válido", nil)
		return
	}
	if id <= 0 {
		WriteError(w, http.StatusBadRequest, "El ID debe ser mayor a 0", nil)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrAlbumNotFound) {
			WriteError(w, http.StatusNotFound, err.Error(), nil) // 404
			return
		}
		log.Printf("[ERROR INTERNO] DELETE /albums/%d: %v\n", id, err)
		WriteError(w, http.StatusInternalServerError, "Error al eliminar el álbum", nil) // 500
		return
	}

	WriteNoContent(w)
}

// Otros

// Agregar nuevo track a album (POST /albums/{id}/tracks)
func (h *AlbumHandler) AddTrack(w http.ResponseWriter, r *http.Request) {
	albumID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "El ID del álbum debe ser un entero válido", nil)
		return
	}
	if albumID <= 0 {
		WriteError(w, http.StatusBadRequest, "El ID del álbum debe ser mayor a 0", nil)
		return
	}

	var input domain.TrackInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteError(w, http.StatusBadRequest, "Formato JSON inválido", err.Error())
		return
	}

	err = h.service.AddTrack(r.Context(), albumID, &input)
	if err != nil {
		var valErrs domain.ValidationError
		if errors.As(err, &valErrs) {
			WriteError(w, http.StatusBadRequest, "Datos de track inválidos", valErrs)
			return
		}

		// Manejo de Conflictos (409) atrapados desde la base de datos, con Errores Centinela
		if errors.Is(err, domain.ErrTrackAlreadyExists) || errors.Is(err, domain.ErrSongAlreadyInAlbum) {
			WriteError(w, http.StatusConflict, err.Error(), nil) // 409 Conflict
			return
		}
		// Manejo de Bad Request (400)
		if errors.Is(err, domain.ErrSongNotInDB) {
			WriteError(w, http.StatusBadRequest, err.Error(), nil) // 400 Bad Request
			return
		}

		log.Printf("[ERROR INTERNO] POST /albums/%d/tracks: %v\n", albumID, err)
		WriteError(w, http.StatusInternalServerError, "Error al agregar el track", nil) // 500
		return
	}

	// Como solo agregamos la relación, un mensaje de éxito simple o un 204 es ideal.
	WriteJSON(w, http.StatusCreated, map[string]string{"message": "Track agregado exitosamente"}) // 200
}

// DELETE (/albums/{id}/tracks/{song_id})
func (h *AlbumHandler) RemoveTrack(w http.ResponseWriter, r *http.Request) {
	albumID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	songID, err2 := strconv.ParseInt(r.PathValue("song_id"), 10, 64)

	if err != nil || err2 != nil || albumID <= 0 || songID <= 0 {
		WriteError(w, http.StatusBadRequest, "Los IDs de la URL deben ser válidos", nil)
		return
	}

	err = h.service.RemoveTrack(r.Context(), albumID, songID)
	if err != nil {
		if errors.Is(err, domain.ErrTrackNotFound) {
			WriteError(w, http.StatusNotFound, err.Error(), nil) // 404
			return
		}

		log.Printf("[ERROR INTERNO] DELETE /albums/%d/tracks/%d: %v\n", albumID, songID, err)
		WriteError(w, http.StatusInternalServerError, "Error al remover el track", nil) // 500
		return
	}

	WriteNoContent(w) // 204
}
