package domain

import "math"

// PaginationParams define lo que entra desde la URL (?page=1&limit=10)
type PaginationParams struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// GetOffset calcula el salto de registros para la base de datos
// Ej: en la pagina 2, deberia omitir los 10 primeros elementos (offset 10)
func (p *PaginationParams) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10 // Límite por defecto
	}
	return (p.Page - 1) * p.Limit
}

// PaginatedResult es una estructura genérica
// [T any] que puede contener cualquier slice de modelos
type PaginatedResult[T any] struct {
	Data       []T `json:"data"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
}

// Helper para calcular las páginas totales automáticamente
func NewPaginatedResult[T any](data []T, totalItems, page, limit int) *PaginatedResult[T] {
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	return &PaginatedResult[T]{
		Data:       data,
		TotalItems: totalItems,
		TotalPages: totalPages,
		Page:       page,
		Limit:      limit,
	}
}
