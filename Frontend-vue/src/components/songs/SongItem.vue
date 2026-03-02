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

const displayArtists = computed(() => {
  if (!props.song.artists || props.song.artists.length === 0) return 'Artista Desconocido';
  return props.song.artists.map(a => a.name).join(', ');
});
</script>

<template>
  <div class="song-item glass-panel">
    <div class="number" v-if="index !== undefined">{{ index + 1 }}</div>
    <div class="info">
      <h4 class="title">{{ song.title }}</h4>
      <p class="artists">{{ displayArtists }}</p>
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

.icon-btn.danger:hover {
  /* No specific color as it's an emoji but we do the scale */
}
</style>
