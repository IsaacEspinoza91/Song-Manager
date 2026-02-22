package validation

import "strings"

// Quita espacios vacios finales y iniciales trim
func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}

// Trim para string opcional. Solo si es "" se cambia a nil
func SanitizeOpcionalString(s *string) *string {
	if s == nil {
		return nil
	}

	clean := strings.TrimSpace(*s)
	if clean == "" {
		return nil
	}

	return &clean
}
