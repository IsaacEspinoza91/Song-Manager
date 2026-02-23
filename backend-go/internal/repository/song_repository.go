package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/IsaacEspinoza91/Song-Manager/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type songRepository struct {
	db *pgxpool.Pool
}

func NewSongRepository(db *pgxpool.Pool) domain.SongRepository {
	return &songRepository{db: db}
}

// Create inserta una canción y sus relaciones de forma transaccional
func (r *songRepository) Create(ctx context.Context, input *domain.SongInput) (*domain.Song, error) {
	// 1. Iniciar la transacción
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error iniciando transacción para create: %w", err)
	}
	defer tx.Rollback(ctx) // Roolback si no se llega al commit

	// 2. Insertar la canción maestra
	querySong := `
		INSERT INTO songs (title, duration, updated_at) 
		VALUES ($1, $2, NOW()) 
		RETURNING id
	`
	var songID int64
	// Usar tx en vez de r.db para consulta
	err = tx.QueryRow(ctx, querySong, input.Title, input.Duration).Scan(&songID)
	if err != nil {
		return nil, fmt.Errorf("error insertando la canción: %w", err)
	}

	// 3. Insertar las relaciones en la tabla intermedia (si es que enviaron artistas)
	if len(input.Artists) > 0 {
		queryArtistSong := `
			INSERT INTO song_artists (song_id, artist_id, role) 
			VALUES ($1, $2, $3)
		`

		for _, a := range input.Artists {
			_, err := tx.Exec(ctx, queryArtistSong, songID, a.ArtistID, a.Role)
			if err != nil {
				// Al salir sin hacer Commit, el defer ejecuta el Rollback.
				return nil, fmt.Errorf("error asociando el artista ID %d a la canción: %w", a.ArtistID, err)
			}
		}
	}

	// 4. Commit si todo sale bien
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error confirmando la transacción: %w", err)
	}

	// 5. Obtener mediante GetByID para devolver el objeto completo de song con artists
	fullSong, err := r.GetByID(ctx, songID)
	if err != nil {
		return nil, fmt.Errorf("canción creada, pero error al obtener detalles: %w", err)
	}

	return fullSong, nil
}

// READ
func (r *songRepository) GetByID(ctx context.Context, id int64) (*domain.Song, error) {
	// 1. Obtener los datos principales de la Canción
	var song domain.Song
	querySong := `
		SELECT id, title, duration, created_at, updated_at
		FROM songs
		WHERE id = $1 AND deleted_at IS NULL
	`
	err := r.db.QueryRow(ctx, querySong, id).Scan(
		&song.ID,
		&song.Title,
		&song.Duration,
		&song.CreatedAt,
		&song.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("canción no encontrada")
		}
		return nil, fmt.Errorf("error obteniendo la canción: %w", err)
	}

	// 2. Obtener los artistas y sus roles usando JOIN
	// Hacemos JOIN entre artists y song_artists para cruzar los datos
	queryArtists := `
		SELECT a.id, a.name, asg.role
		FROM artists a
		INNER JOIN song_artists asg ON a.id = asg.artist_id
		WHERE asg.song_id = $1 AND a.deleted_at IS NULL
	`
	rows, err := r.db.Query(ctx, queryArtists, id)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo artistas de la canción: %w", err)
	}
	defer rows.Close()

	// Inicializamos el slice para evitar que el JSON devuelva "artists": null
	song.Artists = []domain.ArtistWithRole{}

	for rows.Next() {
		var artist domain.ArtistWithRole
		if err := rows.Scan(&artist.ID, &artist.Name, &artist.Role); err != nil {
			return nil, fmt.Errorf("error escaneando artista: %w", err)
		}
		song.Artists = append(song.Artists, artist)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando los artistas: %w", err)
	}

	return &song, nil
}

// Mala practica iterar en un bucle y llamar a GetByID cada uno, problema N+1, mal rendimiento
func (r *songRepository) GetAll(ctx context.Context) ([]domain.Song, error) {
	// 1. Obtener datos songs
	querySongs := `
		SELECT id, title, duration, created_at, updated_at 
		FROM songs 
		WHERE deleted_at IS NULL 
		ORDER BY id ASC
	`
	rows, err := r.db.Query(ctx, querySongs)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo canciones: %w", err)
	}
	defer rows.Close()

	var songs []domain.Song
	var songIDs []int64 // Guardaremos los IDs para buscar sus artistas de golpe

	for rows.Next() {
		var s domain.Song
		if err := rows.Scan(&s.ID, &s.Title, &s.Duration, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error escaneando canción: %w", err)
		}

		s.Artists = []domain.ArtistWithRole{} // Iniciarlizar slice evita null en JSON
		songs = append(songs, s)
		songIDs = append(songIDs, s.ID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando canciones: %w", err)
	}

	// Si no hay canciones, retornamos el arreglo vacío inmediatamente
	if len(songs) == 0 {
		return songs, nil
	}

	// 2. Obtener artitas de canciones en una sola consulta
	// Funcionalidad pgx: Permite pasar slice a ANY($1)
	queryArtists := `
		SELECT asg.song_id, a.id, a.name, asg.role
		FROM artists a
		INNER JOIN song_artists asg ON a.id = asg.artist_id
		WHERE asg.song_id = ANY($1) AND a.deleted_at IS NULL
	`
	artistRows, err := r.db.Query(ctx, queryArtists, songIDs)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo artistas en bloque: %w", err)
	}
	defer artistRows.Close()

	// 3. Unir artistas con canciones
	songMap := make(map[int64]*domain.Song)
	for i := range songs {
		songMap[songs[i].ID] = &songs[i]
	}

	for artistRows.Next() {
		var songID int64
		var artist domain.ArtistWithRole
		if err := artistRows.Scan(&songID, &artist.ID, &artist.Name, &artist.Role); err != nil {
			return nil, fmt.Errorf("error escaneando artista relacional: %w", err)
		}

		// Agregamos el artista al slice de la canción correspondiente
		if song, exists := songMap[songID]; exists {
			song.Artists = append(song.Artists, artist)
		}
	}

	return songs, nil
}

// Utilizar subconsultas y no INNER JOIN, esto evita duplicados que rompan paginacion
func (r *songRepository) GetAllPaginated(ctx context.Context, filter domain.SongFilter, params domain.PaginationParams) (*domain.PaginatedResult[domain.Song], error) {

	baseQuery := `SELECT id, title, duration, created_at, updated_at FROM songs WHERE deleted_at IS NULL`
	countQuery := `SELECT COUNT(*) FROM songs WHERE deleted_at IS NULL`

	var args []interface{}
	argID := 1

	// Busqueda parcial nombre cancion
	if filter.Title != "" {
		condition := fmt.Sprintf(" AND title ILIKE $%d", argID)
		baseQuery += condition
		countQuery += condition
		args = append(args, "%"+filter.Title+"%")
		argID++
	}
	// Busqueda exacta artista id
	if filter.ArtistID > 0 {
		condition := fmt.Sprintf(" AND id IN (SELECT song_id FROM song_artists WHERE artist_id = $%d)", argID)
		baseQuery += condition
		countQuery += condition
		args = append(args, filter.ArtistID)
		argID++
	}
	// Busqueda parcial nombre artista
	if filter.ArtistName != "" {
		condition := fmt.Sprintf(` AND id IN (
			SELECT asg.song_id 
			FROM song_artists asg 
			INNER JOIN artists a ON asg.artist_id = a.id 
			WHERE a.name ILIKE $%d AND a.deleted_at IS NULL
		)`, argID)

		baseQuery += condition
		countQuery += condition
		args = append(args, "%"+filter.ArtistName+"%")
		argID++
	}

	// Conteo total
	var totalItems int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, fmt.Errorf("error contando las canciones para paginación: %w", err)
	}

	// Obtener canciones paginadas
	baseQuery += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando query paginada de canciones: %w", err)
	}
	defer rows.Close()

	// Crear slice de songs (sin artistas todavia)
	var songs []domain.Song
	var songIDs []int64
	for rows.Next() {
		var s domain.Song
		if err := rows.Scan(&s.ID, &s.Title, &s.Duration, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error escaneando canción paginada: %w", err)
		}
		s.Artists = []domain.ArtistWithRole{}
		songs = append(songs, s)
		songIDs = append(songIDs, s.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando canciones paginadas: %w", err)
	}

	// Asignar artistas en bloquea songs
	if len(songs) > 0 {
		queryArtists := `
			SELECT asg.song_id, a.id, a.name, asg.role
			FROM artists a
			INNER JOIN song_artists asg ON a.id = asg.artist_id
			WHERE asg.song_id = ANY($1) AND a.deleted_at IS NULL
		`
		artistRows, err := r.db.Query(ctx, queryArtists, songIDs)
		if err != nil {
			return nil, fmt.Errorf("error obteniendo artistas para canciones paginadas: %w", err)
		}
		defer artistRows.Close()

		songMap := make(map[int64]*domain.Song)
		for i := range songs {
			songMap[songs[i].ID] = &songs[i]
		}
		for artistRows.Next() {
			var songID int64
			var artist domain.ArtistWithRole
			if err := artistRows.Scan(&songID, &artist.ID, &artist.Name, &artist.Role); err != nil {
				return nil, fmt.Errorf("error escaneando artista relacional paginado: %w", err)
			}
			if song, exists := songMap[songID]; exists {
				song.Artists = append(song.Artists, artist)
			}
		}
	} else {
		songs = []domain.Song{} // [] vacio en vez de nil
	}

	return domain.NewPaginatedResult(songs, totalItems, params.Page, params.Limit), nil
}

// UPDATE
// PUT clasico, actualiza todo slos datos de la tabla principal, elimina las relaciones existentes
// y las inserta de nuevo.
func (r *songRepository) Update(ctx context.Context, id int64, input *domain.SongInput) (*domain.Song, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error iniciando transacción para update: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Actualizar tabla principal (song)
	updateSongQuery := `
		UPDATE songs
		SET title = $1, duration = $2, updated_at = NOW()
		WHERE id = $3 AND deleted_at IS NULL
	`
	// Usar Exec porque no necesitamos hacer Scan de nada. Solo saber si afecto una fila
	res, err := tx.Exec(ctx, updateSongQuery, input.Title, input.Duration, id)
	if err != nil {
		return nil, fmt.Errorf("error actualizando los datos de la canción: %w", err)
	}

	// Si no se afectó ninguna fila, es porque el ID no existe o la canción fue borrada lógicamente
	if res.RowsAffected() == 0 {
		return nil, errors.New("canción no encontrada")
	}

	// 2. Limpiar relaciones antiguas
	deletedRelsQuery := `DELETE FROM song_artists WHERE song_id = $1`
	_, err = tx.Exec(ctx, deletedRelsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("error limpiando relaciones antiguas de la canción: %w", err)
	}

	// 3. Insertar nuevas relaciones, solo si vienen en input dto
	if len(input.Artists) > 0 {
		insertRelQuery := `
			INSERT INTO song_artists (song_id, artist_id, role)
			VALUES ($1, $2, $3)
		`
		for _, a := range input.Artists {
			_, err := tx.Exec(ctx, insertRelQuery, id, a.ArtistID, a.Role)
			if err != nil {
				// falta caso retornar error id de artista no existe
				return nil, fmt.Errorf("error insertando nueva relación con artista ID %d: %w", a.ArtistID, err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error confirmando la transacción de actualización: %w", err)
	}

	// Retornar la canción completa con sus datos y artistas actualizados
	// Usar r.db, ya que la transacción ya se cerró y guardó.
	updatedSong, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("canción actualizada, pero error al obtener detalles: %w", err)
	}

	return updatedSong, nil
}

// DELETE
func (r *songRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE songs SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error eliminando a la canción ID %d: %w", id, err)
	}
	if res.RowsAffected() == 0 {
		return errors.New("canción no encontrada")
	}
	return nil
}
