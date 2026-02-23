package service

import (
	"context"
	"errors"

	"github.com/IsaacEspinoza91/Song-Manager/internal/domain"
)

type songService struct {
	repo domain.SongRepository
}

func NewSongService(repo domain.SongRepository) domain.SongService {
	return &songService{repo: repo}
}

func (s *songService) Create(ctx context.Context, input *domain.SongInput) (*domain.Song, error) {
	if err := input.Validate(); err != nil {
		return nil, err // Devolver error de validacion
	}

	song, err := s.repo.Create(ctx, input)
	if err != nil {
		return nil, err
	}
	return song, nil
}

func (s *songService) GetByID(ctx context.Context, id int64) (*domain.Song, error) {
	if id <= 0 {
		return nil, errors.New("ID de cancion inválido")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *songService) GetAll(ctx context.Context) ([]domain.Song, error) {
	return s.repo.GetAll(ctx)
}

func (s *songService) GetAllPaginated(ctx context.Context, filter domain.SongFilter, params domain.PaginationParams) (*domain.PaginatedResult[domain.Song], error) {
	// Dentro de pagination.go se hacen validaciones de page y limit
	return s.repo.GetAllPaginated(ctx, filter, params)
}

func (s *songService) Update(ctx context.Context, id int64, input *domain.SongInput) (*domain.Song, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	if id <= 0 {
		return nil, errors.New("ID de cancion inválido")
	}
	song, err := s.repo.Update(ctx, id, input)
	if err != nil {
		return nil, err
	}
	return song, nil
}

func (s *songService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("ID de cancion inválido")
	}
	return s.repo.Delete(ctx, id)
}
