<script setup>
import { ref, onMounted, reactive, watch } from 'vue';
import { useRoute } from 'vue-router';
import { songService } from '../../services/song.service';
import { artistService } from '../../services/artist.service';
import SongItem from '../../components/songs/SongItem.vue';
import Modal from '../../components/common/Modal.vue';
import ConfirmDeleteModal from '../../components/common/ConfirmDeleteModal.vue';

const songs = ref([]);
const loading = ref(true);
const error = ref(null);

const route = useRoute();

const pagination = reactive({
  page: 1,
  limit: 10,
  total_pages: 1,
  total_items: 0
});

const filters = reactive({
  title: '',
  artist_name: '',
  artist_id: route.query.artist_id || ''
});

const searchArtists = ref([]); // For the combo box filter

const fetchSongs = async () => {
  loading.value = true;
  try {
    const params = {
      page: pagination.page,
      limit: pagination.limit,
      title: filters.title,
      artist_name: filters.artist_name,
      artist_id: filters.artist_id
    };
    const response = await songService.getPaginated(params);
    songs.value = response.data || [];
    pagination.total_pages = response.total_pages || 1;
    pagination.total_items = response.total_items || 0;
  } catch (err) {
    error.value = 'No se pudieron cargar las canciones';
    console.error(err);
  } finally {
    loading.value = false;
  }
};

const handleSearch = () => {
  pagination.page = 1;
  fetchSongs();
};

const changePage = (delta) => {
  const newPage = pagination.page + delta;
  if (newPage >= 1 && newPage <= pagination.total_pages) {
    pagination.page = newPage;
    fetchSongs();
  }
};

// Also load artists for the filter combo box initially
const loadArtistsForFilter = async () => {
    try {
        const resp = await artistService.getAll();
        searchArtists.value = resp.data || resp;
        if (searchArtists.value.data) searchArtists.value = searchArtists.value.data;
    } catch (err) {
        console.error('Failed to load artists for filter', err);
    }
};

// Modal state
const isModalOpen = ref(false);
const isEditing = ref(false);
const formError = ref(null);
const availableArtists = ref([]);
const songForm = reactive({
  id: null,
  title: '',
  duration: 0,
  artists: []
});

const fetchAvailableArtists = async () => {
  try {
    const resp = await artistService.getAll();
    availableArtists.value = resp.data || resp; // API returns directly or inside data
    if (availableArtists.value.data) availableArtists.value = availableArtists.value.data;
  } catch(err) {
    console.error('Could not load artists', err);
  }
};

const openCreateModal = async () => {
  isEditing.value = false;
  Object.assign(songForm, { id: null, title: '', duration: 0, artists: [{ artist_id: '', role: 'main' }] });
  formError.value = null;
  if(availableArtists.value.length === 0) await fetchAvailableArtists();
  isModalOpen.value = true;
};

const handleEdit = async (song) => {
  isEditing.value = true;
  
  const initialArtists = song.artists && song.artists.length > 0 
    ? song.artists.map(a => ({ artist_id: a.id || a.artist_id, role: a.role }))
    : [{ artist_id: '', role: 'main' }];

  Object.assign(songForm, {
    id: song.id,
    title: song.title,
    duration: song.duration,
    artists: [...initialArtists]
  });
  formError.value = null;
  if(availableArtists.value.length === 0) await fetchAvailableArtists();
  isModalOpen.value = true;
};

const saveSong = async () => {
  try {
    formError.value = null;
    let payload = {
      title: songForm.title,
      duration: Number(songForm.duration)
    };
    
    // Filter out rows without a selected artist
    const validArtists = songForm.artists
        .filter(a => a.artist_id)
        .map(a => ({ artist_id: Number(a.artist_id), role: a.role }));
        
    if (validArtists.length > 0) {
       payload.artists = validArtists;
    }

    if (isEditing.value) {
      await songService.update(songForm.id, payload);
    } else {
      await songService.create(payload);
    }
    isModalOpen.value = false;
    await fetchSongs();
  } catch (err) {
    console.error(err);
    formError.value = 'Error al guardar la canción.';
  }
};

const isDeleteModalOpen = ref(false);
const itemToDeleteId = ref(null);
const itemToDeleteName = ref('');

const handleDelete = (id, title) => {
  itemToDeleteId.value = id;
  itemToDeleteName.value = title || 'esta canción';
  isDeleteModalOpen.value = true;
};

const removeArtist = (index) => {
    songForm.artists.splice(index, 1);
};

const addArtist = () => {
    songForm.artists.push({ artist_id: '', role: 'main' });
};

const executeDelete = async () => {
    try {
      await songService.delete(itemToDeleteId.value);
      songs.value = songs.value.filter(s => s.id !== itemToDeleteId.value);
      isDeleteModalOpen.value = false;
    } catch(err) {
      console.error(err);
      alert('Error eliminando la canción');
    }
};

onMounted(() => {
  fetchSongs();
  loadArtistsForFilter();
});
</script>

<template>
  <div class="songs-view">
    <div class="header">
      <h1 class="gradient-text">Canciones</h1>
      <button class="btn btn-primary" @click="openCreateModal">Nueva Canción</button>
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
        <button type="submit" class="btn btn-primary shrink-btn">Buscar</button>
      </form>
    </div>
    
    <div v-if="loading" class="loading">Cargando canciones...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else class="song-list glass-panel">
      <div class="list-header">
        <div class="h-index">#</div>
        <div class="h-title">Título</div>
        <div class="h-duration">Duración</div>
        <div class="h-actions"></div>
      </div>
      <SongItem 
        v-for="(song, index) in songs" 
        :key="song.id" 
        :song="song"
        :index="index"
        @edit="handleEdit"
        @delete="handleDelete(song.id, song.title)"
      />
      <div v-if="songs.length === 0" class="empty-state">
        No se encontraron canciones.
      </div>
    </div>

    <!-- Controles de Paginación -->
    <div v-if="!loading && !error && songs.length > 0" class="pagination">
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
    <Modal :isOpen="isModalOpen" @close="isModalOpen = false" :title="isEditing ? 'Editar Canción' : 'Nueva Canción'">
      <form @submit.prevent="saveSong">
        <div v-if="formError" class="error-msg">{{ formError }}</div>
        
        <div class="form-group">
          <label>Título</label>
          <input type="text" v-model="songForm.title" class="form-input" required />
        </div>
        
        <div class="form-group">
          <label>Duración (Segundos)</label>
          <input type="number" v-model="songForm.duration" class="form-input" required min="1" />
        </div>

        <div class="form-group mb-4">
          <div class="flex justify-between items-center mb-2">
            <label class="mb-0">Artistas (Opcional)</label>
            <button type="button" class="btn btn-secondary btn-sm" @click="addArtist">+ Agregar Artista</button>
          </div>
          
          <div v-for="(artistEntry, index) in songForm.artists" :key="index" class="artist-row flex gap-2 mb-2 items-start">
              <div class="flex-1">
                  <select v-model="artistEntry.artist_id" class="form-input" required>
                    <option value="">Selecciona un Artista</option>
                    <option v-for="artist in availableArtists" :key="artist.id" :value="artist.id">
                      {{ artist.name }}
                    </option>
                  </select>
              </div>
              <div class="w-1/3">
                  <select v-model="artistEntry.role" class="form-input">
                      <option value="main">Main (Principal)</option>
                      <option value="ft">Featuring (Invitado)</option>
                      <option value="producer">Productor</option>
                  </select>
              </div>
              <button 
                  v-if="songForm.artists.length > 1" 
                  type="button" 
                  class="btn btn-danger icon-btn" 
                  @click="removeArtist(index)">
                  🗑️
              </button>
          </div>
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn-secondary" @click="isModalOpen = false">Cancelar</button>
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

.song-list {
  padding: 1rem;
}

.list-header {
  display: flex;
  padding: 0.5rem 1.5rem;
  border-bottom: 1px solid var(--border-color);
  margin-bottom: 1rem;
  color: var(--text-muted);
  font-size: 0.85rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-weight: 600;
}

.h-index { width: 30px; }
.h-title { flex: 1; }
.h-duration { margin: 0 2rem; }
.h-actions { width: 60px; }

.loading, .error, .empty-state {
  text-align: center;
  padding: 3rem;
  color: var(--text-secondary);
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
  margin-top: 2rem;
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

.flex { display: flex; }
.justify-between { justify-content: space-between; }
.items-center { align-items: center; }
.items-start { align-items: flex-start; }
.gap-2 { gap: 0.5rem; }
.mb-2 { margin-bottom: 0.5rem; }
.mb-4 { margin-bottom: 1rem; }
.mb-0 { margin-bottom: 0 !important; }
.flex-1 { flex: 1; }
.w-1\/3 { width: 33.333333%; }

.artist-row {
    background: rgba(255,255,255,0.03);
    padding: 0.5rem;
    border-radius: var(--radius-sm);
    border: 1px solid var(--border-color);
}

.icon-btn {
    padding: 0.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
}
</style>
