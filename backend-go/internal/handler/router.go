package handler

import (
	"net/http"

	"github.com/IsaacEspinoza91/Song-Manager/internal/domain"
	"github.com/IsaacEspinoza91/Song-Manager/internal/middleware"
)

// NewRouter recibe TODOS los servicios y retorna un http.Handler listo para usar
func NewRouter(artistService domain.ArtistService, songService domain.SongService, albumService domain.AlbumService) http.Handler {
	mux := http.NewServeMux()

	// Instanciar los handlers específicos inyectándoles su servicio correspondiente
	artistHandler := NewArtistHandler(artistService)
	songHandler := NewSongHandler(songService)
	albumHandler := NewAlbumHandler(albumService)

	// Registramos las rutas (Requiere Go 1.22+)
	mux.HandleFunc("POST /artists", artistHandler.Create)
	mux.HandleFunc("GET /artists/all", artistHandler.GetAll)
	mux.HandleFunc("GET /artists", artistHandler.GetAllPaginated)
	mux.HandleFunc("GET /artists/{id}", artistHandler.GetByID)
	mux.HandleFunc("PUT /artists/{id}", artistHandler.Update)
	mux.HandleFunc("DELETE /artists/{id}", artistHandler.Delete)

	mux.HandleFunc("POST /songs", songHandler.Create)
	mux.HandleFunc("GET /songs/{id}", songHandler.GetByID)
	mux.HandleFunc("GET /songs/all", songHandler.GetAll)
	mux.HandleFunc("GET /songs", songHandler.GetAllPaginated)
	mux.HandleFunc("PUT /songs/{id}", songHandler.Update)
	mux.HandleFunc("DELETE /songs/{id}", songHandler.Delete)
	mux.HandleFunc("DELETE /songs/{id}/artist/{artist_id}", songHandler.RemoveArtist)
	mux.HandleFunc("POST /songs/{id}/artist", songHandler.AddArtist)

	mux.HandleFunc("POST /albums", albumHandler.Create)
	mux.HandleFunc("GET /albums/{id}", albumHandler.GetByID)
	mux.HandleFunc("GET /albums", albumHandler.GetAllPaginated)
	mux.HandleFunc("GET /albums/artist/{artist_id}", albumHandler.GetAlbumsByArtistID)
	mux.HandleFunc("PUT /albums/{id}", albumHandler.Update)
	mux.HandleFunc("DELETE /albums/{id}", albumHandler.Delete)
	mux.HandleFunc("POST /albums/{id}/tracks", albumHandler.AddTrack)
	mux.HandleFunc("DELETE /albums/{id}/tracks/{song_id}", albumHandler.RemoveTrack)

	// Middleware

	// El orden importa:
	// Primero el Logger anota la entrada.
	// Luego el CORS revisa los permisos.
	// Rate Limiting. Contra ataques masivos de una IP
	// Recovery en caso de panic
	// Finalmente, llega al Mux (enrutador).

	handlerConCORS := middleware.CORS(mux)
	handlerConLogger := middleware.Logger(handlerConCORS)
	handlerConRateLimit := middleware.RateLimit(handlerConLogger)
	handlerFinal := middleware.Recovery(handlerConRateLimit)

	return handlerFinal
}
