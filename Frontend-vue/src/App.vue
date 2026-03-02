<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <div class="app-layout">
    <header class="sidebar glass-panel">
      <div class="logo">
        <h2 class="gradient-text">SongManager</h2>
      </div>
      <nav class="nav-menu">
        <RouterLink to="/" class="nav-item">Inicio</RouterLink>
        <RouterLink to="/artists" class="nav-item">Artistas</RouterLink>
        <RouterLink to="/songs" class="nav-item">Canciones</RouterLink>
        <RouterLink to="/albums" class="nav-item">Álbumes</RouterLink>
      </nav>
    </header>

    <main class="content-area">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
}

.sidebar {
  width: 260px;
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  padding: 2rem 1rem;
  display: flex;
  flex-direction: column;
  z-index: 10;
  border-radius: 0; /* Override glass-panel radius for sidebar attached to left */
  border-left: none;
  border-top: none;
  border-bottom: none;
}

.logo {
  margin-bottom: 3rem;
  padding-left: 1rem;
}

.logo h2 {
  font-size: 1.5rem;
  font-weight: 700;
  letter-spacing: -0.5px;
}

.nav-menu {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.nav-item {
  padding: 0.75rem 1rem;
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-weight: 500;
  transition: all var(--transition-fast);
}

.nav-item:hover, .nav-item.router-link-active {
  background: rgba(255, 255, 255, 0.05);
  color: var(--text-primary);
}

.nav-item.router-link-active {
  background: rgba(30, 215, 96, 0.15); /* Slightly green tint */
  color: var(--accent-primary);
  border-left: 3px solid var(--accent-primary);
  border-top-left-radius: 0;
  border-bottom-left-radius: 0;
}

.content-area {
  flex: 1;
  margin-left: 260px;
  padding: 2rem;
  min-height: 100vh;
  position: relative;
  overflow-y: auto;
}
</style>
