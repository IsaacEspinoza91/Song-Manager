package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

/*
Middleware funciona como portero para las peticiones. se encarga de todo aquello que es transversal a la app.
Las responsabilidades tipicas son:   (OJO EL ORDEN)
- Logging de peticiones
- CORS
- Recovery (Anti-Pánico): en caso de panic, middleware lo atrapa y evita una caida del sistema
- Autenticación / Autorización: ej JWT
*/

// Logger envuelve el handler para medir el tiempo y registrar la ruta. Print en consola
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r) // Pasar la petición al siguiente nivel (tu router/handler)

		// Cuando el handler termina, calcular el tiempo que tomó
		duration := time.Since(start)
		log.Printf("[%s] %s - %v\n", r.Method, r.URL.Path, duration)
	})
}

// CORS configura las cabeceras para permitir que un frontend consuma la API
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Cabeceras de permiso
		w.Header().Set("Access-Control-Allow-Origin", "*") // En producción, cambiar "*" por dominio
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 2. Manejo del "Preflight Request"
		// Los navegadores envían una petición OPTIONS antes de un POST/PUT para ver si tienen permiso.
		// Si es OPTIONS, le decimos "sí tienes permiso" y cortamos el flujo aquí (Status 204 No Content).
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// 3. Si no es OPTIONS, dejamos que la petición siga hacia el Handler real
		next.ServeHTTP(w, r)
	})
}

// Recovery (Anti-Pánico) atrapa cualquier panic que ocurra en los handlers y evita que la petición muera
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// defer se ejecuta SIEMPRE justo antes de que la función retorne o muera
		defer func() {
			// recover() detiene el pánico y nos devuelve el error que lo causó
			if err := recover(); err != nil {
				log.Printf("[PANIC RECOVERED] %v\n", err) // Print el error del pánico

				// Print el Stack Trace para saber en qué línea exacta del código explotó
				log.Printf("[STACK TRACE]\n%s\n", debug.Stack())

				// Devolver un error 500 JSON al cliente
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				// Usamos una estructura rápida para el JSON de error
				json.NewEncoder(w).Encode(map[string]interface{}{
					"status":  http.StatusInternalServerError,
					"message": "Ocurrió un error interno crítico en el servidor",
				})
			}
		}()

		// Pasamos la petición al siguiente middleware o handler
		next.ServeHTTP(w, r)
	})
}

// Autenticación / Autorización
