package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/IsaacEspinoza91/Song-Manager/pkg/validation"
)

// MODELOS

// Representa a artista dentro de album
type AlbumArtist struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	IsPrimary bool   `json:"is_primary"`
}

// Representa a cancion dentro de album
type Track struct {
	TrackNumber int    `json:"track_number"`
	SongID      int64  `json:"song_id"`
	Title       string `json:"title"`    // Info extraída de la tabla songs
	Duration    int    `json:"duration"` // Info extraída de la tabla songs
}

type Album struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	ReleaseDate time.Time  `json:"release_date"`
	Type        string     `json:"type"` // e.g. 'EP', 'LP', 'Single'
	CoverURL    *string    `json:"cover_url"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`

	Artists []AlbumArtist `json:"artists,omitempty"` // Mapeados para dto respuesta
	Tracks  []Track       `json:"tracks,omitempty"`
}

type AlbumArtistInput struct {
	ArtistID  int64 `json:"artist_id"`
	IsPrimary bool  `json:"is_primary"`
}

type TrackInput struct {
	SongID      int64 `json:"song_id"`
	TrackNumber int   `json:"track_number"`
}

type AlbumInput struct {
	Title       string             `json:"title"`
	ReleaseDate string             `json:"release_date"` // Formato "YYYY-MM-DD"
	Type        string             `json:"type"`
	CoverURL    *string            `json:"cover_url"`
	Artists     []AlbumArtistInput `json:"artists"`
	Tracks      []TrackInput       `json:"tracks"`
}

type AlbumFilter struct {
	Title      string
	Type       string
	ArtistID   int64
	ArtistName string
}

func (input *AlbumInput) Sanitize() {
	input.Title = validation.SanitizeString(input.Title)
	input.ReleaseDate = validation.SanitizeString(input.ReleaseDate)
	input.Type = validation.SanitizeString(input.Type)
	input.CoverURL = validation.SanitizeOpcionalString(input.CoverURL)
}

func (input *AlbumInput) Validate() error {
	input.Validate()
	errs := make(ValidationError)

	if input.Title == "" {
		errs["title"] = "el título del álbum es obligatorio"
	}
	if input.Type != "EP" && input.Type != "LP" && input.Type != "Single" {
		errs["type"] = "el tipo de álbum debe ser EP, LP o Single"
	}
	if input.ReleaseDate == "" {
		errs["release_date"] = "la fecha de lanzamiento es obligatoria"
	} else {
		// time.DateOnly equivale internamente a "2006-01-02"
		_, err := time.Parse(time.DateOnly, input.ReleaseDate)
		if err != nil {
			errs["release_date"] = "el formato de la fecha debe ser YYYY-MM-DD"
		}
	}

	// Validacion negocio
	if len(input.Artists) == 0 {
		errs["artists"] = "el álbum debe tener al menos un artista asociado"
	} else {
		// Al menos un artista principal
		hasPrimary := false
		for _, a := range input.Artists {
			if a.ArtistID <= 0 {
				errs["artists"] = "uno de los artistas tiene un ID inválido"
				break
			}
			if a.IsPrimary {
				hasPrimary = true
			}
		}

		if !hasPrimary && errs["artists"] == "" {
			errs["artists"] = "el álbum debe tener al menos un artista marcado como principal (is_primary: true)"
		}
	}

	if len(input.Tracks) > 0 {
		seenTrackNumbers := make(map[int]bool)
		seenSongIDs := make(map[int64]bool) // Evitar que manden la misma canción 2 veces

		for _, t := range input.Tracks {
			if t.TrackNumber <= 0 {
				errs["tracks"] = "los números de pista deben ser mayores a 0"
				break
			}
			if t.SongID <= 0 {
				errs["tracks"] = "uno de los IDs de canción es inválido"
				break
			}

			// Buscar num de track duplicados
			if seenTrackNumbers[t.TrackNumber] {
				errs["tracks"] = fmt.Sprintf("el número de pista %d está duplicado", t.TrackNumber)
				break
			}
			seenTrackNumbers[t.TrackNumber] = true // Add Track al mapa

			// Buscar id de song
			if seenSongIDs[t.SongID] {
				errs["tracks"] = fmt.Sprintf("la canción con ID %d está duplicada en el tracklist", t.SongID)
				break
			}
			seenSongIDs[t.SongID] = true
		}
	}

	// Si hay errores, retornamos el mapa
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (input *TrackInput) Validate() error {
	errs := make(ValidationError)
	if input.SongID <= 0 {
		errs["song_id"] = "el ID de la canción es obligatorio y debe ser mayor a 0"
	}
	if input.TrackNumber <= 0 {
		errs["track_number"] = "el número de pista debe ser mayor a 0"
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// INTERFACES
type AlbumRepository interface {
	Create(ctx context.Context, input *AlbumInput) (*Album, error)
	AddTrack(ctx context.Context, albumID int64, input *TrackInput) error
	RemoveTrack(ctx context.Context, albumID int64, songID int64) error
	GetByID(ctx context.Context, id int64) (*Album, error)
	GetAlbumsByArtistID(ctx context.Context, artistID int64) ([]Album, error)
	GetAllPaginated(ctx context.Context, filter AlbumFilter, params PaginationParams) (*PaginatedResult[Album], error)
	Update(ctx context.Context, albumID int64, input *AlbumInput) (*Album, error)
	Delete(ctx context.Context, id int64) error
}

type AlbumService interface {
	Create(ctx context.Context, input *AlbumInput) (*Album, error)
	GetByID(ctx context.Context, albumID int64) (*Album, error)
	GetAllPaginated(ctx context.Context, filter AlbumFilter, params PaginationParams) (*PaginatedResult[Album], error)
	GetAlbumsByArtistID(ctx context.Context, artistID int64) ([]Album, error)
	Update(ctx context.Context, id int64, input *AlbumInput) (*Album, error)
	AddTrack(ctx context.Context, albumID int64, input *TrackInput) error
	RemoveTrack(ctx context.Context, albumID int64, songID int64) error
	Delete(ctx context.Context, albumID int64) error
}
