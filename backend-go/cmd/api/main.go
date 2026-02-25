package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IsaacEspinoza91/Song-Manager/internal/config"
	"github.com/IsaacEspinoza91/Song-Manager/internal/database"
	"github.com/IsaacEspinoza91/Song-Manager/internal/handler"
	"github.com/IsaacEspinoza91/Song-Manager/internal/repository"
	"github.com/IsaacEspinoza91/Song-Manager/internal/service"
)

func main() {
	// 1. Cargar Configuración Centralizada
	cfg := config.Load()

	// 2. Inicializar DB
	ctx := context.Background()
	dbPool, err := database.NewPostgresConnection(ctx, cfg.DBUrl)
	if err != nil {
		log.Fatalf("Error fatal conectando a la base de datos: %v", err)
	}
	defer dbPool.Close()
	log.Println("Conectado a PostgreSQL exitosamente")

	// 3. Crear repositorios (Inyectar DB)
	artistRepo := repository.NewArtistRepository(dbPool)
	songRepo := repository.NewSongRepository(dbPool)
	albumRepo := repository.NewAlbumRepository(dbPool)

	// 4. Crear servicios (Inyectar repo)
	artistService := service.NewArtistService(artistRepo)
	songService := service.NewSongService(songRepo)
	albumService := service.NewAlbumService(albumRepo)

	// 5. Crar enrutador (Inyectar services). Middleware: Log, CORS, recovery
	router := handler.NewRouter(artistService, songService, albumService)

	// 6. Config servidor HTTP con Graceful Shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router, // Envuelto en middleware

		// Buena práctica de seguridad, evitar que clientes lentos saturen la API
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Canal para escuchar señales del S.O. (Ctrl+C, Docker Stop, etc)
	// Util para apagado suave. Usa paralelismo.
	quit := make(chan os.Signal, 1) // Crear canal. Pipe
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Levantar server con goroutine (2do plano)
	go func() {
		log.Printf("Servidor corriendo en el puerto %s...\n", cfg.Port)
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error crítico en el servidor HTTP: %v", err)
		}
	}()

	// El hilo principal se bloquea aquí hasta que reciba una señal en el canal 'quit'
	<-quit
	log.Println("Señal de apagado recibida. Iniciando Graceful Shutdown...")

	// Crear un contexto con un límite de tiempo (ej. 10seg)
	// Si el servidor tarda más de 10 seg en terminar peticiones pendientes, forzamos apagado
	// Evitar ataque Slowloris
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("El servidor forzó el apagado debido a un error: %v", err)
	}

	log.Println("Servidor apagado correctamente.")
}
