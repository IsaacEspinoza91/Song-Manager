package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/IsaacEspinoza91/Song-Manager/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type albumRepository struct {
	db *pgxpool.Pool
}

func NewAlbumRepository(db *pgxpool.Pool) domain.AlbumRepository {
	return &albumRepository{db: db}
}

func (r *albumRepository) Create(ctx context.Context, input *domain.AlbumInput) (*domain.Album, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error iniciando transacción para crear álbum: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insertar album. CreateAt y UpdateAt NOW() por def en DB al crear
	queryAlbum := `
		INSERT INTO albums (title, release_date, type, cover_url, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id
	`
	var albumID int64
	// input.ReleaseDate (string "YYYY-MM-DD") es casteado automáticamente a DATE por Postgres.
	err = tx.QueryRow(ctx, queryAlbum, input.Title, input.ReleaseDate, input.Type, input.CoverURL).Scan(&albumID)
	if err != nil {
		return nil, fmt.Errorf("error insertando el album: %w", err)
	}

	// Insertar relacion artista. En caso de no incluir artistas, solo inserta album
	if len(input.Artists) > 0 {
		queryAlbumArtist := `
			INSERT INTO album_artists (album_id, artist_id, is_primary)
			VALUES ($1, $2, $3)
		`
		for _, a := range input.Artists {
			_, err := tx.Exec(ctx, queryAlbumArtist, albumID, a.ArtistID, a.IsPrimary)
			if err != nil {
				return nil, fmt.Errorf("error asociando el artista ID %d al album: %w", a.ArtistID, err)
			}
		}
	}

	// Insertar relacion cancion (track)
	if len(input.Tracks) > 0 {
		queryTrack := `
			INSERT INTO tracks (album_id, song_id, track_number)
			VALUES ($1, $2, $3)
		`
		for _, t := range input.Tracks {
			_, err := tx.Exec(ctx, queryTrack, albumID, t.SongID, t.TrackNumber)
			if err != nil {
				return nil, fmt.Errorf("error asociando la cancion ID %d como track %d: %w", t.SongID, t.TrackNumber, err)
			}
		}
	}

	// Hacer commit
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error confirmando la transacción del álbum: %w", err)
	}

	// Obtener album segun ID
	fullAlbum, err := r.GetByID(ctx, albumID)
	if err != nil {
		return nil, fmt.Errorf("álbum creado con éxito, pero falló al obtener los detalles: %w", err)
	}
	return fullAlbum, nil
}

// POST /albums/{id}/tracks
func (r *albumRepository) AddTrack(ctx context.Context, albumID int64, input *domain.TrackInput) error {
	query := `
		INSERT INTO tracks (album_id, song_id, track_number)
		VALUES ($1, $2, $3)
	`
	// Uso de Exec, porque solo insertamos en tabla intermedia
	_, err := r.db.Exec(ctx, query, albumID, input.SongID, input.TrackNumber)
	if err != nil {
		var pgErr *pgconn.PgError // Verificar si error viene de postgres

		if errors.As(err, &pgErr) {
			// 23505: unique_violation
			if pgErr.Code == "23505" {
				// Evaluamos qué restricción falló
				if pgErr.ConstraintName == "tracks_album_id_track_number_key" {
					return errors.New("este número de pista ya está ocupado en el álbum")
				}
				if pgErr.ConstraintName == "tracks_pkey" {
					return errors.New("esta canción ya existe en este álbum")
				}
			}
			// Código 23503: foreign_key_violation (song_id no existe)
			if pgErr.Code == "23503" {
				return errors.New("la canción indicada no existe en la base de datos")
			}
		}

		return fmt.Errorf("error agregando el track %d al álbum %d: %w", input.SongID, albumID, err)
	}
	return nil
}

// DELETE /albums/{id}/tracks/{song_id}
func (r *albumRepository) RemoveTrack(ctx context.Context, albumID int64, songID int64) error {
	query := `DELETE FROM tracks WHERE album_id = $1 AND song_id = $2`

	res, err := r.db.Exec(ctx, query, albumID, songID)
	if err != nil {
		return fmt.Errorf("error eliminando el track %d del álbum %d: %w", songID, albumID, err)
	}
	if res.RowsAffected() == 0 {
		return errors.New("track no encontrado en este álbum")
	}
	return nil
}

func (r *albumRepository) GetByID(ctx context.Context, id int64) (*domain.Album, error) {
	var album domain.Album

	// 1. Obtener los datos principales del Álbum
	queryAlbum := `
		SELECT id, title, release_date, type, cover_url, created_at, updated_at
		FROM albums
		WHERE id = $1 AND deleted_at IS NULL
	`
	err := r.db.QueryRow(ctx, queryAlbum, id).Scan(
		&album.ID,
		&album.Title,
		&album.ReleaseDate,
		&album.Type,
		&album.CoverURL, // Soporta *string automáticamente
		&album.CreatedAt,
		&album.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("álbum no encontrado")
		}
		return nil, fmt.Errorf("error obteniendo el álbum principal: %w", err)
	}

	// 2. Obtener artistas de Album
	album.Artists = []domain.AlbumArtist{} // Inicializamos vacío para el JSON
	queryArtist := `
		SELECT a.id, a.name, aa.is_primary
		FROM artists a
		INNER JOIN album_artists aa ON a.id = aa.artist_id
		WHERE aa.album_id = $1 AND a.deleted_at IS NULL
	`
	artistRows, err := r.db.Query(ctx, queryArtist, id)
	if err != nil {
		return nil, fmt.Errorf("error consultando los artistas del álbum: %w", err)
	}
	defer artistRows.Close()

	for artistRows.Next() {
		var artist domain.AlbumArtist
		if err := artistRows.Scan(&artist.ID, &artist.Name, &artist.IsPrimary); err != nil {
			return nil, fmt.Errorf("error escaneando artista del álbum: %w", err)
		}
		album.Artists = append(album.Artists, artist)
	}
	if err := artistRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando los artistas del álbum: %w", err)
	}

	// 3. Obtener los tracks
	album.Tracks = []domain.Track{}
	// JOIN con songs para traernos el título y la duración. Es vital el ORDER BY track_number
	queryTracks := `
		SELECT t.track_number, s.id, s.title, s.duration
		FROM tracks t
		INNER JOIN songs s ON t.song_id = s.id
		WHERE t.album_id = $1 AND s.deleted_at IS NULL
		ORDER BY t.track_number ASC
	`
	trackRows, err := r.db.Query(ctx, queryTracks, id)
	if err != nil {
		return nil, fmt.Errorf("error consultando los tracks del álbum: %w", err)
	}
	defer trackRows.Close()

	for trackRows.Next() {
		var track domain.Track
		err := trackRows.Scan(&track.TrackNumber, &track.SongID, &track.Title, &track.Duration)
		if err != nil {
			return nil, fmt.Errorf("error escaneando track del álbum: %w", err)
		}
		album.Tracks = append(album.Tracks, track)
	}
	if err := trackRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando los tracks del álbum: %w", err)
	}

	return &album, nil
}

func (r *albumRepository) GetAllPaginated(ctx context.Context, filter domain.AlbumFilter, params domain.PaginationParams) (*domain.PaginatedResult[domain.Album], error) {
	baseQuery := `SELECT id, title, release_date, type, cover_url, created_at, updated_at FROM albums WHERE deleted_at IS NULL`
	countQuery := `SELECT COUNT(*) FROM albums WHERE deleted_at IS NULL`

	var args []interface{}
	var argID int64 = 1
	// 1. Filtro por Título
	if filter.Title != "" {
		condition := fmt.Sprintf(" AND title ILIKE $%d", argID)
		baseQuery += condition
		countQuery += condition
		args = append(args, "%"+filter.Title+"%")
		argID++
	}

	// 2. Filtro por Tipo (EP, LP, Single)
	if filter.Type != "" {
		condition := fmt.Sprintf(" AND type = $%d", argID)
		baseQuery += condition
		countQuery += condition
		args = append(args, filter.Type)
		argID++
	}

	// 3. Filtro por ID de Artista (Subconsulta)
	if filter.ArtistID > 0 {
		condition := fmt.Sprintf(" AND id IN (SELECT album_id FROM album_artists WHERE artist_id = $%d)", argID)
		baseQuery += condition
		countQuery += condition
		args = append(args, filter.ArtistID)
		argID++
	}

	// 4. Filtro por Nombre de Artista (Subconsulta con JOIN)
	if filter.ArtistName != "" {
		condition := fmt.Sprintf(` AND id IN (
			SELECT aa.album_id 
			FROM album_artists aa 
			INNER JOIN artists a ON aa.artist_id = a.id 
			WHERE a.name ILIKE $%d AND a.deleted_at IS NULL
		)`, argID)

		baseQuery += condition
		countQuery += condition
		args = append(args, "%"+filter.ArtistName+"%")
		argID++
	}

	// Contar items
	var totalItems int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, fmt.Errorf("error contando álbumes: %w", err)
	}

	// Consulta de paginacion
	baseQuery += fmt.Sprintf(" ORDER BY release_date DESC, id ASC LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo álbumes paginados: %w", err)
	}
	defer rows.Close()

	// Mapear rows a respuesta
	var albums []domain.Album
	var albumIDs []int64
	for rows.Next() {
		var a domain.Album
		err := rows.Scan(&a.ID, &a.Title, &a.ReleaseDate, &a.Type, &a.CoverURL, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error escaneando álbum: %w", err)
		}
		a.Artists = []domain.AlbumArtist{} // Inicializar
		a.Tracks = []domain.Track{}        // Inicializar vacío intencionalmente (Summary View)
		albums = append(albums, a)
		albumIDs = append(albumIDs, a.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando álbumes: %w", err)
	}

	// Traer bloque de artistas
	if len(albums) > 0 {
		queryArtists := `
			SELECT aa.album_id, a.id, a.name, aa.is_primary
			FROM artists a
			INNER JOIN album_artists aa ON a.id = aa.artist_id
			WHERE aa.album_id = ANY($1) AND a.deleted_at IS NULL
		`
		artistRows, err := r.db.Query(ctx, queryArtists, albumIDs)
		if err != nil {
			return nil, fmt.Errorf("error obteniendo artistas en bloque para álbumes: %w", err)
		}
		defer artistRows.Close()

		albumMap := make(map[int64]*domain.Album)
		for i := range albums {
			albumMap[albums[i].ID] = &albums[i]
		}
		for artistRows.Next() {
			var albumID int64
			var artist domain.AlbumArtist
			if err := artistRows.Scan(&albumID, &artist.ID, &artist.Name, &artist.IsPrimary); err != nil {
				return nil, fmt.Errorf("error escaneando artista relacional: %w", err)
			}
			if album, exists := albumMap[albumID]; exists {
				album.Artists = append(album.Artists, artist)
			}
		}
		if err := artistRows.Err(); err != nil {
			return nil, fmt.Errorf("error iterando artistas: %w", err)
		}

	} else {
		albums = []domain.Album{}
	}

	return domain.NewPaginatedResult(albums, totalItems, params.Page, params.Limit), nil
}

func (r *albumRepository) GetAlbumsByArtistID(ctx context.Context, artistID int64) ([]domain.Album, error) {
	queryAlbums := `
		SELECT al.id, al.title, al.release_date, al.type, al.cover_url, al.created_at, al.updated_at
		FROM albums al
		INNER JOIN album_artists aa ON al.album_id = aa.album_id
		WHERE aa.artist_id = $1 AND al.deleted_at IS NULL
		ORDER BY al.release_date DESC
	`
	rows, err := r.db.Query(ctx, queryAlbums, artistID)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando consulta de álbumes por artista: %w", err)
	}
	defer rows.Close()

	var albums []domain.Album
	var albumsIDs []int64
	for rows.Next() {
		var a domain.Album
		if err := rows.Scan(&a.ID, &a.Title, &a.ReleaseDate, &a.Type, &a.CoverURL, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error escaneando álbum de artista: %w", err)
		}
		// definir slices vacios para llenar
		a.Artists = []domain.AlbumArtist{}
		a.Tracks = []domain.Track{}
		albums = append(albums, a)
		albumsIDs = append(albumsIDs, a.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando los álbumes del artista: %w", err)
	}

	// Artista no tiene albums
	if len(albums) == 0 {
		return []domain.Album{}, nil
	}

	// Obtener todos los artistas dela album
	queryArtist := `
		SELECT aa.album_id, a.id, a.name, aa.is_primary
		FROM artist a
		INNER JOIN album_artists aa ON a.id = aa.artist_id
		WHERE aa.album_id = ANY($1) AND a.deleted_at IS NULL
	`
	// Consulta recibe conjunto de IDs
	artistRows, err := r.db.Query(ctx, queryArtist, albumsIDs)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo artistas colaboradores: %w", err)
	}
	defer artistRows.Close()

	albumMap := make(map[int64]*domain.Album) // Mapa, id de album apunta a album
	for i := range albums {
		albumMap[albums[i].ID] = &albums[i]
	}

	for artistRows.Next() {
		var albumID int64
		var artist domain.AlbumArtist
		if err := artistRows.Scan(&albumID, &artist.ID, &artist.Name, &artist.IsPrimary); err != nil {
			return nil, fmt.Errorf("error escaneando colaborador: %w", err)
		}
		if album, exists := albumMap[albumID]; exists {
			album.Artists = append(album.Artists, artist)
		}
	}

	return albums, nil
}

// Editar Albumn con artistas (no tracks)
// PUT /albums/{id}
func (r *albumRepository) Update(ctx context.Context, albumID int64, input *domain.AlbumInput) (*domain.Album, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error iniciando transacción para update de álbum: %w", err)
	}
	defer tx.Rollback(ctx)

	updateAlbumQuery := `
		UPDATE albums
		SET title = $1, release_date = $2, type = $3, cover_url = $4, updated_at = NOW()
		WHERE id = $5 AND deleted_at IS NULL
	`
	res, err := tx.Exec(ctx, updateAlbumQuery, input.Title, input.ReleaseDate, input.Type, input.CoverURL, albumID)
	if err != nil {
		return nil, fmt.Errorf("error actualizando los datos del album: %w", err)
	}
	if res.RowsAffected() == 0 { // Si no afecta filas es porque la cancion no existe o fue borrada
		return nil, errors.New("álbum no encontrado")
	}

	// Limpiar relaciones antiguas
	deletedArtistsQuery := `DELETE FROM album_artist WHERE album_id = $1`
	_, err = tx.Exec(ctx, deletedArtistsQuery, albumID)
	if err != nil {
		return nil, fmt.Errorf("error limpiando relaciones antiguas de artitas del álbum: %w", err)
	}
	deletedTracksQuery := `DELETE FROM tracks WHERE album_id = $1`
	_, err = tx.Exec(ctx, deletedTracksQuery, albumID)
	if err != nil {
		return nil, fmt.Errorf("error limpiando relaciones antiguas de tracks del álbum: %w", err)
	}

	// Insertar nuevas relaciones
	if len(input.Artists) > 0 {
		insertArtistQuery := `
			INSERT INTO album_artists (album_id, artist_id, is_primary)
			VALUES ($1, $2, $3)
		`
		for _, a := range input.Artists {
			_, err := tx.Exec(ctx, insertArtistQuery, albumID, a.ArtistID, a.IsPrimary)
			if err != nil {
				return nil, fmt.Errorf("error insertando nueva relación con artista ID %d: %w", a.ArtistID, err)
			}
		}
	}
	if len(input.Tracks) > 0 {
		insertTracksQuery := `
			INSERT INTO tracks (album_id, song_id, track_number)
			VALUES ($1, $2, $3)
		`
		for _, t := range input.Tracks {
			_, err := tx.Exec(ctx, insertTracksQuery, albumID, t.SongID, t.TrackNumber)
			if err != nil {
				return nil, fmt.Errorf("error insertando nueva relación con concion ID %d: %w", t.SongID, err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error confirmando la transacción de actualización: %w", err)
	}

	updatedAlbum, err := r.GetByID(ctx, albumID)
	if err != nil {
		return nil, fmt.Errorf("álbum actualizado, pero error al obtener detalles: %w", err)
	}
	return updatedAlbum, nil
}

func (r *albumRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE albums SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error eliminando al álbum ID %d: %w", id, err)
	}
	if res.RowsAffected() == 0 {
		return errors.New("álbum no encontrada")
	}
	return nil
}

// (Opcional) Reordenar Tracks (PUT /albums/{id}/tracks)
