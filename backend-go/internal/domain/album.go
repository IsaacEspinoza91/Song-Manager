package domain

import (
	"context"
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

func (input *AlbumInput) Sanitize() {
	input.Title = validation.SanitizeString(input.Title)
	input.ReleaseDate = validation.SanitizeString(input.ReleaseDate)
	input.Type = validation.SanitizeString(input.Type)
	input.CoverURL = validation.SanitizeOpcionalString(input.CoverURL)
}

// Validated implementa que no haya track_numbers duplicados

// INTERFACES
type AlbumRepository interface {
	GetByID(ctx context.Context, id int64) (*Album, error)
}
