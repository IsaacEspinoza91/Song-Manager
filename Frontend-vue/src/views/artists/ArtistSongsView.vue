<script setup>
import { ref, onMounted, reactive, computed } from 'vue';
import { useRoute } from 'vue-router';
import { songService } from '../../services/song.service';
import { artistService } from '../../services/artist.service';
import SongItem from '../../components/songs/SongItem.vue';
import Pagination from '../../components/common/Pagination.vue';
import Breadcrumbs from '../../components/common/Breadcrumbs.vue';
import SongFormModal from '../../components/songs/SongFormModal.vue';
import ConfirmDeleteModal from '../../components/common/ConfirmDeleteModal.vue';
import { useToast } from '../../composables/useToast';

const toast = useToast();

const route = useRoute();
const artistId = computed(() => route.params.id);

const songs = ref([]);
const artist = ref(null);
const loading = ref(true);
const error = ref(null);

const pagination = reactive({
  page: 1,
  limit: 15,
  total_pages: 1,
  total_items: 0
});

const filters = reactive({
  title: ''
});

const fetchArtistInfo = async () => {
    try {
        const resp = await artistService.getById(artistId.value);
        artist.value = resp.data || resp;
    } catch(err) {
        toast.handleApiError(err, 'Error al obtener artista');
    }
};

const breadcrumbItems = computed(() => {
  return [
    { label: 'Inicio', to: '/' },
    { label: 'Artistas', to: '/artists' },
    { label: artist.value ? artist.value.name : 'Cargando...', to: `/artists/${artistId.value}` },
    { label: 'Canciones' }
  ];
});

const fetchSongs = async () => {
  loading.value = true;
  try {
    const params = {
      artist_id: artistId.value,
      page: pagination.page,
      limit: pagination.limit,
      title: filters.title
    };
    const response = await songService.getPaginated(params);
    songs.value = response.data || [];
    pagination.total_pages = response.total_pages || 1;
    pagination.total_items = response.total_items || 0;
  } catch (err) {
    toast.handleApiError(err, 'No se pudieron cargar las canciones');
    error.value = 'No se pudieron cargar las canciones';
  } finally {
    loading.value = false;
  }
};

const handleSearch = () => {
  pagination.page = 1;
  fetchSongs();
};

const changePage = (newPage) => {
  if (newPage >= 1 && newPage <= pagination.total_pages) {
    pagination.page = newPage;
    fetchSongs();
  }
};

// ========================
// SONG CRUD (EDIT & DELETE)
// ========================
const isSongModalOpen = ref(false);
const selectedSong = ref(null);

const handleEditSong = async (songStub) => {
  try {
      const fullSongResp = await songService.getById(songStub.id);
      selectedSong.value = fullSongResp.data || fullSongResp;
      isSongModalOpen.value = true;
  } catch (err) {
      toast.handleApiError(err, 'Error al obtener datos de la canción');
  }
};

const handleSongSaved = async () => {
    await fetchSongs();
};

const isDeleteModalOpen = ref(false);
const itemToDeleteId = ref(null);
const itemToDeleteName = ref('');

const handleDeleteSong = (id, title) => {
  itemToDeleteId.value = id;
  itemToDeleteName.value = title || 'esta canción';
  isDeleteModalOpen.value = true;
};

const executeDelete = async () => {
    try {
      await songService.delete(itemToDeleteId.value);
      songs.value = songs.value.filter(s => s.id !== itemToDeleteId.value);
      isDeleteModalOpen.value = false;
      toast.success('Canción eliminada exitosamente');
    } catch(err) {
      toast.handleApiError(err, 'Error eliminando la canción');
    }
};

onMounted(() => {
  fetchArtistInfo();
  fetchSongs();
});
</script>

<template>
  <div class="songs-view">
    <Breadcrumbs :items="breadcrumbItems" />
    
    <div class="header">
      <h1 class="gradient-text">Todas las Canciones</h1>
    </div>

    <!-- Filtro de Búsqueda (Sólo título) -->
    <div class="filters glass-panel">
      <form @submit.prevent="handleSearch" class="filter-form">
        <input type="text" v-model="filters.title" placeholder="Buscar por título de canción..." class="form-input" />
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
        :readonly="false"
        @edit="handleEditSong"
        @delete="handleDeleteSong(song.id, song.title)"
      />
      <div v-if="songs.length === 0" class="empty-state">
        No se encontraron canciones.
      </div>
    </div>

    <!-- Controles de Paginación -->
    <Pagination 
      v-if="!loading && !error && songs.length > 0"
      :currentPage="pagination.page"
      :totalPages="pagination.total_pages"
      :totalItems="pagination.total_items"
      @page-change="changePage"
    />

    <!-- Song Edit Modal -->
    <SongFormModal 
        :isOpen="isSongModalOpen" 
        :song="selectedSong" 
        @close="isSongModalOpen = false" 
        @saved="handleSongSaved" 
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
.songs-view {
  padding-bottom: 2rem;
}

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
  font-size: 1.25rem;
}

.error { color: var(--danger); }

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
</style>
