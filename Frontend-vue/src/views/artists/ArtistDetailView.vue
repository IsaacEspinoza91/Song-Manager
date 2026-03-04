<script setup>
import { ref, onMounted, onUnmounted, computed, reactive, watch, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { artistService } from '../../services/artist.service';
import { songService } from '../../services/song.service';
import { albumService } from '../../services/album.service';
import Breadcrumbs from '../../components/common/Breadcrumbs.vue';
import SongItem from '../../components/songs/SongItem.vue';
import AlbumCard from '../../components/albums/AlbumCard.vue';
import SongFormModal from '../../components/songs/SongFormModal.vue';
import AlbumFormModal from '../../components/albums/AlbumFormModal.vue';
import ArtistFormModal from '../../components/artists/ArtistFormModal.vue';

const route = useRoute();
const router = useRouter();
const artistId = computed(() => route.params.id);

watch(artistId, (newId) => {
  if (newId) {
    fetchData();
  }
});

const loading = ref(true);
const error = ref(null);
const artist = ref(null);
const songs = ref([]);
const albums = ref([]);
const activeAlbumTab = ref(''); // '' = Todos, 'LP' = Álbumes, 'EP' = EPs, 'Single' = Singles

const albumsContainer = ref(null);
const albumLimit = ref(5);

const calculateAlbumLimit = () => {
    if (!albumsContainer.value) return 5;
    // albumsContainer has p-4 (1.5rem padding = 24px each side -> 48px total inner reduction)
    const availableWidth = albumsContainer.value.clientWidth - 48;
    let cols = Math.floor((availableWidth + 24) / 224); // Based on minmax(200px) + gap(1.5rem=24px)
    if (cols < 1) cols = 1;
    return cols; // 1 row only
};

let resizeTimer;
const handleResize = () => {
    clearTimeout(resizeTimer);
    resizeTimer = setTimeout(() => {
        const newLimit = calculateAlbumLimit();
        if (newLimit !== albumLimit.value && newLimit > 0) {
            albumLimit.value = newLimit;
            fetchAlbums();
        }
    }, 300);
};

const fetchAlbums = async () => {
    try {
        const params = {
            artist_id: artistId.value,
            limit: albumLimit.value,
            page: 1,
            type: activeAlbumTab.value
        };
        const albumsResp = await albumService.getPaginated(params);
        albums.value = albumsResp.data || [];
    } catch(err) {
        console.error('Failed to load albums', err);
    }
};

const setAlbumTab = (tab) => {
    activeAlbumTab.value = tab;
    fetchAlbums();
};

const breadcrumbItems = computed(() => {
  return [
    { label: 'Inicio', to: '/' },
    { label: 'Artistas', to: '/artists' },
    { label: artist.value ? artist.value.name : 'Cargando...' }
  ];
});

const formatDate = (dateString) => {
    if (!dateString) return 'Desconocido';
    return new Date(dateString).toLocaleDateString();
};

const fetchData = async () => {
  loading.value = true;
  try {
    const artistResp = await artistService.getById(artistId.value);
    artist.value = artistResp.data || artistResp;

    const songsResp = await songService.getPaginated({ artist_id: artistId.value, limit: 5 });
    songs.value = songsResp.data || [];

  } catch(err) {
    console.error(err);
    error.value = 'No se pudieron cargar los datos del artista.';
  } finally {
    loading.value = false;
    await nextTick();
    albumLimit.value = calculateAlbumLimit();
    await fetchAlbums();
  }
};

const goToAllSongs = () => {
    router.push(`/artists/${artistId.value}/songs`);
};

const goToDiscography = () => {
    router.push(`/artists/${artistId.value}/discography`);
};

// ========================
// ARTIST SUMMARY (EDIT)
// ========================
const isArtistModalOpen = ref(false);

const openEditArtist = () => {
    isArtistModalOpen.value = true;
};

const handleArtistSaved = (updatedArtist) => {
    artist.value = updatedArtist;
    // update url if id changed? (id won't change)
};

// ========================
// SONG CRUD (EDIT & DELETE)
// ========================
const isSongModalOpen = ref(false);
const selectedSong = ref(null);

const handleEditSong = (song) => {
  selectedSong.value = song;
  isSongModalOpen.value = true;
};

const openCreateSong = () => {
  selectedSong.value = null;
  isSongModalOpen.value = true;
};

const handleSongSaved = async () => {
    const songsResp = await songService.getPaginated({ artist_id: artistId.value, limit: 5 });
    songs.value = songsResp.data || [];
};

const isDeleteModalOpen = ref(false);
const isAlbumDelete = ref(false);
const itemToDeleteId = ref(null);
const itemToDeleteName = ref('');

const handleDeleteSong = (id, title) => {
  isAlbumDelete.value = false;
  itemToDeleteId.value = id;
  itemToDeleteName.value = title || 'esta canción';
  isDeleteModalOpen.value = true;
};

const executeDelete = async () => {
    try {
        if (isAlbumDelete.value) {
            await albumService.delete(itemToDeleteId.value);
            await fetchAlbums();
        } else {
            await songService.delete(itemToDeleteId.value);
            const songsResp = await songService.getPaginated({ artist_id: artistId.value, limit: 5 });
            songs.value = songsResp.data || [];
        }
        isDeleteModalOpen.value = false;
    } catch(err) {
        console.error(err);
        alert('Error en la eliminación');
    }
};

// ========================
// ALBUM CRUD (EDIT & DELETE)
// ========================
const isAlbumModalOpen = ref(false);
const selectedAlbum = ref(null);

const openCreateAlbum = () => {
    selectedAlbum.value = null;
    isAlbumModalOpen.value = true;
};

const handleEditAlbum = (album) => {
    selectedAlbum.value = album;
    isAlbumModalOpen.value = true;
};

const handleAlbumSaved = async () => {
    await fetchAlbums();
};

const handleDeleteAlbum = (id, title) => {
  isAlbumDelete.value = true;
  itemToDeleteId.value = id;
  itemToDeleteName.value = title || 'este álbum';
  isDeleteModalOpen.value = true;
};

onMounted(() => {
  albumLimit.value = calculateAlbumLimit();
  fetchData();
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});
</script>

<template>
  <div class="artist-detail-view">
    <Breadcrumbs :items="breadcrumbItems" />
    
    <div v-if="loading" class="loading">Cargando perfil...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="artist">
      
      <!-- Hero Section -->
      <div class="hero glass-panel">
          <div class="hero-bg" :style="{ backgroundImage: `url(${artist.image_url || ''})` }"></div>
          <div class="hero-content">
              <div class="profile-img-wrapper">
                  <img v-if="artist.image_url" :src="artist.image_url" :alt="artist.name" class="profile-img" />
                  <div v-else class="placeholder-profile">{{ artist.name.charAt(0) }}</div>
              </div>
              <div class="hero-info">
                  <h1 class="artist-name">{{ artist.name }}</h1>
                  <p class="artist-meta">{{ artist.genre }} • {{ artist.country }}</p>
                  <p class="artist-bio" v-if="artist.bio">{{ artist.bio }}</p>
                  <div class="flex justify-between items-end w-full">
                      <p class="artist-dates mb-0">
                          Añadido: {{ formatDate(artist.created_at) }}<br/>
                          Última actualización: {{ formatDate(artist.updated_at) }}
                      </p>
                      <button class="btn btn-secondary btn-sm flex items-center gap-2" @click="openEditArtist">
                          <span>✏️</span> Editar Artista
                      </button>
                  </div>
              </div>
          </div>
      </div>

      <div class="content-grid mt-4">
          <!-- Top Songs -->
          <div class="songs-section glass-panel p-4">
              <div class="flex justify-between items-center mb-4">
                  <h2 class="section-title mb-0">Canciones</h2>
                  <button class="btn btn-primary btn-sm" @click="openCreateSong">+ Nueva Canción</button>
              </div>
              <div v-if="songs.length === 0" class="empty-state mb-4">No hay canciones registradas.</div>
              <div v-else class="songs-list mb-4">
                  <SongItem 
                    v-for="(song, index) in songs" 
                    :key="song.id" 
                    :song="song" 
                    :index="index" 
                    :readonly="false"
                    @edit="handleEditSong"
                    @delete="handleDeleteSong(song.id, song.title)"
                  />
              </div>
              <button class="btn btn-secondary w-full" @click="goToAllSongs">Ver todas las canciones</button>
          </div>

          <!-- Albums -->
          <div class="albums-section glass-panel p-4 mt-4" ref="albumsContainer">
              <div class="flex justify-between items-center mb-4">
                  <h2 class="section-title mb-0">Discografía</h2>
                  <button class="btn btn-primary btn-sm" @click="openCreateAlbum">+ Nuevo Álbum</button>
              </div>
              
              <div class="tabs mb-4">
                  <button class="tab-btn" :class="{ active: activeAlbumTab === '' }" @click="setAlbumTab('')">Todos</button>
                  <button class="tab-btn" :class="{ active: activeAlbumTab === 'LP' }" @click="setAlbumTab('LP')">Álbumes</button>
                  <button class="tab-btn" :class="{ active: activeAlbumTab === 'EP' }" @click="setAlbumTab('EP')">EPs</button>
                  <button class="tab-btn" :class="{ active: activeAlbumTab === 'Single' }" @click="setAlbumTab('Single')">Singles</button>
              </div>

              <div v-if="albums.length === 0" class="empty-state mb-4">No hay álbumes en esta categoría.</div>
              <div v-else class="albums-grid">
                  <AlbumCard 
                    v-for="album in albums" 
                    :key="album.id" 
                    :album="album" 
                    :readonly="false"
                    @edit="handleEditAlbum"
                    @delete="handleDeleteAlbum(album.id, album.title)"
                  />
              </div>
              
              <button class="btn btn-secondary w-full" style="margin-top: 1rem" @click="goToDiscography">Ver toda la discografía</button>
          </div>
      </div>
    </div>

    <!-- Modals -->
    <ArtistFormModal
        :isOpen="isArtistModalOpen"
        :artist="artist"
        @close="isArtistModalOpen = false"
        @saved="handleArtistSaved"
    />

    <SongFormModal 
        :isOpen="isSongModalOpen" 
        :song="selectedSong" 
        :initialArtist="artist"
        @close="isSongModalOpen = false" 
        @saved="handleSongSaved" 
    />

    <AlbumFormModal 
        :isOpen="isAlbumModalOpen" 
        :album="selectedAlbum" 
        :initialArtist="artist"
        @close="isAlbumModalOpen = false" 
        @saved="handleAlbumSaved" 
    />

    <!-- Confirm Delete Modal -->
    <ConfirmDeleteModal 
      :isOpen="isDeleteModalOpen" 
      :itemName="itemToDeleteName"
      @close="isDeleteModalOpen = false"
      @confirm="executeDelete"
    />
  </div>
</template>

<style scoped>
.artist-detail-view {
  padding-bottom: 2rem;
}

.hero {
  position: relative;
  overflow: hidden;
  border-radius: var(--radius-lg);
  min-height: 300px;
  display: flex;
  align-items: flex-end;
  padding: 2rem;
  margin-bottom: 2rem;
}

.hero-bg {
    position: absolute;
    top: 0; left: 0; right: 0; bottom: 0;
    background-size: cover;
    background-position: center;
    filter: blur(20px) brightness(0.3);
    z-index: 0;
}

.hero-content {
    position: relative;
    z-index: 1;
    display: flex;
    gap: 2rem;
    align-items: center;
    width: 100%;
}

.profile-img-wrapper {
    width: 200px;
    height: 200px;
    border-radius: 50%;
    overflow: hidden;
    border: 4px solid var(--border-color);
    box-shadow: var(--shadow-neon);
    flex-shrink: 0;
}

.profile-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.placeholder-profile {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 5rem;
    font-weight: 700;
    background: linear-gradient(135deg, var(--accent-primary), var(--accent-cyan));
    color: white;
}

.hero-info {
    flex: 1;
}

.artist-name {
    font-size: 3.5rem;
    font-weight: 800;
    margin-bottom: 0.5rem;
    text-shadow: 2px 2px 4px rgba(0,0,0,0.5);
    line-height: 1.1;
}

.artist-meta {
    font-size: 1.25rem;
    color: var(--accent-primary);
    margin-bottom: 1rem;
    font-weight: 500;
}

.artist-bio {
    font-size: 1rem;
    color: var(--text-secondary);
    line-height: 1.5;
    max-width: 800px;
    margin-bottom: 1rem;
}

.artist-dates {
  font-size: 0.85rem;
  color: var(--text-muted);
}

.section-title {
    font-size: 1.5rem;
    font-weight: 600;
    margin-bottom: 1.5rem;
    color: var(--text-primary);
}

.albums-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 1.5rem;
}

.tabs {
    display: flex;
    gap: 0.5rem;
    border-bottom: 1px solid var(--border-color);
    padding-bottom: 0.5rem;
}

.tab-btn {
    background: transparent;
    color: var(--text-secondary);
    border: none;
    padding: 0.5rem 1rem;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    border-radius: var(--radius-sm);
    transition: all var(--transition-fast);
}

.tab-btn:hover {
    color: var(--text-primary);
    background: rgba(255,255,255,0.05);
}

.tab-btn.active {
    background: var(--accent-primary);
    color: var(--bg-primary);
}

.w-full { width: 100%; }
.mt-4 { margin-top: 1rem; }
.mb-4 { margin-bottom: 1rem; }
.mb-0 { margin-bottom: 0 !important; }
.p-4 { padding: 1.5rem; }

/* Utilities fallback */
.flex { display: flex; }
.justify-between { justify-content: space-between; }
.items-center { align-items: center; }
.items-end { align-items: flex-end; }
.gap-2 { gap: 0.5rem; }

.loading, .error, .empty-state {
    text-align: center;
    color: var(--text-secondary);
    padding: 2rem;
}

.error { color: var(--danger); }


/* Responsive Adjustments */
@media (max-width: 768px) {
    .hero-content {
        flex-direction: column;
        align-items: flex-start;
        text-align: left;
    }
    
    .profile-img-wrapper {
        width: 150px;
        height: 150px;
    }
    
    .artist-name {
        font-size: 2.5rem;
    }
}
</style>
