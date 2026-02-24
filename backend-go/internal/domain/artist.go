package domain

import (
	"context"
	"time"

	"github.com/IsaacEspinoza91/Song-Manager/pkg/validation"
)

// MODELOS

type Artist struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Genre     string     `json:"genre"`
	Country   string     `json:"country"`
	Bio       *string    `json:"bio"` // Puntero puede ser nulo
	ImageURL  *string    `json:"image_url"`
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

// Contiene los campos opcionales para buscar artistas
type ArtistFilter struct {
	Name    string
	Genre   string
	Country string
}

// VALIDACIONES Y LIMPIEZA
func (input *ArtistFilter) Sanitize() {
	input.Name = validation.SanitizeString(input.Name)
	input.Genre = validation.SanitizeString(input.Genre)
	input.Country = validation.SanitizeString(input.Country)
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

// ValidationError mapa personalizado para acumular errores por campo
type ValidationError map[string]string

// Implementamos la interfaz 'error' nativa de Go
// para que ValidationError pueda ser retornado como un error normal
func (e ValidationError) Error() string {
	return "errores de validación en los datos de entrada"
}

func (input *ArtistInput) Validate() error {
	input.Sanitize()
	errs := make(ValidationError) // Crea mapa para acum errores

	if input.Name == "" {
		errs["name"] = "el nombre es obligatorio"
	}
	if input.Genre == "" {
		errs["genre"] = "el género es obligatorio"
	}
	if input.Country == "" {
		errs["country"] = "el país es obligatorio"
	}

	// Si el mapa tiene elementos, significa que hubo errores
	if len(errs) > 0 {
		return errs // Retornamos el mapa de errores
	}
	return nil
}

// INTERFACES

type ArtistRepository interface {
	Create(ctx context.Context, input *ArtistInput) (*Artist, error)
	Update(ctx context.Context, id int64, input *ArtistInput) (*Artist, error)
	GetAll(ctx context.Context) ([]Artist, error)
	GetAllPaginated(ctx context.Context, filter ArtistFilter, params PaginationParams) (*PaginatedResult[Artist], error)
	GetByID(ctx context.Context, id int64) (*Artist, error)
	Delete(ctx context.Context, id int64) error
}

type ArtistService interface {
	Create(ctx context.Context, input *ArtistInput) (*Artist, error)
	Update(ctx context.Context, id int64, input *ArtistInput) (*Artist, error)
	GetAll(ctx context.Context) ([]Artist, error)
	GetAllPaginated(ctx context.Context, filter ArtistFilter, params PaginationParams) (*PaginatedResult[Artist], error)
	GetByID(ctx context.Context, id int64) (*Artist, error)
	Delete(ctx context.Context, id int64) error
}
