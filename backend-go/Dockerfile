# Etapa 1: Construcción (Builder)
# Imagen oficial de Go basada en Alpine (muy ligera)
FROM golang:1.22-alpine AS builder

# Seteamos el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiamos los archivos de dependencias primero (para aprovechar la caché de Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto del código fuente
COPY . .

# Compilamos la aplicación. 
# CGO_ENABLED=0 asegura que el binario sea 100% estático y funcione en cualquier Linux.
RUN CGO_ENABLED=0 GOOS=linux go build -o song-manager .



# Etapa 2: Producción (Release)
# Usamos Alpine puro
FROM alpine:latest

# Añadimos certificados de seguridad (necesarios si API hace peticiones HTTPS a terceros)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copiamos SOLO el archivo binario compilado desde la Etapa 1
COPY --from=builder /app/song-manager .

# Exponemos el puerto 
EXPOSE 8080

# Comando para arrancar el servidor
CMD ["./song-manager"]