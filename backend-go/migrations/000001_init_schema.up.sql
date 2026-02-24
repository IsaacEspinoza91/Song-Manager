-- Extension Búsqueda Difusa (Fuzzy Search) soporta errores ortograficos
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- 1. Tabla Artists
CREATE TABLE IF NOT EXISTS artists (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    genre VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL,
    bio TEXT,
    image_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- 2. Tabla Songs
CREATE TABLE IF NOT EXISTS songs (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    duration INT NOT NULL CHECK (duration > 0), -- Validación simple
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- 3. Tabla Albums
CREATE TABLE IF NOT EXISTS albums (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    release_date DATE NOT NULL,
    type VARCHAR(50) NOT NULL,
    cover_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- 4. Tabla Intermedia: Song Artists (Roles)
CREATE TABLE IF NOT EXISTS song_artists (
    song_id BIGINT NOT NULL,
    artist_id BIGINT NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'main',
    
    PRIMARY KEY (song_id, artist_id),
    CONSTRAINT fk_song FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE,
    CONSTRAINT fk_artist FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE
);

-- 5. Tabla Intermedia: Album Artists (Colaboraciones en Albums)
CREATE TABLE IF NOT EXISTS album_artists (
    album_id BIGINT NOT NULL,
    artist_id BIGINT NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT true,

    PRIMARY KEY (album_id, artist_id),
    CONSTRAINT fk_album FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE,
    CONSTRAINT fk_artist FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE
);

-- 6. Tabla Intermedia: Tracks (Relación Album - Song)
CREATE TABLE IF NOT EXISTS tracks (
    album_id BIGINT NOT NULL,
    song_id BIGINT NOT NULL,
    track_number INT NOT NULL,

    PRIMARY KEY (album_id, song_id),
    CONSTRAINT fk_album FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE,
    CONSTRAINT fk_song FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE,
    
    -- Restricción única: No puede haber dos canciones con el track_number 1 en el mismo album
    CONSTRAINT unique_track_number_per_album UNIQUE (album_id, track_number)
);


-- Indices para mejorar performance en búsquedas frecuentes
CREATE INDEX idx_artists_deleted_at ON artists(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_songs_deleted_at ON songs(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_albums_deleted_at ON albums(deleted_at) WHERE deleted_at IS NULL;
-- Busquedas con errores ortograficos mas rapidas. GIN Indice invertido
CREATE INDEX songs_title_trgm_idx ON songs USING GIN (title gin_trgm_ops);