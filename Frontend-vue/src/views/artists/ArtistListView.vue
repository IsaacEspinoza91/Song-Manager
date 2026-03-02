<script setup>
import { ref, onMounted, reactive, watch } from 'vue';
import { artistService } from '../../services/artist.service';
import ArtistCard from '../../components/artists/ArtistCard.vue';
import Modal from '../../components/common/Modal.vue';
import ConfirmDeleteModal from '../../components/common/ConfirmDeleteModal.vue';

const artists = ref([]);
const loading = ref(true);
const error = ref(null);

const pagination = reactive({
  page: 1,
  limit: 10,
  total_pages: 1,
  total_items: 0
});

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

const changePage = (delta) => {
  const newPage = pagination.page + delta;
  if (newPage >= 1 && newPage <= pagination.total_pages) {
    pagination.page = newPage;
    fetchArtists();
  }
};

// Modal state
const isModalOpen = ref(false);
const isEditing = ref(false);
const formError = ref(null);
const artistForm = reactive({
  id: null,
  name: '',
  genre: '',
  country: '',
  bio: '',
  image_url: ''
});

const openCreateModal = () => {
  isEditing.value = false;
  Object.assign(artistForm, { id: null, name: '', genre: '', country: '', bio: '', image_url: '' });
  formError.value = null;
  isModalOpen.value = true;
};

const handleEdit = (artist) => {
  isEditing.value = true;
  Object.assign(artistForm, artist);
  formError.value = null;
  isModalOpen.value = true;
};

const saveArtist = async () => {
  try {
    formError.value = null;
    if (isEditing.value) {
      const resp = await artistService.update(artistForm.id, artistForm);
      const index = artists.value.findIndex(a => a.id === artistForm.id);
      if (index !== -1) artists.value[index] = resp.data || resp;
    } else {
      const resp = await artistService.create(artistForm);
      artists.value.push(resp.data || resp);
    }
    isModalOpen.value = false;
    await fetchArtists(); // Refresh to ensure backend consistency
  } catch (err) {
    console.error(err);
    formError.value = 'Error al guardar el artista. Revisa la consola y conexión.';
  }
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
    } catch(err) {
      console.error(err);
      alert('Error eliminando el artista');
    }
};

onMounted(() => {
  fetchArtists();
});
</script>

<template>
  <div class="artists-view">
    <div class="header">
      <h1 class="gradient-text">Artistas</h1>
      <button class="btn btn-primary" @click="openCreateModal">Nuevo Artista</button>
    </div>

    <!-- Filtros de Búsqueda -->
    <div class="filters glass-panel">
      <form @submit.prevent="handleSearch" class="filter-form">
        <input type="text" v-model="filters.name" placeholder="Buscar por nombre..." class="form-input" />
        <input type="text" v-model="filters.genre" placeholder="Filtro de género" class="form-input" />
        <input type="text" v-model="filters.country" placeholder="Filtro de país" class="form-input" />
        <button type="submit" class="btn btn-primary shrink-btn">Buscar</button>
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
    <div v-if="!loading && !error && artists.length > 0" class="pagination">
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
    <Modal :isOpen="isModalOpen" @close="isModalOpen = false" :title="isEditing ? 'Editar Artista' : 'Nuevo Artista'">
      <form @submit.prevent="saveArtist">
        <div v-if="formError" class="error-msg">{{ formError }}</div>
        <div class="form-group">
          <label>Nombre</label>
          <input type="text" v-model="artistForm.name" class="form-input" required />
        </div>
        <div class="form-group">
          <label>Género</label>
          <input type="text" v-model="artistForm.genre" class="form-input" required />
        </div>
        <div class="form-group">
          <label>País</label>
          <input type="text" v-model="artistForm.country" class="form-input" required />
        </div>
        <div class="form-group">
          <label>Biografía (Opcional)</label>
          <textarea v-model="artistForm.bio" class="form-textarea" rows="3"></textarea>
        </div>
        <div class="form-group">
          <label>URL de Imagen (Opcional)</label>
          <input type="url" v-model="artistForm.image_url" class="form-input" />
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
</style>
