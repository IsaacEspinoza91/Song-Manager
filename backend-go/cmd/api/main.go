package main

import (
	"context"
	"log"
	"net/http"

	"github.com/IsaacEspinoza91/Song-Manager/internal/database"
	"github.com/IsaacEspinoza91/Song-Manager/internal/handler"
	"github.com/IsaacEspinoza91/Song-Manager/internal/repository"
	"github.com/IsaacEspinoza91/Song-Manager/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	ctx := context.Background()

	// Inicializar DB
	dbPool, err := database.NewPostgresConnection(ctx)
	if err != nil {
		log.Fatalf("Advertencia: No se encontró el archivo .env o hubo un error al leerlo: %v", err)
	}
	defer dbPool.Close()
	log.Println("Conectado a PostgreSQL exitosamente")

	// Crear repositorios (Inyectar DB)
	artistRepo := repository.NewArtistRepository(dbPool)
	songRepo := repository.NewSongRepository(dbPool)
	albumRepo := repository.NewAlbumRepository(dbPool)

	// Crear servicios (Inyectar repo)
	artistService := service.NewArtistService(artistRepo)
	songService := service.NewSongService(songRepo)
	albumService := service.NewAlbumService(albumRepo)

	// Crar enrutador (Inyectar services). Middleware: Log, CORS, recovery
	router := handler.NewRouter(artistService, songService, albumService)

	// Levantar Server
	log.Println("Servidor corriendo en el puerto 8080...")
	http.ListenAndServe(":8080", router) // Envuelto en middleware

}
