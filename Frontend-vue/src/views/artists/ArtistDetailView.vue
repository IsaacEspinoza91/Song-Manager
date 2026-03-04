<script setup>
import { ref, onMounted, onUnmounted, computed, reactive, watch, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { artistService } from '../../services/artist.service';
import { songService } from '../../services/song.service';
import { albumService } from '../../services/album.service';
import Breadcrumbs from '../../components/common/Breadcrumbs.vue';
import SongItem from '../../components/songs/SongItem.vue';
import AlbumCard from '../../components/albums/AlbumCard.vue';
import Modal from '../../components/common/Modal.vue';
import ConfirmDeleteModal from '../../components/common/ConfirmDeleteModal.vue';
import SearchSelect from '../../components/common/SearchSelect.vue';

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
// SONG CRUD (EDIT & DELETE)
// ========================
const isSongModalOpen = ref(false);
const songFormError = ref(null);
const songForm = reactive({
  id: null,
  title: '',
  duration: 0,
  artist_id: '',
  artist_name: '',
  role: 'main'
});

const handleEditSong = (song) => {
  Object.assign(songForm, {
    id: song.id,
    title: song.title,
    duration: song.duration,
    artist_id: song.artists && song.artists.length > 0 ? song.artists[0].id : '',
    artist_name: song.artists && song.artists.length > 0 ? song.artists[0].name || song.artists[0].artist_name : '',
    role: song.artists && song.artists.length > 0 ? song.artists[0].role : 'main'
  });
  songFormError.value = null;
  isSongModalOpen.value = true;
};

const saveSong = async () => {
  try {
    songFormError.value = null;
    let payload = {
      title: songForm.title,
      duration: Number(songForm.duration)
    };
    if (songForm.artist_id) {
       payload.artists = [{
           artist_id: Number(songForm.artist_id),
           role: songForm.role
       }];
    }
    await songService.update(songForm.id, payload);
    isSongModalOpen.value = false;
    
    // Refresh songs for this artist
    const songsResp = await songService.getPaginated({ artist_id: artistId.value, limit: 5 });
    songs.value = songsResp.data || [];
  } catch (err) {
    console.error(err);
    songFormError.value = 'Error al editar la canción.';
  }
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
const albumForm = reactive({
  id: null,
  title: '',
  release_date: '',
  type: 'LP',
  cover_url: '',
  artist_id: '',
  artist_name: '',
  is_primary: true,
  artists: []
});

const isEditingAlbum = ref(false);
const isAlbumModalOpen = ref(false);
const albumFormError = ref(null);

const openCreateAlbum = () => {
    isEditingAlbum.value = false;
    Object.assign(albumForm, {
        id: null,
        title: '',
        release_date: '',
        type: 'LP',
        cover_url: '',
        artist_id: artistId.value,
        artist_name: artist.value ? artist.value.name : '',
        is_primary: true,
        artists: [{ artist_id: artistId.value, is_primary: true }]
    });
    albumFormError.value = null;
    isAlbumModalOpen.value = true;
};

const handleEditAlbum = (album) => {
  isEditingAlbum.value = true;
  Object.assign(albumForm, {
    id: album.id,
    title: album.title,
    release_date: album.release_date ? album.release_date.split('T')[0] : '', // format for date input
    type: album.type || 'LP',
    cover_url: album.cover_url || '',
    artist_id: album.artists && album.artists.length > 0 ? album.artists[0].id : '',
    artist_name: album.artists && album.artists.length > 0 ? album.artists[0].name || album.artists[0].artist_name : '',
    is_primary: album.artists && album.artists.length > 0 ? album.artists[0].is_primary : true
  });
  albumFormError.value = null;
  isAlbumModalOpen.value = true;
};

const saveAlbum = async () => {
  try {
    albumFormError.value = null;
    let payload = {
      title: albumForm.title,
      release_date: albumForm.release_date || undefined,
      type: albumForm.type,
      cover_url: albumForm.cover_url || undefined,
      artists: []
    };

    if (albumForm.artist_id) {
       payload.artists = [{
           artist_id: Number(albumForm.artist_id),
           is_primary: albumForm.is_primary
       }];
    } else {
        albumFormError.value = 'Un álbum requiere de un artista.';
        return;
    }

    if (isEditingAlbum.value) {
        await albumService.update(albumForm.id, payload);
    } else {
        await albumService.create(payload);
    }
    
    isAlbumModalOpen.value = false;
    
    // Refresh albums for this artist
    await fetchAlbums();
  } catch (err) {
    console.error(err);
    albumFormError.value = 'Error al editar el álbum.';
  }
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
                  <p class="artist-dates">
                      Añadido: {{ formatDate(artist.created_at) }}<br/>
                      Última actualización: {{ formatDate(artist.updated_at) }}
                  </p>
              </div>
          </div>
      </div>

      <div class="content-grid mt-4">
          <!-- Top Songs -->
          <div class="songs-section glass-panel p-4">
              <h2 class="section-title mb-4">Canciones</h2>
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

    <!-- Song Edit Modal -->
    <Modal :isOpen="isSongModalOpen" @close="isSongModalOpen = false" title="Editar Canción">
      <form @submit.prevent="saveSong">
        <div v-if="songFormError" class="error-msg">{{ songFormError }}</div>
        
        <div class="form-group">
          <label>Título</label>
          <input type="text" v-model="songForm.title" class="form-input" required />
        </div>
        
        <div class="form-group">
          <label>Duración (Segundos)</label>
          <input type="number" v-model="songForm.duration" class="form-input" required min="1" />
        </div>

        <div class="form-group">
          <label>Artista Principal</label>
          <SearchSelect 
              v-model="songForm.artist_id"
              :initialName="songForm.artist_name"
              :searchFn="artistService.search"
              :formatDisplay="(a) => a.artist_name || a.name"
              placeholder="Busca un artista..."
              @select="(item) => songForm.artist_name = (item.artist_name || item.name)"
          />
        </div>

        <div class="form-group" v-if="songForm.artist_id">
            <label>Rol del Artista</label>
            <select v-model="songForm.role" class="form-input">
                <option value="main">Main (Principal)</option>
                <option value="ft">Featuring (Invitado)</option>
                <option value="producer">Productor</option>
            </select>
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn-secondary" @click="isSongModalOpen = false">Cancelar</button>
          <button type="submit" class="btn btn-primary">Guardar</button>
        </div>
      </form>
    </Modal>

    <!-- Album Edit Modal -->
    <Modal :isOpen="isAlbumModalOpen" @close="isAlbumModalOpen = false" :title="isEditingAlbum ? 'Editar Álbum' : 'Crear Nuevo Álbum'">
      <form @submit.prevent="saveAlbum">
        <div v-if="albumFormError" class="error-msg">{{ albumFormError }}</div>
        
        <div class="form-group">
          <label>Título</label>
          <input type="text" v-model="albumForm.title" class="form-input" required />
        </div>

        <div class="form-group">
          <label>Tipo</label>
          <select v-model="albumForm.type" class="form-input" required>
            <option value="LP">LP (Long Play)</option>
            <option value="EP">EP (Extended Play)</option>
            <option value="Single">Single (Sencillo)</option>
          </select>
        </div>
        
        <div class="form-group">
          <label>Fecha de Lanzamiento</label>
          <input type="date" v-model="albumForm.release_date" class="form-input" required />
        </div>

        <div class="form-group">
          <label>URL de Portada (Opcional)</label>
          <input type="url" v-model="albumForm.cover_url" class="form-input" />
        </div>

        <div class="form-group">
          <label>Artista Principal</label>
          <SearchSelect 
              v-model="albumForm.artist_id"
              :initialName="albumForm.artist_name"
              :searchFn="artistService.search"
              :formatDisplay="(a) => a.artist_name || a.name"
              placeholder="Busca un artista..."
              @select="(item) => albumForm.artist_name = (item.artist_name || item.name)"
          />
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn-secondary" @click="isAlbumModalOpen = false">Cancelar</button>
          <button type="submit" class="btn btn-primary">Guardar</button>
        </div>
      </form>
    </Modal>

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
.p-4 { padding: 1.5rem; }

.loading, .error, .empty-state {
    text-align: center;
    color: var(--text-secondary);
    padding: 2rem;
}

.error { color: var(--danger); }
.error-msg {
  color: var(--danger);
  margin-bottom: 1rem;
  font-size: 0.9rem;
  background: rgba(239, 68, 68, 0.1);
  padding: 0.5rem;
  border-radius: var(--radius-sm);
}

select.form-input option {
  background: var(--bg-primary);
  color: var(--text-primary);
}

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
