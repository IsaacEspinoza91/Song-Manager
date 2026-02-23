package handler

import (
	"net/http"

	"github.com/IsaacEspinoza91/Song-Manager/internal/domain"
)

// NewRouter recibe TODOS los servicios y retorna un http.Handler listo para usar
func NewRouter(artistService domain.ArtistService, songService domain.SongService) http.Handler {
	mux := http.NewServeMux()

	// Instanciar los handlers específicos inyectándoles su servicio correspondiente
	artistHandler := NewArtistHandler(artistService)
	songHandler := NewSongHandler(songService)

	// 2. Registramos las rutas (Requiere Go 1.22+)
	mux.HandleFunc("POST /artists", artistHandler.Create)
	mux.HandleFunc("GET /artists/all", artistHandler.GetAll)
	mux.HandleFunc("GET /artists", artistHandler.GetAllPaginated)
	mux.HandleFunc("GET /artists/{id}", artistHandler.GetByID)
	mux.HandleFunc("PUT /artists/{id}", artistHandler.Update)
	mux.HandleFunc("DELETE /artists/{id}", artistHandler.Delete)

	mux.HandleFunc("POST /songs", songHandler.Create)
	mux.HandleFunc("GET /songs/{id}", songHandler.GetByID)

	return mux
}
