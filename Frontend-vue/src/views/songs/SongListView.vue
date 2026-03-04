<script setup>
import { ref, onMounted, reactive, watch } from 'vue';
import { useRoute } from 'vue-router';
import { songService } from '../../services/song.service';
import { artistService } from '../../services/artist.service';
import { useToast } from '../../composables/useToast';
import SongItem from '../../components/songs/SongItem.vue';
import SongFormModal from '../../components/songs/SongFormModal.vue';
import ConfirmDeleteModal from '../../components/common/ConfirmDeleteModal.vue';
import SearchSelect from '../../components/common/SearchSelect.vue';
import Pagination from '../../components/common/Pagination.vue';

const toast = useToast();

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

const filterArtistName = ref('');

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

const changePage = (newPage) => {
  if (newPage >= 1 && newPage <= pagination.total_pages) {
    pagination.page = newPage;
    fetchSongs();
  }
};

// Modal state
const isModalOpen = ref(false);
const selectedSong = ref(null);

const openCreateModal = async () => {
  selectedSong.value = null;
  isModalOpen.value = true;
};

const handleEdit = async (song) => {
  selectedSong.value = song;
  isModalOpen.value = true;
};

const handleSaved = async () => {
    await fetchSongs();
};

const isDeleteModalOpen = ref(false);
const itemToDeleteId = ref(null);
const itemToDeleteName = ref('');

const handleDelete = (id, title) => {
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
  fetchSongs();
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
        <SearchSelect 
            v-model="filters.artist_id"
            :initialName="filterArtistName"
            :searchFn="artistService.search"
            :formatDisplay="(a) => a.artist_name || a.name"
            placeholder="Todos los artistas (Búsqueda)"
            @select="(item) => filterArtistName = (item.artist_name || item.name)"
        />
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
    <Pagination 
      v-if="!loading && !error && songs.length > 0"
      :currentPage="pagination.page"
      :totalPages="pagination.total_pages"
      :totalItems="pagination.total_items"
      @page-change="changePage"
    />

    <!-- Modal Form -->
    <SongFormModal 
        :isOpen="isModalOpen"
        :song="selectedSong"
        @close="isModalOpen = false"
        @saved="handleSaved"
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

.filter-form .form-input,
.filter-form:deep(.search-select-wrapper) {
  flex: 1;
}

.shrink-btn {
  flex-shrink: 0;
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
