package handler

import (
	"encoding/json"
	"errors"
	"log"
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
		WriteError(w, http.StatusBadRequest, "Formato JSON inválido", err.Error())
		return
	}

	// 2. Llamar a la capa de servicio
	// Contexto (r.Context()) viaja desde aquí hasta la base de datos
	artist, err := h.service.Create(r.Context(), &input)
	if err != nil {
		// Verificamos si el error es de tipo ValidationError (acumulación de errores)
		var valErrs domain.ValidationError
		if errors.As(err, &valErrs) {
			// Si lo es, pasamos el mapa completo a los Details
			WriteError(w, http.StatusBadRequest, "Datos de entrada inválidos", valErrs)
			return
		}

		// cualquier otro error (ej: base de datos caida), responder 400
		WriteError(w, http.StatusBadRequest, "No se pudo crear el artista", err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, artist) // 201 Created
}

// GET ALL (GET /artists/all)
func (h *ArtistHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	artists, err := h.service.GetAll(r.Context())
	if err != nil {
		// Si falla la base de datos, es un error interno del servidor (500)
		log.Printf("[ERROR INTERNO en Handler] %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Error interno obteniendo artistas", nil)
		return
	}

	// Si no hay artistas, devolver array vacio en vez de nil
	if artists == nil {
		artists = []domain.Artist{}
	}

	WriteJSON(w, http.StatusOK, artists) // 200 OK
}

// GET ID (GET /artists/{id})
func (h *ArtistHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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

	artist, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrArtistNotFound) {
			WriteError(w, http.StatusNotFound, err.Error(), nil) // 404
			return
		}

		log.Printf("[ERROR INTERNO en Handler] %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Error al buscar el artista", nil) // 500
		return
	}

	WriteJSON(w, http.StatusOK, artist) // 200 OK
}

// GET ALL PAG (GET /artists?page=2&limit=5&genre=rock&country=chile)
func (h *ArtistHandler) GetAllPaginated(w http.ResponseWriter, r *http.Request) {
	// Extraer query params
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10 // valor default
	}
	pagination := domain.PaginationParams{
		Page:  page,
		Limit: limit,
	}

	filter := domain.ArtistFilter{
		Name:    r.URL.Query().Get("name"),
		Genre:   r.URL.Query().Get("genre"),
		Country: r.URL.Query().Get("country"),
	}

	// Llamar servicio
	paginatedData, err := h.service.GetAllPaginated(r.Context(), filter, pagination)
	if err != nil {
		log.Printf("[ERROR INTERNO en Handler] %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Error obteniendo la lista de artistas", nil)
		return
	}

	WriteJSON(w, http.StatusOK, paginatedData) // 200
}

// UPDATE (PUT /artists/{id})
func (h *ArtistHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Extraer ID, definido en router {id}
	idString := r.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 64) // Convertir string a int64 (base 10, 64 bits)
	if err != nil {
		log.Printf("[ERROR en Handler] %v\n", err)
		WriteError(w, http.StatusBadRequest, "El ID de la URL debe ser un número entero válido mayor a 0", nil)
		return
	}
	if id <= 0 {
		WriteError(w, http.StatusBadRequest, "El ID debe ser mayor a 0", nil)
		return
	}

	// Decodificar el JSON entrante al DTO (ArtistInput)
	var input domain.ArtistInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteError(w, http.StatusBadRequest, "Formato JSON inválido", err.Error())
		return
	}

	// Pasar el ID y los datos a la capa de Servicio
	artist, err := h.service.Update(r.Context(), id, &input)
	if err != nil {
		// Evaluar si es error de validación de campos
		var valErrs domain.ValidationError
		if errors.As(err, &valErrs) {
			WriteError(w, http.StatusBadRequest, "Datos de actualización inválidos", valErrs)
			return
		}

		// Caso error fue porque no se encontró el artista
		if errors.Is(err, domain.ErrArtistNotFound){
			WriteError(w, http.StatusNotFound, err.Error(), nil) // 404
			return
		}
		// Cualquier otro error de validación o base de datos
		log.Printf("[ERROR INTERNO en Handler] %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Error actualizando al artista", nil) // 500
		return
	}

	WriteJSON(w, http.StatusOK, artist) // 200 OK
}

// DELETE (DELETE /artists/{id})
func (h *ArtistHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Extraer y convertir el ID
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

	// Llamar al servicio. Soft Delete
	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrArtistNotFound) {
			WriteError(w, http.StatusNotFound, err.Error(), nil) // 404
			return
		}
		log.Printf("[ERROR INTERNO en Handler] %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Error al eliminar el artista", nil) // 500
		return
	}

	// Responder al frontend indicando exito, sin contenido en el body
	// 204 No Content es el código HTTP estándar para un DELETE exitoso
	WriteNoContent(w)
}
