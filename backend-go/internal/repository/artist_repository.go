package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/IsaacEspinoza91/Song-Manager/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// artistRepository Maneja comunicacion con db
type artistRepository struct {
	db *pgxpool.Pool
}

// Al implementar la interfaz, cualquier otra parte del codigo solo puede usar los metodos definidos ahi
// Ademas, si en el futuro cambia la db, solo debo crear otro archivo repo y cumplir la interfaz
func NewArtistRepository(db *pgxpool.Pool) domain.ArtistRepository {
	return &artistRepository{db: db}
}

// 1. Create
func (r *artistRepository) Create(ctx context.Context, input *domain.ArtistInput) (*domain.Artist, error) {
	var artist domain.Artist
	query := `
		INSERT INTO artists (name, genre, country, bio, image_url) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, name, genre, country, bio, image_url, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query, input.Name, input.Genre, input.Country, input.Bio, input.ImageURL).
		Scan(
			&artist.ID,
			&artist.Name,
			&artist.Genre,
			&artist.Country,
			&artist.Bio,
			&artist.ImageURL,
			&artist.CreatedAt,
			&artist.UpdatedAt,
		)
	if err != nil {
		return nil, fmt.Errorf("error creando al artista: %w", err)
	}
	return &artist, nil
}

// 2. READ
func (r *artistRepository) GetByID(ctx context.Context, id int64) (*domain.Artist, error) {
	// Funciona sin RETURNING porque no modifica datos
	query := `
		SELECT id, name, genre, country, bio, image_url, created_at, updated_at
		FROM artists 
		WHERE id = $1 AND deleted_at IS NULL
	`

	var a domain.Artist
	err := r.db.QueryRow(ctx, query, id).Scan(&a.ID, &a.Name, &a.Genre, &a.Country, &a.Bio, &a.ImageURL, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("artista no encontrado")
		}
		return nil, fmt.Errorf("error obteniendo al artista: %w", err)
	}
	return &a, nil
}

// Retorno de slice, los artistas se guardan juntos en memoria y es mas facil de procesar.
// Usa slice de punteros cuando la estructura es gigante
func (r *artistRepository) GetAll(ctx context.Context) ([]domain.Artist, error) {
	query := `
		SELECT id, name, genre, country, bio, image_url, created_at, updated_at 
		FROM artists
		WHERE deleted_at IS NULL
		ORDER BY id ASC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando query para obtener los artistas: %w", err)
	}
	defer rows.Close()

	// No usar porque no inicializa slice, si no hay objetos retorna nil
	//var artists []domain.Artist
	artists := make([]domain.Artist, 0) // o también: artists := []domain.Artist{}
	for rows.Next() {
		var a domain.Artist
		err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Genre,
			&a.Country,
			&a.Bio,
			&a.ImageURL,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando fila de artista: %w", err)
		}

		artists = append(artists, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando las filas de artistas: %w", err)
	}

	return artists, nil
}

func (r *artistRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM artists WHERE deleted_at IS NULL`

	var total int64
	err := r.db.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// No es recomendable usar transaccion en este metodo porque no modifican la base de datos.
// Efectivamente puede haber un pequeño error de datos, pero nada grave. Es mejor asumir ese riesgo
// que usar recursos adicionales de memoria para bloquear la db en ese momento
func (r *artistRepository) GetAllPaginated(ctx context.Context, filter domain.ArtistFilter, params domain.PaginationParams) (*domain.PaginatedResult[domain.Artist], error) {
	// 1. Consultas base
	baseQuery := `SELECT id, name, genre, country, bio, image_url, created_at, updated_at FROM artists WHERE deleted_at IS NULL`
	countQuery := `SELECT COUNT(*) FROM artists WHERE deleted_at IS NULL`

	var args []interface{} // Slice de cualquier tipo
	argID := 1             // Rastreador de la posición del parámetro ($1, $2, etc.)

	// 2. Construir filtros dinámicamente
	if filter.Name != "" {
		// Agregamos la condición a ambas consultas
		condition := fmt.Sprintf(" AND name ILIKE $%d", argID) // ILIKE ignora mayusculas o minusculas
		baseQuery += condition
		countQuery += condition

		// Agregamos el valor envuelto en '%' para la búsqueda parcial
		args = append(args, "%"+filter.Name+"%")
		argID++
	}
	if filter.Genre != "" {
		condition := fmt.Sprintf(" AND genre ILIKE $%d", argID)
		baseQuery += condition
		countQuery += condition
		args = append(args, "%"+filter.Genre+"%")
		argID++
	}
	if filter.Country != "" {
		condition := fmt.Sprintf(" AND country ILIKE $%d", argID)
		baseQuery += condition
		countQuery += condition
		args = append(args, "%"+filter.Country+"%")
		argID++
	}

	// 3. Obtener el total de elementos (para calcular las páginas)
	var totalItems int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, fmt.Errorf("error contando los artistas para paginación: %w", err)
	}

	// 4. Agregar Paginación a la consulta principal
	baseQuery += fmt.Sprintf(" ORDER BY id ASC LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, params.Limit, params.GetOffset())

	// 5. Ejecutar la consulta final
	rows, err := r.db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando query paginada de artistas: %w", err)
	}
	defer rows.Close()

	var artists []domain.Artist
	for rows.Next() {
		var a domain.Artist
		err := rows.Scan(&a.ID, &a.Name, &a.Genre, &a.Country, &a.Bio, &a.ImageURL, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error escaneando artista en paginación: %w", err)
		}
		artists = append(artists, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando filas en paginación: %w", err)
	}

	// 6. Retornar la estructura genérica armada
	if artists == nil {
		artists = []domain.Artist{} // Evitar null en JSON
	}
	return domain.NewPaginatedResult(artists, totalItems, params.Page, params.Limit), nil
}

// 3. Update
func (r *artistRepository) Update(ctx context.Context, id int64, input *domain.ArtistInput) (*domain.Artist, error) {
	var artist domain.Artist
	query := `
		UPDATE artists 
		SET name = $1, genre = $2, country = $3, bio = $4, image_url = $5, updated_at = NOW() 
		WHERE id = $6 AND deleted_at IS NULL
		RETURNING id, name, genre, country, bio, image_url, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, input.Name, input.Genre, input.Country, input.Bio, input.ImageURL, id).
		Scan(
			&artist.ID,
			&artist.Name,
			&artist.Genre,
			&artist.Country,
			&artist.Bio,
			&artist.ImageURL,
			&artist.CreatedAt,
			&artist.UpdatedAt,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("artista no encontrado")
		}
		return nil, fmt.Errorf("error actualizando al artista ID %d: %w", id, err)
	}
	return &artist, nil
}

// 4. Delete
func (r *artistRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE artists SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("artista no encontrado")
		}
		return fmt.Errorf("error eliminando al artista ID %d: %w", id, err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("error ningun artista eliminado: %w", err)
	}
	return nil
}

// Nota que las transacciones se usan cuando se van a ejecutar 2 o mas operaciones sql
// Para una operacion simple (1 query) no es necesario porque en si ya es una transaccion implicita

// GetALL no pag deleted  SOLO ADMIN
// Get pag
// get pag con filtros  preguntar ia si es buena practica dar la opcion de seleccionar deleted o no
// Get discos con canciones
