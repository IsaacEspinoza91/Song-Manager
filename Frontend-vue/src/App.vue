<script setup>
import { RouterLink, RouterView } from 'vue-router'
import ToastContainer from './components/common/ToastContainer.vue';
import Icon from './components/common/Icon.vue';
</script>

<template>
  <div class="app-layout">
    <header class="sidebar glass-panel">
      <div class="logo">
        <h2 class="gradient-text">SongManager</h2>
      </div>
      <nav class="nav-menu">
        <RouterLink to="/" class="nav-item" title="Ir al Inicio">
          <Icon name="home" class="nav-icon" /> Inicio
        </RouterLink>
        <RouterLink to="/artists" class="nav-item" title="Explorar Artistas">
          <Icon name="users" class="nav-icon" /> Artistas
        </RouterLink>
        <RouterLink to="/songs" class="nav-item" title="Explorar Canciones">
          <Icon name="music" class="nav-icon" /> Canciones
        </RouterLink>
        <RouterLink to="/albums" class="nav-item" title="Explorar Álbumes">
          <Icon name="disc" class="nav-icon" /> Álbumes
        </RouterLink>
      </nav>

      <div class="sidebar-spacer"></div>

      <div class="user-footer">
        <div class="user-info">
          <div class="user-avatar">
            <span>IE</span>
          </div>
          <div class="user-details">
            <span class="user-name"><small>Creado por</small></span>
            <span class="user-name">Isaac Espinoza</span>
            <div class="user-links">
              <!-- GitHub -->
              <a href="https://github.com/IsaacEspinoza91/Song-Manager" target="_blank" rel="noopener noreferrer" title="GitHub">
                <svg viewBox="0 0 24 24" class="social-icon">
                  <path fill="currentColor" d="M12 2C6.477 2 2 6.477 2 12c0 4.42 2.865 8.166 6.839 9.489.5.092.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.603-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.464-1.11-1.464-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.831.092-.646.35-1.086.636-1.336-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.025A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.294 2.747-1.025 2.747-1.025.546 1.377.203 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.578.688.48C19.138 20.161 22 16.418 22 12c0-5.523-4.477-10-10-10z"></path>
                </svg>
              </a>
              <!-- LinkedIn -->
              <!--<a href="https://linkedin.com/in/tu-perfil" target="_blank" rel="noopener noreferrer" title="LinkedIn">
                <svg viewBox="0 0 24 24" class="social-icon">
                  <path fill="currentColor" d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/>
                </svg>
              </a> -->
              <!-- Email -->
              <a href="mailto:isaac.espinoza.barria@gmail.com" title="Email">
                <svg viewBox="0 0 24 24" class="social-icon">
                  <path fill="currentColor" d="M20 4H4C2.9 4 2.01 4.9 2.01 6L2 18C2 19.1 2.9 20 4 20H20C21.1 20 22 19.1 22 18V6C22 4.9 21.1 4 20 4ZM20 8L12 13L4 8V6L12 11L20 6V8Z"/>
                </svg>
              </a>
            </div>
          </div>
        </div>
      </div>
    </header>

    <main class="content-area">
      <RouterView />
    </main>

    <ToastContainer />
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
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-weight: 500;
  transition: all var(--transition-fast);
}

.nav-icon {
  opacity: 0.7;
  transition: opacity var(--transition-fast);
}

.nav-item:hover, .nav-item.router-link-active {
  background: rgba(255, 255, 255, 0.05);
  color: var(--text-primary);
}

.nav-item:hover .nav-icon, .nav-item.router-link-active .nav-icon {
  opacity: 1;
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

/* User Footer Styles */
.sidebar-spacer {
  flex: 1;
}

.user-footer {
  margin-top: 2rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--border-color);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--accent-primary) 0%, #1e3c72 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 1.1rem;
  color: white;
  flex-shrink: 0;
  box-shadow: 0 4px 10px rgba(30, 215, 96, 0.3);
}

.user-details {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  overflow: hidden;
}

.user-name {
  font-size: 0.85rem;
  font-weight: 500;
  color: var(--text-primary);
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
}

.user-links {
  display: flex;
  gap: 0.6rem;
  align-items: center;
}

.social-icon {
  width: 18px;
  height: 18px;
  color: var(--text-muted);
  transition: all var(--transition-fast);
}

.social-icon:hover {
  color: var(--accent-primary);
  transform: translateY(-2px);
}
</style>
