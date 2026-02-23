package domain

import (
	"context"
	"errors"
	"time"

	"github.com/IsaacEspinoza91/Song-Manager/pkg/validation"
)

type Artist struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Genre     string     `json:"genre"`
	Country   string     `json:"country"`
	Bio       *string    `json:"bio,omitempty"` // Puntero puede ser nulo
	ImageURL  *string    `json:"image_url,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	Role      string `json:"role,omitempty"`       // Song relationships
	IsPrimary *bool  `json:"is_primary,omitempty"` // Album relationships
}

type ArtistInput struct {
	Name     string  `json:"name"`
	Genre    string  `json:"genre"`
	Country  string  `json:"country"`
	Bio      *string `json:"bio"`
	ImageURL *string `json:"image_url"`
}

type ArtistRepository interface {
	Create(ctx context.Context, input *ArtistInput) (*Artist, error)
	Update(ctx context.Context, id int64, input *ArtistInput) (*Artist, error)
	GetAll(ctx context.Context) ([]Artist, error)
	GetByID(ctx context.Context, id int64) (*Artist, error)
	Count(ctx context.Context) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type ArtistService interface {
	Create(ctx context.Context, input *ArtistInput) (*Artist, error)
	Update(ctx context.Context, id int64, input *ArtistInput) (*Artist, error)
	GetAll(ctx context.Context) ([]Artist, error)
	GetByID(ctx context.Context, id int64) (*Artist, error)
	Delete(ctx context.Context, id int64) error
}

func (input *ArtistInput) Sanitize() {
	// En caso de ser nulos, go asigna valor default del tipo
	input.Name = validation.SanitizeString(input.Name)
	input.Genre = validation.SanitizeString(input.Genre)
	input.Country = validation.SanitizeString(input.Country)

	// Validacion opcional, si es string vacio se cambia a nil
	input.Bio = validation.SanitizeOpcionalString(input.Bio)
	input.ImageURL = validation.SanitizeOpcionalString(input.ImageURL)
}

func (input *ArtistInput) Validate() error {
	input.Sanitize()
	if input.Name == "" {
		return errors.New("el nombre del artista es obligatorio")
	}
	if input.Genre == "" {
		return errors.New("el género es obligatorio")
	}
	if input.Country == "" {
		return errors.New("el país es obligatorio")
	}
	return nil
}
