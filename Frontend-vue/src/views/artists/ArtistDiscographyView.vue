<script setup>
import { ref, onMounted, onUnmounted, reactive, watch, computed } from 'vue';
import { useRoute } from 'vue-router';
import { albumService } from '../../services/album.service';
import { artistService } from '../../services/artist.service';
import AlbumCard from '../../components/albums/AlbumCard.vue';
import Pagination from '../../components/common/Pagination.vue';
import Breadcrumbs from '../../components/common/Breadcrumbs.vue';
import Icon from '../../components/common/Icon.vue';

const route = useRoute();
const artistId = computed(() => route.params.id);

const albums = ref([]);
const artist = ref(null);
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
            pagination.page = 1; 
            fetchAlbums();
        }
    }, 300);
};

const filters = reactive({
  title: '',
  type: ''
});

const fetchArtistInfo = async () => {
    try {
        const resp = await artistService.getById(artistId.value);
        artist.value = resp.data || resp;
    } catch(err) {
        console.error('Error fetching artist', err);
    }
};

const breadcrumbItems = computed(() => {
  return [
    { label: 'Inicio', to: '/' },
    { label: 'Artistas', to: '/artists' },
    { label: artist.value ? artist.value.name : 'Cargando...', to: `/artists/${artistId.value}` },
    { label: 'Discografía' }
  ];
});

const fetchAlbums = async () => {
  loading.value = true;
  try {
    const params = {
      artist_id: artistId.value,
      page: pagination.page,
      limit: pagination.limit,
      title: filters.title,
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

onMounted(() => {
  fetchArtistInfo();
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
    <Breadcrumbs :items="breadcrumbItems" />
    
    <div class="header">
      <h1 class="gradient-text">Discografía</h1>
    </div>

    <!-- Filtros de Búsqueda -->
    <div class="filters glass-panel">
      <form @submit.prevent="handleSearch" class="filter-form">
        <input type="text" v-model="filters.title" placeholder="Buscar por título..." class="form-input" />
        <select v-model="filters.type" class="form-input">
            <option value="">Todos los tipos</option>
            <option value="LP">Álbumes (LP)</option>
            <option value="EP">EPs</option>
            <option value="Single">Singles</option>
        </select>
        <button type="submit" class="btn btn-primary shrink-btn" title="Buscar álbumes">
          <Icon name="search" /> Buscar
        </button>
      </form>
    </div>
    
    <div v-if="loading" class="loading">Cargando discografía...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else class="grid">
      <div v-for="album in albums" :key="album.id" class="album-card-wrapper">
          <AlbumCard 
            :album="album"
            :readonly="true"
          />
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
  </div>
</template>

<style scoped>
.albums-view {
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

.album-card-wrapper {
    display: flex;
    flex-direction: column;
}
</style>
