package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/IsaacEspinoza91/Song-Manager/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type albumRepository struct {
	db *pgxpool.Pool
}

func NewAlbumRepository(db *pgxpool.Pool) domain.AlbumRepository {
	return &albumRepository{db: db}
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
