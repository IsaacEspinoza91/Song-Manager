package domain

import "errors"

// Errores comunes y transversales
var (
	ErrInvalidID = errors.New("el ID proporcionado es inválido")
)

// Errores de Artistas
var (
	ErrArtistNotFound  = errors.New("artista no encontrado")
	ErrArtistIDInvalid = errors.New("ID de artista inválido")
)

// Errores de Canciones
var (
	ErrSongNotFound  = errors.New("canción no encontrada")
	ErrSongIDInvalid = errors.New("ID de canción inválido")
)

// Errores de Álbumes y Tracks
var (
	ErrAlbumNotFound      = errors.New("álbum no encontrado")
	ErrAlbumIDInvalid     = errors.New("ID de álbum inválido")
	ErrTrackNotFound      = errors.New("track no encontrado en este álbum")
	ErrTrackAlreadyExists = errors.New("este número de pista ya está ocupado en el álbum")
	ErrSongAlreadyInAlbum = errors.New("esta canción ya existe en este álbum")
	ErrSongNotInDB        = errors.New("la canción indicada no existe en la base de datos")
)
