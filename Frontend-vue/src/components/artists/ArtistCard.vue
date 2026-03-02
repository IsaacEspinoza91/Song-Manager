<script setup>
import { defineProps } from 'vue';

const props = defineProps({
  artist: {
    type: Object,
    required: true
  }
});
</script>

<template>
  <div class="artist-card glass-panel" @click="$router.push(`/artists/${artist.id}`)">
    <div class="image-wrapper">
      <img v-if="artist.image_url" :src="artist.image_url" :alt="artist.name" class="artist-image" />
      <div v-else class="placeholder-image">
        <span>{{ artist.name.charAt(0) }}</span>
      </div>
    </div>
    <div class="content">
      <h3 class="name">{{ artist.name }}</h3>
      <p class="genre">{{ artist.genre }} • {{ artist.country }}</p>
      <div class="actions">
        <button class="btn btn-secondary btn-sm" @click.stop="$emit('edit', artist)">✏️ Editar</button>
        <button class="btn btn-danger btn-sm" @click.stop="$emit('delete', artist.id)">🗑️ Eliminar</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.artist-card {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  transition: transform var(--transition-normal), box-shadow var(--transition-normal);
  cursor: pointer;
}

.artist-card:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-neon);
  border-color: var(--accent-primary);
}

.image-wrapper {
  aspect-ratio: 1;
  width: 100%;
  overflow: hidden;
  position: relative;
}

.artist-image, .placeholder-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform var(--transition-normal);
}

.artist-card:hover .artist-image {
  transform: scale(1.05);
}

.placeholder-image {
  background: linear-gradient(135deg, var(--bg-tertiary), var(--bg-secondary));
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 3rem;
  font-weight: 700;
  color: var(--text-muted);
}

.content {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  flex: 1;
}

.name {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 0.25rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.genre {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin-bottom: 1.5rem;
  flex: 1;
}

.actions {
  display: flex;
  gap: 0.5rem;
  margin-top: auto;
}

.btn-sm {
  padding: 0.25rem 0.75rem;
  font-size: 0.875rem;
  flex: 1;
}

.btn-danger {
  background-color: transparent;
  color: var(--danger);
  border: 1px solid var(--danger);
  border-radius: var(--radius-md);
}

.btn-danger:hover {
  background-color: var(--danger);
  color: white;
}
</style>
