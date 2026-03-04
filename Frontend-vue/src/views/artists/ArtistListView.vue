<script setup>
import { ref, onMounted, onUnmounted, reactive, watch } from 'vue';
import { artistService } from '../../services/artist.service';
import { useToast } from '../../composables/useToast';
import Icon from '../../components/common/Icon.vue';
import ArtistCard from '../../components/artists/ArtistCard.vue';
import ArtistFormModal from '../../components/artists/ArtistFormModal.vue';
import ConfirmDeleteModal from '../../components/common/ConfirmDeleteModal.vue';
import Pagination from '../../components/common/Pagination.vue';

const artists = ref([]);
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
    // Calculate columns based on minmax(240px) + gap(24px)
    let cols = Math.floor((width + 24) / 264);
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
            pagination.page = 1; // Reset to page 1 to avoid offset bugs
            fetchArtists();
        }
    }, 300);
};

const filters = reactive({
  name: '',
  genre: '',
  country: ''
});

const fetchArtists = async () => {
  loading.value = true;
  try {
    const params = {
      page: pagination.page,
      limit: pagination.limit,
      name: filters.name,
      genre: filters.genre,
      country: filters.country
    };
    const response = await artistService.getPaginated(params);
    // Standard pagination response format
    artists.value = response.data || [];
    pagination.total_pages = response.total_pages || 1;
    pagination.total_items = response.total_items || 0;
  } catch (err) {
    error.value = 'No se pudieron cargar los artistas';
    console.error(err);
  } finally {
    loading.value = false;
  }
};

const handleSearch = () => {
  pagination.page = 1;
  fetchArtists();
};

const changePage = (newPage) => {
  if (newPage >= 1 && newPage <= pagination.total_pages) {
    pagination.page = newPage;
    fetchArtists();
  }
};

// Modal state
const isModalOpen = ref(false);
const selectedArtist = ref(null);

const openCreateModal = () => {
  selectedArtist.value = null;
  isModalOpen.value = true;
};

const handleEdit = (artist) => {
  selectedArtist.value = artist;
  isModalOpen.value = true;
};

const handleSaved = async () => {
    await fetchArtists();
};

const isDeleteModalOpen = ref(false);
const itemToDeleteId = ref(null);
const itemToDeleteName = ref('');

const handleDelete = (id, name) => {
  itemToDeleteId.value = id;
  itemToDeleteName.value = name || 'este artista';
  isDeleteModalOpen.value = true;
};

const executeDelete = async () => {
    try {
      await artistService.delete(itemToDeleteId.value);
      artists.value = artists.value.filter(a => a.id !== itemToDeleteId.value);
      isDeleteModalOpen.value = false;
      toast.success('Artista eliminado exitosamente');
    } catch(err) {
      toast.handleApiError(err, 'Error eliminando el artista');
    }
};

onMounted(() => {
  pagination.limit = calculateLimit();
  fetchArtists();
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});
</script>

<template>
  <div class="artists-view" ref="viewContainer">
    <div class="header">
      <h1 class="gradient-text">Artistas</h1>
      <button class="btn btn-primary" @click="openCreateModal" title="Crear un nuevo artista">
        <Icon name="user-plus" /> Nuevo Artista
      </button>
    </div>

    <!-- Filtros de Búsqueda -->
    <div class="filters glass-panel">
      <form @submit.prevent="handleSearch" class="filter-form">
        <input type="text" v-model="filters.name" placeholder="Buscar por nombre..." class="form-input" />
        <input type="text" v-model="filters.genre" placeholder="Filtro de género" class="form-input" />
        <input type="text" v-model="filters.country" placeholder="Filtro de país" class="form-input" />
        <button type="submit" class="btn btn-primary shrink-btn" title="Buscar artistas">
          <Icon name="search" /> Buscar
        </button>
      </form>
    </div>
    
    <div v-if="loading" class="loading">Cargando artistas...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else class="grid">
      <ArtistCard 
        v-for="artist in artists" 
        :key="artist.id" 
        :artist="artist"
        @edit="handleEdit"
        @delete="handleDelete(artist.id, artist.name)"
      />
    </div>

    <!-- Controles de Paginación -->
    <Pagination 
      v-if="!loading && !error && artists.length > 0"
      :currentPage="pagination.page"
      :totalPages="pagination.total_pages"
      :totalItems="pagination.total_items"
      @page-change="changePage"
    />

    <!-- Modal Form -->
    <ArtistFormModal 
        :isOpen="isModalOpen"
        :artist="selectedArtist"
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

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 1.5rem;
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
