package middleware

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

/*
Rate Limiting. Limita las peticiones por IP en un tiempo.
Evita que usuario mande 10.000 peticiones por seg y ataque el sistema.

Este rate limiting solo funciona para un molitico de 1 replica.
En caso de replicas (escalado horizontal - load balancer), no funciona porque
no sabe si las peticiones van a un server u otro por lo que las cuenta mal. (usar Redis)
*/

// getVisitor obtiene o crea un limitador para una IP específica
func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		// r = Tasa de recarga (ej. 2 peticiones por segundo)
		// b = Ráfaga máxima permitida de golpe (ej. 5 peticiones)
		limiter = rate.NewLimiter(2, 5)
		visitors[ip] = limiter
	}

	return limiter
}

// RateLimit intercepta las peticiones y bloquea las que superen el límite
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener la IP real del cliente
		ip := extractIP(r)

		// Obtener el limitador para esa IP
		limiter := getVisitor(ip)

		// Evaluar si se le permite pasar
		if !limiter.Allow() {
			log.Printf("[RATE LIMIT] IP bloqueada temporalmente: %s\n", ip)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests) // Código HTTP 429

			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  http.StatusTooManyRequests,
				"message": "Has superado el límite de peticiones. Por favor, intenta más tarde.",
			})
			return
		}

		next.ServeHTTP(w, r) // Si tiene permisos, la petición continúa
	})
}

// extractIP limpia la IP
func extractIP(r *http.Request) string {
	// Si la API está detrás de un Nginx o Load Balancer (común en producción)
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		return forwardedFor
	}

	// Si es conexión directa, r.RemoteAddr viene como "IP:Puerto" (ej. "192.168.1.5:45321")
	ip, _, err := net.SplitHostPort(r.RemoteAddr) // SplitHostPort para obtener solo con la IP
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
