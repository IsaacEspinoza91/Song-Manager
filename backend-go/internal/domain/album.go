package domain

import (
	"time"
)

type Album struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	ReleaseDate time.Time  `json:"release_date"`
	Type        string     `json:"type"` // e.g. 'EP', 'LP', 'Single'
	CoverURL    *string    `json:"cover_url,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`

	Artists []Artist `json:"artists,omitempty"` // Artists mapped from album_artists
	Tracks  []Track  `json:"tracks,omitempty"`
}

// Tabla intermedia Albums-Songs
type Track struct {
	AlbumID     int64 `json:"album_id"`
	SongID      int64 `json:"song_id"`
	TrackNumber int   `json:"track_number"`

	Song *Song `json:"song,omitempty"`
}

// Tabla intermedia Albums-Artists. Puede trabajarse en input
type AlbumArtist struct {
	AlbumID   int64 `json:"album_id"`
	ArtistID  int64 `json:"artist_id"`
	IsPrimary bool  `json:"is_primary"`
}
