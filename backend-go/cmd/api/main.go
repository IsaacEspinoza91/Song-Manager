package main

import (
	"context"
	"log"

	"github.com/IsaacEspinoza91/Song-Manager/internal/database"
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

}
