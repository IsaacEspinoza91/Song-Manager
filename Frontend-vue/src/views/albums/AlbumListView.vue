<script setup>
import { ref, onMounted, reactive, watch } from 'vue';
import { albumService } from '../../services/album.service';
import { artistService } from '../../services/artist.service';
import { songService } from '../../services/song.service';
import AlbumCard from '../../components/albums/AlbumCard.vue';
import Modal from '../../components/common/Modal.vue';

const albums = ref([]);
const loading = ref(true);
const error = ref(null);

const pagination = reactive({
  page: 1,
  limit: 10,
  total_pages: 1,
  total_items: 0
});

const filters = reactive({
  title: '',
  artist_name: '',
  artist_id: '',
  type: ''
});

const searchArtists = ref([]); // For the combo box filter

const fetchAlbums = async () => {
  loading.value = true;
  try {
    const params = {
      page: pagination.page,
      limit: pagination.limit,
      title: filters.title,
      artist_name: filters.artist_name,
      artist_id: filters.artist_id,
      type: filters.type
    };
    const response = await albumService.getPaginated(params);
    albums.value = response.data || [];
    pagination.total_pages = response.total_pages || 1;
    pagination.total_items = response.total_items || 0;
  } catch (err) {
    error.value = 'No se pudieron cargar los álbumes';
    console.error(err);
  } finally {
    loading.value = false;
  }
};

const handleSearch = () => {
  pagination.page = 1;
  fetchAlbums();
};

const changePage = (delta) => {
  const newPage = pagination.page + delta;
  if (newPage >= 1 && newPage <= pagination.total_pages) {
    pagination.page = newPage;
    fetchAlbums();
  }
};

const loadArtistsForFilter = async () => {
    try {
        const resp = await artistService.getAll();
        searchArtists.value = resp.data || resp;
        if (searchArtists.value.data) searchArtists.value = searchArtists.value.data;
    } catch (err) {
        console.error('Failed to load artists for filter', err);
    }
};

// Modal State
const isModalOpen = ref(false);
const isEditing = ref(false);
const formError = ref(null);
const availableArtists = ref([]);

const albumForm = reactive({
  id: null,
  title: '',
  release_date: '',
  type: 'LP',
  cover_url: '',
  artist_id: '',
  is_primary: true
});

const fetchAvailableArtists = async () => {
  try {
    const resp = await artistService.getAll();
    availableArtists.value = resp.data || resp;
    if (availableArtists.value.data) availableArtists.value = availableArtists.value.data;
  } catch(err) {
    console.error('Could not load artists', err);
  }
};

const openCreateModal = async () => {
  isEditing.value = false;
  Object.assign(albumForm, { id: null, title: '', release_date: '', type: 'LP', cover_url: '', artist_id: '', is_primary: true });
  formError.value = null;
  if(availableArtists.value.length === 0) await fetchAvailableArtists();
  isModalOpen.value = true;
};

const handleEdit = async (album) => {
  isEditing.value = true;
  Object.assign(albumForm, {
    id: album.id,
    title: album.title,
    release_date: album.release_date ? album.release_date.split('T')[0] : '', // format for date input
    type: album.type || 'LP',
    cover_url: album.cover_url || '',
    artist_id: album.artists && album.artists.length > 0 ? album.artists[0].id : '',
    is_primary: album.artists && album.artists.length > 0 ? album.artists[0].is_primary : true
  });
  formError.value = null;
  if(availableArtists.value.length === 0) await fetchAvailableArtists();
  isModalOpen.value = true;
};

const saveAlbum = async () => {
  try {
    formError.value = null;
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
        formError.value = 'Un álbum requiere de un artista.';
        return;
    }

    if (isEditing.value) {
      await albumService.update(albumForm.id, payload);
    } else {
      await albumService.create(payload);
    }
    isModalOpen.value = false;
    await fetchAlbums();
  } catch (err) {
    console.error(err);
    formError.value = 'Error al guardar el álbum.';
  }
};

const handleDelete = async (id) => {
  if(confirm('¿Estás seguro de que quieres eliminar este álbum?')) {
      try {
        await albumService.delete(id);
        albums.value = albums.value.filter(a => a.id !== id);
      } catch(err) {
        console.error(err);
        alert('Error eliminando el álbum');
      }
    }
};

// Track Management State
const isManageTracksOpen = ref(false);
const activeAlbum = ref(null);
const albumTracks = ref([]);
const availableSongs = ref([]);
const trackForm = reactive({
  song_id: '',
  track_number: 1
});

const openManageTracks = async (album) => {
  activeAlbum.value = album;
  formError.value = null;
  isManageTracksOpen.value = true;
  await loadAlbumTracks(album.id);
  await loadAvailableSongs();
};

const loadAlbumTracks = async (albumId) => {
    try {
        // Detailed album fetch includes tracks
        const resp = await albumService.getById(albumId);
        albumTracks.value = resp.tracks || [];
        trackForm.track_number = albumTracks.value.length + 1;
    } catch(err) {
        console.error('Failed to load tracks', err);
        formError.value = 'No se pudieron cargar las pistas del álbum.';
    }
};

const loadAvailableSongs = async () => {
    try {
        const resp = await songService.getAll();
        availableSongs.value = resp.data || resp;
    } catch(err) {
        console.error('Failed to load songs', err);
    }
};

const addTrack = async () => {
  try {
      formError.value = null;
      if (!trackForm.song_id) {
          formError.value = 'Selecciona una canción.';
          return;
      }
      
      await albumService.addTrack(activeAlbum.value.id, {
          song_id: Number(trackForm.song_id),
          track_number: Number(trackForm.track_number)
      });
      
      trackForm.song_id = '';
      await loadAlbumTracks(activeAlbum.value.id);
  } catch(err) {
      console.error(err);
      formError.value = 'Error al agregar la pista.';
  }
};

const removeTrack = async (songId) => {
    if(confirm('¿Eliminar esta pista del álbum?')) {
        try {
            await albumService.removeTrack(activeAlbum.value.id, songId);
            await loadAlbumTracks(activeAlbum.value.id);
        } catch(err) {
            console.error(err);
            formError.value = 'Error al eliminar la pista.';
        }
    }
};

onMounted(() => {
  fetchAlbums();
  loadArtistsForFilter();
});
</script>

<template>
  <div class="albums-view">
    <div class="header">
      <h1 class="gradient-text">Álbumes</h1>
      <button class="btn btn-primary" @click="openCreateModal">Nuevo Álbum</button>
    </div>

    <!-- Filtros de Búsqueda -->
    <div class="filters glass-panel">
      <form @submit.prevent="handleSearch" class="filter-form">
        <input type="text" v-model="filters.title" placeholder="Buscar por título..." class="form-input" />
        <input type="text" v-model="filters.artist_name" placeholder="Filtro Artista (Nombre)" class="form-input" />
        <select v-model="filters.artist_id" class="form-input">
            <option value="">Todos los artistas (ID)</option>
            <option v-for="a in searchArtists" :key="a.id" :value="a.id">{{ a.name }}</option>
        </select>
        <select v-model="filters.type" class="form-input">
            <option value="">Todos los tipos</option>
            <option value="LP">LP</option>
            <option value="EP">EP</option>
            <option value="Single">Single</option>
        </select>
        <button type="submit" class="btn btn-primary shrink-btn">Buscar</button>
      </form>
    </div>
    
    <div v-if="loading" class="loading">Cargando álbumes...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else class="grid">
      <div v-for="album in albums" :key="album.id" class="album-card-wrapper">
          <AlbumCard 
            :album="album"
            @edit="handleEdit"
            @delete="handleDelete"
          />
          <button class="btn btn-secondary w-full mt-2" @click="openManageTracks(album)">
              Gestionar Pistas
          </button>
      </div>
      <div v-if="albums.length === 0" class="empty-state">
        No se encontraron álbumes.
      </div>
    </div>

    <!-- Controles de Paginación -->
    <div v-if="!loading && !error && albums.length > 0" class="pagination">
      <button class="btn btn-secondary" :disabled="pagination.page === 1" @click="changePage(-1)">
        Anterior
      </button>
      <span class="page-info">
        Página {{ pagination.page }} de {{ pagination.total_pages }} ({{ pagination.total_items }} elementos)
      </span>
      <button class="btn btn-secondary" :disabled="pagination.page === pagination.total_pages" @click="changePage(1)">
        Siguiente
      </button>
    </div>

    <!-- Modal Form -->
    <Modal :isOpen="isModalOpen" @close="isModalOpen = false" :title="isEditing ? 'Editar Álbum' : 'Nuevo Álbum'">
      <form @submit.prevent="saveAlbum">
        <div v-if="formError" class="error-msg">{{ formError }}</div>
        
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
          <select v-model="albumForm.artist_id" class="form-input" required>
            <option value="">Selecciona un Artista</option>
            <option v-for="artist in availableArtists" :key="artist.id" :value="artist.id">
              {{ artist.name }}
            </option>
          </select>
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn-secondary" @click="isModalOpen = false">Cancelar</button>
          <button type="submit" class="btn btn-primary">Guardar</button>
        </div>
      </form>
    </Modal>

    <!-- Track Management Modal -->
    <Modal :isOpen="isManageTracksOpen" @close="isManageTracksOpen = false" :title="`Pistas: ${activeAlbum?.title}`">
        <div v-if="formError" class="error-msg">{{ formError }}</div>
        
        <div class="tracks-list mb-4">
            <h4 class="mb-2 text-sm text-secondary">Pistas Actuales</h4>
            <div v-if="albumTracks.length === 0" class="text-sm text-muted mb-2">
                No hay pistas en este álbum.
            </div>
            <ul class="track-items">
                <li v-for="track in albumTracks" :key="track.song_id" class="track-item flex justify-between items-center glass-panel p-2 mb-2 rounded">
                    <span>
                        <span class="text-muted mr-2">{{ track.track_number }}.</span>
                        {{ track.title }}
                    </span>
                    <button class="icon-btn text-danger" @click="removeTrack(track.song_id)">🗑️</button>
                </li>
            </ul>
        </div>

        <form @submit.prevent="addTrack" class="add-track-form border-t pt-4">
            <h4 class="mb-2 text-sm text-secondary">Agregar Pista</h4>
            <div class="form-group">
                <label>Canción</label>
                <select v-model="trackForm.song_id" class="form-input" required>
                    <option value="">Seleccionar canción...</option>
                    <option v-for="song in availableSongs" :key="song.id" :value="song.id">
                        {{ song.title }}
                    </option>
                </select>
            </div>
            <div class="form-group flex gap-2">
                 <div class="w-1/3">
                    <label>Número</label>
                    <input type="number" v-model="trackForm.track_number" class="form-input" required min="1" />
                 </div>
                 <div class="w-2/3 flex items-end">
                    <button type="submit" class="btn btn-primary w-full h-11">Agregar Pista</button>
                 </div>
            </div>
        </form>
    </Modal>
  </div>
</template>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

h1 {
  font-size: 2.5rem;
  font-weight: 700;
}

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 2rem;
}

.loading, .error, .empty-state {
  text-align: center;
  padding: 3rem;
  color: var(--text-secondary);
  font-size: 1.25rem;
  grid-column: 1 / -1;
}

.error {
  color: var(--danger);
}

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

.filters {
  padding: 1rem 1.5rem;
  margin-bottom: 2rem;
  display: flex;
}

.filter-form {
  display: flex;
  gap: 1rem;
  width: 100%;
}

.filter-form .form-input {
  flex: 1;
}

.shrink-btn {
  flex-shrink: 0;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1.5rem;
  margin-top: 3rem;
  padding-bottom: 2rem;
}

.page-info {
  color: var(--text-secondary);
  font-size: 0.9rem;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.album-card-wrapper {
    display: flex;
    flex-direction: column;
}

.w-full { width: 100%; }
.mt-2 { margin-top: 0.5rem; }
.mb-2 { margin-bottom: 0.5rem; }
.mb-4 { margin-bottom: 1rem; }
.pt-4 { padding-top: 1rem; }
.p-2 { padding: 0.5rem; }
.border-t { border-top: 1px solid var(--border-color); }
.text-sm { font-size: 0.875rem; }
.text-secondary { color: var(--text-secondary); }
.text-muted { color: var(--text-muted); }
.flex { display: flex; }
.justify-between { justify-content: space-between; }
.items-center { align-items: center; }
.items-end { align-items: flex-end; }
.gap-2 { gap: 0.5rem; }
.w-1\/3 { width: 33.333333%; }
.w-2\/3 { width: 66.666667%; }
.rounded { border-radius: var(--radius-sm); }
.mr-2 { margin-right: 0.5rem; }

.track-item {
    transition: background-color var(--transition-fast);
}

.track-item:hover {
    background: rgba(255,255,255,0.05);
}

.icon-btn {
    background: none;
    border: none;
    cursor: pointer;
    font-size: 1.1rem;
    transition: transform 0.2s;
}

.icon-btn:hover {
    transform: scale(1.2);
}
</style>
