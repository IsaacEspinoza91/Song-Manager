<script setup>
</script>

<template>
  <main class="home-view">
    <div class="hero glass-panel">
      <h1 class="gradient-text">Bienvenido a Song Manager</h1>
      <p>Tu plataforma premium para gestionar artistas, canciones y álbumes.</p>
    </div>
  </main>
</template>

<style scoped>
.home-view {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 80vh;
}

.hero {
  text-align: center;
  padding: 4rem;
  max-width: 800px;
}

h1 {
  font-size: 3rem;
  margin-bottom: 1rem;
}

p {
  font-size: 1.25rem;
  color: var(--text-secondary);
}
</style>
