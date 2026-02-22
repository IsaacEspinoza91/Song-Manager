package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool" // Mejor rendimiento que lib/pq
)

// Abre conexion, hace ping y retorna pool de conexiones
func NewPostgresConnection(ctx context.Context) (*pgxpool.Pool, error) {
	// Contruccion URL usando var entorno
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// pgxpool maneja el conjunto de conexiones abiertas
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("error al crear el pool de conexiones: %v", err)
	}

	// Verificar conexion ping
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error al hacer ping a la base de datos: %v", err)
	}

	return pool, nil
}
