package domain

import "errors"

// Errores comunes y transversales
var (
	ErrInvalidID = errors.New("el ID proporcionado es inválido")
)

// Errores de Artistas
var (
	ErrArtistNotFound = errors.New("artista no encontrado")
)

// Errores de Canciones
var (
	ErrSongNotFound = errors.New("canción no encontrada")
)

// Errores de Álbumes y Tracks
var (
	ErrAlbumNotFound      = errors.New("álbum no encontrado")
	ErrTrackNotFound      = errors.New("track no encontrado en este álbum")
	ErrTrackAlreadyExists = errors.New("este número de pista ya está ocupado en el álbum")
	ErrSongAlreadyInAlbum = errors.New("esta canción ya existe en este álbum")
	ErrSongNotInDB        = errors.New("la canción indicada no existe en la base de datos")
)