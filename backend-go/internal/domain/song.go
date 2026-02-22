package domain

import (
	"time"
)

type Song struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	Duration  int        `json:"duration"` // En segundos
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	Artists []Artist `json:"artists,omitempty"` // Para relationships
}

// Tabla intermedia entre songs-artists. Se puede trabajar desde el input
type SongArtist struct {
	SongID   int64
	ArtistID int64
	Role     string
}
