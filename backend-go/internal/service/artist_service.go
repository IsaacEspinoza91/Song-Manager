package service

import (
	"context"
	"errors"

	"github.com/IsaacEspinoza91/Song-Manager/internal/domain"
)

type artistService struct {
	repo domain.ArtistRepository
}

// Constructor: Recibe la interfaz del repo y devuelve la interfaz del servicio
func NewArtistService(repo domain.ArtistRepository) domain.ArtistService {
	return &artistService{repo: repo}
}

// 1. Create
func (s *artistService) Create(ctx context.Context, input *domain.ArtistInput) (*domain.Artist, error) {
	// Validate() hace internamente Saniteze()
	if err := input.Validate(); err != nil {
		// Devolver error de validación, Handler debe mostrar 400 Bad Request
		return nil, err
	}

	artist, err := s.repo.Create(ctx, input)
	if err != nil {
		// Create ya retorna ("error creando al artista")
		return nil, err
	}

	return artist, nil
}

// 2. Read
func (s *artistService) GetByID(ctx context.Context, id int64) (*domain.Artist, error) {
	if id <= 0 {
		return nil, errors.New("ID de artista inválido")
	}

	artist, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return artist, nil
}

func (s *artistService) GetAll(ctx context.Context) ([]domain.Artist, error) {
	return s.repo.GetAll(ctx)
}

// 3. Update
func (s *artistService) Update(ctx context.Context, id int64, input *domain.ArtistInput) (*domain.Artist, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	if id <= 0 {
		return nil, errors.New("ID de artista inválido")
	}

	artist, err := s.repo.Update(ctx, id, input)
	if err != nil {
		return nil, err
	}

	return artist, nil
}

// 4. Delete
func (s *artistService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("ID de artista inválido")
	}

	return s.repo.Delete(ctx, id)
}
