<script setup>
import { ref, onMounted, onUnmounted, reactive, watch } from 'vue';
import { albumService } from '../../services/album.service';
import { artistService } from '../../services/artist.service';
import { songService } from '../../services/song.service';
import AlbumCard from '../../components/albums/AlbumCard.vue';
import Modal from '../../components/common/Modal.vue';
import AlbumFormModal from '../../components/albums/AlbumFormModal.vue';
import SongFormModal from '../../components/songs/SongFormModal.vue';
import ConfirmDeleteModal from '../../components/common/ConfirmDeleteModal.vue';
import SearchSelect from '../../components/common/SearchSelect.vue';
import Pagination from '../../components/common/Pagination.vue';
import { useToast } from '../../composables/useToast';

const toast = useToast();

const albums = ref([]);
const loading = ref(true);
const error = ref(null);

const pagination = reactive({
  page: 1,
  limit: 10,
  total_pages: 1,
  total_items: 0
});

const viewContainer = ref(null);

const calculateLimit = () => {
    if (!viewContainer.value) return 10;
    const width = viewContainer.value.clientWidth;
    // Calculate columns based on minmax(220px) + gap(32px)
    let cols = Math.floor((width + 32) / 252);
    if (cols < 1) cols = 1;
    return cols * 2; // 2 rows
};

let resizeTimer;
const handleResize = () => {
    clearTimeout(resizeTimer);
    resizeTimer = setTimeout(() => {
        const newLimit = calculateLimit();
        if (newLimit !== pagination.limit && newLimit > 0) {
            pagination.limit = newLimit;
            pagination.page = 1; // Reset to page 1
            fetchAlbums();
        }
    }, 300);
};

const filters = reactive({
  title: '',
  artist_name: '',
  artist_id: '',
  type: ''
});

const filterArtistName = ref('');

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

const changePage = (newPage) => {
  if (newPage >= 1 && newPage <= pagination.total_pages) {
    pagination.page = newPage;
    fetchAlbums();
  }
};

// Modal State
const isModalOpen = ref(false);
const selectedAlbum = ref(null);
const formError = ref(null);

const openCreateModal = async () => {
  selectedAlbum.value = null;
  isModalOpen.value = true;
};

const handleEdit = async (album) => {
  selectedAlbum.value = album;
  isModalOpen.value = true;
};

const handleSaved = async () => {
    await fetchAlbums();
};

const isDeleteModalOpen = ref(false);
const isTrackDelete = ref(false);
const itemToDeleteId = ref(null);
const itemToDeleteName = ref('');

const handleDelete = (id, title) => {
  isTrackDelete.value = false;
  itemToDeleteId.value = id;
  itemToDeleteName.value = title || 'este álbum';
  isDeleteModalOpen.value = true;
};



const executeDelete = async () => {
    try {
        if (isTrackDelete.value) {
           await albumService.removeTrack(activeAlbum.value.id, itemToDeleteId.value);
           toast.success('Pista removida del álbum');
           await loadAlbumTracks(activeAlbum.value.id);
        } else {
           await albumService.delete(itemToDeleteId.value);
           albums.value = albums.value.filter(a => a.id !== itemToDeleteId.value);
           toast.success('Álbum eliminado exitosamente');
        }
        isDeleteModalOpen.value = false;
    } catch(err) {
        toast.handleApiError(err, 'Error en la eliminación');
    }
};

// Track Management State
const isManageTracksOpen = ref(false);
const activeAlbum = ref(null);
const albumTracks = ref([]);
const trackForm = reactive({
  song_id: '',
  song_title: '',
  track_number: 1
});

const formatSongDisplay = (song) => {
    const artistNames = song.artists ? song.artists.map(a => a.artist_name || a.name).join(', ') : '';
    return artistNames ? `${song.title} - ${artistNames}` : song.title;
};

const openManageTracks = async (album) => {
  activeAlbum.value = album;
  formError.value = null;
  isManageTracksOpen.value = true;
  await loadAlbumTracks(album.id);
};

// ========================
// TRACK SONG CREATION
// ========================
const isSongModalOpen = ref(false);

const openCreateTrackSong = () => {
    isSongModalOpen.value = true;
};

const handleSongSaved = (newSong) => {
    if (newSong) {
        trackForm.song_id = newSong.id;
        trackForm.song_title = formatSongDisplay(newSong);
    }
};

const loadAlbumTracks = async (albumId) => {
    try {
        // Detailed album fetch includes tracks
        const resp = await albumService.getById(albumId);
        albumTracks.value = resp.tracks || [];
        trackForm.track_number = albumTracks.value.length + 1;
    } catch(err) {
        toast.handleApiError(err, 'No se pudieron cargar las pistas del álbum');
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
      toast.success('Pista agregada al álbum');
      await loadAlbumTracks(activeAlbum.value.id);
  } catch(err) {
      toast.handleApiError(err, 'Error al agregar la pista');
  }
};

const removeTrack = (songId) => {
    const track = albumTracks.value.find(t => t.song_id === songId);
    isTrackDelete.value = true;
    itemToDeleteId.value = songId;
    itemToDeleteName.value = track ? `la pista "${track.title}" del álbum (la canción original no será eliminada)` : 'esta pista del álbum';
    isDeleteModalOpen.value = true;
};

onMounted(() => {
  pagination.limit = calculateLimit();
  fetchAlbums();
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});
</script>

<template>
  <div class="albums-view" ref="viewContainer">
    <div class="header">
      <h1 class="gradient-text">Álbumes</h1>
      <button class="btn btn-primary" @click="openCreateModal">Nuevo Álbum</button>
    </div>

    <!-- Filtros de Búsqueda -->
    <div class="filters glass-panel">
      <form @submit.prevent="handleSearch" class="filter-form">
        <input type="text" v-model="filters.title" placeholder="Buscar por título..." class="form-input" />
        <SearchSelect 
            v-model="filters.artist_id"
            :initialName="filterArtistName"
            :searchFn="artistService.search"
            :formatDisplay="(a) => a.artist_name || a.name"
            placeholder="Todos los artistas (Búsqueda)"
            @select="(item) => filterArtistName = (item.artist_name || item.name)"
        />
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
            @delete="handleDelete(album.id, album.title)"
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
    <Pagination 
      v-if="!loading && !error && albums.length > 0"
      :currentPage="pagination.page"
      :totalPages="pagination.total_pages"
      :totalItems="pagination.total_items"
      @page-change="changePage"
    />

    <!-- Modal Form for Album Create/Edit -->
    <AlbumFormModal
        :isOpen="isModalOpen"
        :album="selectedAlbum"
        @close="isModalOpen = false"
        @saved="handleSaved"
    />

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
                        {{ track.title }}<span v-if="track.artists && track.artists.length > 0" class="text-secondary"> - {{ track.artists.map(a => a.name).join(', ') }}</span>
                    </span>
                    <button class="icon-btn text-danger" @click="removeTrack(track.song_id)">🗑️</button>
                </li>
            </ul>
        </div>

        <form @submit.prevent="addTrack" class="add-track-form border-t pt-4">
            <h4 class="mb-2 text-sm text-secondary">Agregar Pista</h4>
            <div class="form-group">
                <label>Canción</label>
                <div class="flex gap-2">
                    <div class="flex-1">
                        <SearchSelect 
                            v-model="trackForm.song_id"
                            :initialName="trackForm.song_title"
                            :searchFn="songService.search"
                            :formatDisplay="formatSongDisplay"
                            placeholder="Busca una canción..."
                            @select="(item) => trackForm.song_title = formatSongDisplay(item)"
                        />
                    </div>
                    <button type="button" class="btn btn-secondary icon-btn" title="Crear Nueva Canción" @click="openCreateTrackSong">
                        ➕🎵
                    </button>
                </div>
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

    <!-- Confirm Delete Modal -->
    <ConfirmDeleteModal 
      :isOpen="isDeleteModalOpen" 
      :itemName="itemToDeleteName"
      @close="isDeleteModalOpen = false"
      @confirm="executeDelete"
    />

    <!-- Song Form Modal for Track Management -->
    <SongFormModal 
        :isOpen="isSongModalOpen" 
        @close="isSongModalOpen = false" 
        @saved="handleSongSaved" 
    />
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

.filter-form .form-input,
.filter-form:deep(.search-select-wrapper) {
  flex: 1;
}

.shrink-btn {
  flex-shrink: 0;
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
    padding: 0.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
}

.artist-row {
    background: rgba(255,255,255,0.03);
    padding: 0.5rem;
    border-radius: var(--radius-sm);
    border: 1px solid var(--border-color);
}
</style>
