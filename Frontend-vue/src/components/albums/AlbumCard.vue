<script setup>
import { defineProps, computed } from 'vue';
import Icon from '../common/Icon.vue';

const props = defineProps({
  album: {
    type: Object,
    required: true
  },
  readonly: {
    type: Boolean,
    default: false
  }
});

const releaseYear = computed(() => {
  if (!props.album.release_date) return 'Desconocido';
  return new Date(props.album.release_date).getFullYear();
});

const primaryArtistsList = computed(() => {
  if (!props.album.artists || props.album.artists.length === 0) return [];
  // Get all primary artists, or default to all artists if none are marked primary
  const primaries = props.album.artists.filter(a => a.is_primary);
  return primaries.length > 0 ? primaries : props.album.artists;
});
</script>

<template>
  <div class="album-card glass-panel" @click="$router.push(`/albums/${album.id}`)">
    <div class="cover-wrapper">
      <img v-if="album.cover_url" :src="album.cover_url" :alt="album.title" class="cover-image" />
      <div v-else class="placeholder-cover">
        <span>{{ album.title.charAt(0) }}</span>
      </div>
      <div class="type-badge">{{ album.type }}</div>
    </div>
    <div class="details">
      <h3 class="title">{{ album.title }}</h3>
      <p class="artist">
          <span v-if="!primaryArtistsList.length">Artista Desconocido</span>
          <template v-else>
              <template v-for="(artist, index) in primaryArtistsList" :key="artist.id || artist.artist_id || index">
                  <router-link :to="`/artists/${artist.id || artist.artist_id}`" class="artist-link" @click.stop>
                      {{ artist.name || artist.artist_name }}
                  </router-link><span v-if="index < primaryArtistsList.length - 1">, </span>
              </template>
          </template>
      </p>
      <div class="meta">
        <span class="year">{{ releaseYear }}</span>
        <div class="actions" v-if="!readonly">
          <button class="icon-btn" @click.stop="$emit('edit', album)" title="Editar Álbum">
             <Icon name="edit" />
          </button>
          <button class="icon-btn danger" @click.stop="$emit('delete', album.id)" title="Eliminar Álbum">
             <Icon name="trash" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.album-card {
  display: flex;
  flex-direction: column;
  border-radius: var(--radius-lg);
  overflow: hidden;
  transition: all var(--transition-normal);
  cursor: pointer;
}

.album-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 10px 20px rgba(0,0,0,0.2), 0 0 15px rgba(6, 182, 212, 0.3); /* Cyan shadow hint */
  border-color: rgba(255,255,255,0.2);
}

.cover-wrapper {
  width: 100%;
  aspect-ratio: 1;
  position: relative;
  overflow: hidden;
}

.cover-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform var(--transition-normal);
}

.album-card:hover .cover-image {
  transform: scale(1.03);
}

.placeholder-cover {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #1f2937, #111827);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 3rem;
  color: var(--text-muted);
  font-weight: bold;
}

.type-badge {
  position: absolute;
  top: 0.75rem;
  right: 0.75rem;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  padding: 0.25rem 0.5rem;
  border-radius: var(--radius-sm);
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: white;
  border: 1px solid rgba(255,255,255,0.1);
}

.details {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  flex: 1;
}

.title {
  font-size: 1.1rem;
  font-weight: 600;
  margin-bottom: 0.2rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.artist {
  font-size: 0.9rem;
  color: var(--text-secondary);
  margin-bottom: 0.75rem;
}

.artist-link {
  color: var(--text-secondary);
  text-decoration: none;
  transition: color var(--transition-fast);
}

.artist-link:hover {
  text-decoration: underline;
  color: var(--text-primary);
}

.meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: auto;
}

.year {
  font-size: 0.8rem;
  color: var(--text-muted);
}

.actions {
  display: flex;
  gap: 0.25rem;
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.album-card:hover .actions {
  opacity: 1;
}

.icon-btn {
  padding: 0.25rem;
  background: none;
  transition: transform var(--transition-fast);
}

.icon-btn:hover {
  transform: scale(1.2);
}
</style>
