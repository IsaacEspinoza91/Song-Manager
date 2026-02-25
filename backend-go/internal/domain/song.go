package domain

import (
	"context"
	"time"

	"github.com/IsaacEspinoza91/Song-Manager/pkg/validation"
)

// MODELOS
// ArtistWithRole datos de artistas para response
type ArtistWithRole struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type Song struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	Duration  int        `json:"duration"` // En segundos
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	Artists []ArtistWithRole `json:"artists,omitempty"`
}

type ArtistSongInput struct {
	ArtistID int64  `json:"artist_id"`
	Role     string `json:"role"`
}

type SongInput struct {
	Title    string            `json:"title"`
	Duration int               `json:"duration"`
	Artists  []ArtistSongInput `json:"artists"` // si frontend no lo manda, asigna array default limpio
}

type SongFilter struct {
	Title      string
	ArtistID   int64
	ArtistName string
}

// VALIDACIONES
func (input *ArtistSongInput) Sanitize() {
	input.Role = validation.SanitizeString(input.Role)
}

func (input *ArtistSongInput) Validate() error {
	input.Sanitize()
	errs := make(ValidationError)

	if input.ArtistID <= 0 {
		errs["artist"] = "el artista tiene un ID inválido"
	}
	if input.Role != "producer" && input.Role != "ft" && input.Role != "main" {
		errs["role"] = "el rol de artistas debe ser main, ft o producer"
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (input *SongFilter) Sanitize() {
	input.Title = validation.SanitizeString(input.Title)
	input.ArtistName = validation.SanitizeString(input.ArtistName)
}

func (input *SongInput) Sanitize() {
	input.Title = validation.SanitizeString(input.Title)
}

// Uso de mapa ValidationError definido en artist.go
func (input *SongInput) Validate() error {
	input.Sanitize()
	errs := make(ValidationError)

	if input.Title == "" {
		errs["title"] = "el nombre es obligatorio"
	}
	if input.Duration <= 0 {
		errs["duration"] = "la duración debe ser mayor a 0 segundos"
	}
	if len(input.Artists) > 0 {
		for _, a := range input.Artists {
			if a.ArtistID <= 0 {
				errs["artists"] = "uno de los artistas tiene un ID inválido"
				break // Salimos del bucle para no saturar con errores
			}

			if a.Role != "main" && a.Role != "ft" && a.Role != "producer" {
				errs["role"] = "el rol de artistas debe ser main, ft o producer"
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// INTERFACES
type SongRepository interface {
	Create(ctx context.Context, input *SongInput) (*Song, error)
	GetByID(ctx context.Context, id int64) (*Song, error)
	GetAll(ctx context.Context) ([]Song, error)
	GetAllPaginated(ctx context.Context, filter SongFilter, params PaginationParams) (*PaginatedResult[Song], error)
	Update(ctx context.Context, id int64, input *SongInput) (*Song, error)
	Delete(ctx context.Context, id int64) error
	AddArtist(ctx context.Context, songID int64, input *ArtistSongInput) error
	RemoveArtist(ctx context.Context, songID, artistID int64) error
}

type SongService interface {
	Create(ctx context.Context, input *SongInput) (*Song, error)
	GetByID(ctx context.Context, id int64) (*Song, error)
	GetAll(ctx context.Context) ([]Song, error)
	GetAllPaginated(ctx context.Context, filter SongFilter, params PaginationParams) (*PaginatedResult[Song], error)
	Update(ctx context.Context, id int64, input *SongInput) (*Song, error)
	Delete(ctx context.Context, id int64) error
	AddArtist(ctx context.Context, songID int64, input *ArtistSongInput) error
	RemoveArtist(ctx context.Context, songID, artistID int64) error
}
