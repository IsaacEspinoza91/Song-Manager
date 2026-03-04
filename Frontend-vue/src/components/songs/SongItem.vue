<script setup>
import { computed } from 'vue';

const props = defineProps({
  song: {
    type: Object,
    required: true
  },
  index: {
    type: Number,
    required: false
  },
  readonly: {
    type: Boolean,
    default: false
  }
});

const formatDuration = (seconds) => {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${secs.toString().padStart(2, '0')}`;
};

const mainArtists = computed(() => {
  if (!props.song.artists) return [];
  return props.song.artists.filter(a => a.role === 'main' || !a.role);
});

const ftArtists = computed(() => {
  if (!props.song.artists) return [];
  return props.song.artists.filter(a => a.role === 'ft');
});
</script>

<template>
  <div class="song-item glass-panel">
    <div class="number" v-if="index !== undefined">{{ index + 1 }}</div>
    <img v-if="song.cover_url" :src="song.cover_url" alt="Cover" class="album-cover" />
    <div class="info">
      <h4 class="title">{{ song.title }}</h4>
      <p class="artists">
        <span v-if="mainArtists.length === 0 && ftArtists.length === 0">Artista Desconocido</span>
        
        <span v-for="(artist, idx) in mainArtists" :key="artist.id">
            <router-link :to="`/artists/${artist.id}`" class="artist-link" @click.stop>{{ artist.name }}</router-link><span v-if="idx < mainArtists.length - 1">, </span>
        </span>
        
        <span v-if="ftArtists.length > 0">
            <span class="ft-text"> ft. </span>
            <span v-for="(artist, idx) in ftArtists" :key="'ft-'+artist.id">
                <router-link :to="`/artists/${artist.id}`" class="artist-link" @click.stop>{{ artist.name }}</router-link><span v-if="idx < ftArtists.length - 1">, </span>
            </span>
        </span>
      </p>
    </div>
    <div class="duration">{{ formatDuration(song.duration) }}</div>
    <div class="actions" v-if="!readonly">
       <button class="icon-btn" @click="$emit('edit', song)" title="Edit">✏️</button>
       <button class="icon-btn danger" @click="$emit('delete', song.id)" title="Delete">🗑️</button>
    </div>
  </div>
</template>

<style scoped>
.song-item {
  display: flex;
  align-items: center;
  padding: 1rem 1.5rem;
  margin-bottom: 0.5rem;
  transition: all var(--transition-fast);
}

.song-item:hover {
  background: rgba(255, 255, 255, 0.05);
  transform: translateX(4px);
  border-color: rgba(255, 255, 255, 0.2);
}

.number {
  width: 30px;
  color: var(--text-muted);
  font-variant-numeric: tabular-nums;
  font-size: 0.9rem;
}

.info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.title {
  font-size: 1rem;
  font-weight: 500;
  margin-bottom: 0.1rem;
}

.album-cover {
  width: 38px;
  height: 38px;
  border-radius: 0;
  object-fit: cover;
  margin-right: 1rem;
}

.artists {
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.duration {
  font-variant-numeric: tabular-nums;
  color: var(--text-secondary);
  font-size: 0.9rem;
  margin: 0 2rem;
}

.actions {
  display: flex;
  gap: 0.5rem;
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.song-item:hover .actions {
  opacity: 1;
}

.icon-btn {
  background: none;
  border: none;
  font-size: 1.1rem;
  cursor: pointer;
  padding: 0.25rem;
  border-radius: var(--radius-sm);
  transition: transform var(--transition-fast);
}

.icon-btn:hover {
  transform: scale(1.2);
}

.artist-link {
  color: var(--text-secondary);
  text-decoration: none;
  transition: color var(--transition-fast);
}

.artist-link:hover {
  color: var(--text-primary);
  text-decoration: underline;
}

.ft-text {
  color: var(--text-muted);
  font-size: 0.8rem;
  margin: 0 0.1rem;
}
</style>
