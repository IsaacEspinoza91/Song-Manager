package service

import (
	"context"
	"errors"

	"github.com/IsaacEspinoza91/Song-Manager/internal/domain"
)

type albumService struct {
	repo domain.AlbumRepository
}

func NewAlbumService(repo domain.AlbumRepository) domain.AlbumService {
	return &albumService{repo: repo}
}

// CREATE
func (s *albumService) Create(ctx context.Context, input *domain.AlbumInput) (*domain.Album, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, input)
}

// READ
func (s *albumService) GetByID(ctx context.Context, albumID int64) (*domain.Album, error) {
	if albumID <= 0 {
		return nil, domain.ErrAlbumIDInvalid
	}
	return s.repo.GetByID(ctx, albumID)
}

func (s *albumService) GetAllPaginated(ctx context.Context, filter domain.AlbumFilter, params domain.PaginationParams) (*domain.PaginatedResult[domain.Album], error) {
	return s.repo.GetAllPaginated(ctx, filter, params)
}

func (s *albumService) GetAlbumsByArtistID(ctx context.Context, artistID int64) ([]domain.Album, error) {
	if artistID <= 0 {
		return nil, domain.ErrArtistIDInvalid
	}
	return s.repo.GetAlbumsByArtistID(ctx, artistID)
}

// UPDATE
func (s *albumService) Update(ctx context.Context, id int64, input *domain.AlbumInput) (*domain.Album, error) {
	// Validaciones defensivas
	if id <= 0 {
		return nil, domain.ErrAlbumIDInvalid
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}

	return s.repo.Update(ctx, id, input)
}

func (s *albumService) AddTrack(ctx context.Context, albumID int64, input *domain.TrackInput) error {
	// Validaciones defensivas del ID y los datos de entrada
	if albumID <= 0 {
		return domain.ErrAlbumIDInvalid
	}
	if err := input.Validate(); err != nil {
		return err // Mapa de ValidationError
	}

	// En caso de error se retorna, sino hay error se retorna nil
	return s.repo.AddTrack(ctx, albumID, input)
}

func (s *albumService) RemoveTrack(ctx context.Context, albumID int64, songID int64) error {
	if albumID <= 0 || songID <= 0 {
		return errors.New("IDs inválidos")
	}

	return s.repo.RemoveTrack(ctx, albumID, songID)
}

// DELETE
func (s *albumService) Delete(ctx context.Context, albumID int64) error {
	if albumID <= 0 {
		return domain.ErrAlbumIDInvalid
	}
	return s.repo.Delete(ctx, albumID)
}
