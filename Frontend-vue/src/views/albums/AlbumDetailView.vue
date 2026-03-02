<script setup>
import { ref, onMounted, computed, reactive } from 'vue';
import { useRoute } from 'vue-router';
import { albumService } from '../../services/album.service';
import { artistService } from '../../services/artist.service';
import { songService } from '../../services/song.service';
import Breadcrumbs from '../../components/common/Breadcrumbs.vue';
import SongItem from '../../components/songs/SongItem.vue';
import Modal from '../../components/common/Modal.vue';
import ConfirmDeleteModal from '../../components/common/ConfirmDeleteModal.vue';

const route = useRoute();
const albumId = route.params.id;

const loading = ref(true);
const error = ref(null);
const album = ref(null);
const primaryArtist = ref(null);

const breadcrumbItems = computed(() => {
  return [
    { label: 'Inicio', to: '/' },
    { label: 'Álbumes', to: '/albums' },
    { label: album.value ? album.value.title : 'Cargando...' }
  ];
});

const formatDate = (dateString) => {
    if (!dateString) return 'Desconocido';
    return new Date(dateString).toLocaleDateString();
};

const fetchData = async () => {
  loading.value = true;
  try {
    const albumResp = await albumService.getById(albumId);
    album.value = albumResp.data || albumResp;

    // Fetch primary artist for the photo
    if (album.value.artists && album.value.artists.length > 0) {
        const pArtist = album.value.artists.find(a => a.is_primary) || album.value.artists[0];
        const artistResp = await artistService.getById(pArtist.artist_id || pArtist.id);
        primaryArtist.value = artistResp.data || artistResp;
    }

  } catch(err) {
    console.error(err);
    error.value = 'No se pudieron cargar los datos del álbum.';
  } finally {
    loading.value = false;
  }
};

// ========================
// SONG CRUD (EDIT & DELETE)
// ========================
const isSongModalOpen = ref(false);
const songFormError = ref(null);
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
    availableArtists.value = resp.data || resp;
    if (availableArtists.value.data) availableArtists.value = availableArtists.value.data;
  } catch(err) {
    console.error('Could not load artists', err);
  }
};

const handleEditSong = async (song) => {
  const initialArtists = song.artists && song.artists.length > 0 
    ? song.artists.map(a => ({ artist_id: a.id || a.artist_id, role: a.role }))
    : [{ artist_id: '', role: 'main' }];

  Object.assign(songForm, {
    id: song.id,
    title: song.title,
    duration: song.duration,
    artists: [...initialArtists]
  });
  songFormError.value = null;
  if(availableArtists.value.length === 0) await fetchAvailableArtists();
  isSongModalOpen.value = true;
};

const saveSong = async () => {
  try {
    songFormError.value = null;
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
    
    // We only update here, no creation from this view
    await songService.update(songForm.id, payload);
    isSongModalOpen.value = false;
    
    // Refresh album data to get updated tracks
    await fetchData();
  } catch (err) {
    console.error(err);
    songFormError.value = 'Error al editar la canción.';
  }
};

// Add/Remove artist handlers
const removeArtist = (index) => {
    songForm.artists.splice(index, 1);
};

const addArtist = () => {
    songForm.artists.push({ artist_id: '', role: 'main' });
};

// Delete Confirmation
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
        await fetchData(); // Refresh album data
        isDeleteModalOpen.value = false;
    } catch(err) {
        console.error(err);
        alert('Error en la eliminación');
    }
};

onMounted(() => {
  fetchData();
});
</script>

<template>
  <div class="album-detail-view">
    <Breadcrumbs :items="breadcrumbItems" />
    
    <div v-if="loading" class="loading">Cargando álbum...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="album">
      
      <!-- Album Hero Section -->
      <div class="hero glass-panel">
          <div class="hero-bg" :style="{ backgroundImage: `url(${album.cover_url || ''})` }"></div>
          <div class="hero-content">
              <div class="cover-wrapper">
                  <img v-if="album.cover_url" :src="album.cover_url" :alt="album.title" class="cover-img" />
                  <div v-else class="placeholder-cover">{{ album.title.charAt(0) }}</div>
              </div>
              <div class="hero-info">
                  <span class="type-badge">{{ album.type }}</span>
                  <h1 class="album-title">{{ album.title }}</h1>
                  
                  <div class="artist-block" v-if="primaryArtist">
                      <div class="artist-thumb-wrapper">
                          <img v-if="primaryArtist.image_url" :src="primaryArtist.image_url" :alt="primaryArtist.name" class="artist-thumb" />
                          <div v-else class="placeholder-thumb">{{ primaryArtist.name.charAt(0) }}</div>
                      </div>
                      <router-link :to="`/artists/${primaryArtist.id}`" class="artist-link">
                          {{ primaryArtist.name }}
                      </router-link>
                      <span class="meta-dot">•</span>
                      <span class="release-year">{{ new Date(album.release_date).getFullYear() }}</span>
                  </div>
                  <div v-else class="artist-block">
                     <span class="text-muted">Artista Desconocido</span>
                  </div>

                  <p class="album-dates">
                      Añadido: {{ formatDate(album.created_at) }}<br/>
                      Última actualización: {{ formatDate(album.updated_at) }}
                  </p>
              </div>
          </div>
      </div>

      <!-- Tracks List -->
      <div class="tracks-section glass-panel p-4 mt-4">
          <h2 class="section-title mb-4">Pistas</h2>
          <div v-if="!album.tracks || album.tracks.length === 0" class="empty-state">
              No hay pistas registradas en este álbum.
          </div>
          <div v-else class="tracks-list">
              <!-- SongItem handling full CRUD mappings -->
              <SongItem 
                v-for="track in album.tracks" 
                :key="track.song_id" 
                :song="{ id: track.song_id, title: track.title, duration: track.duration, artists: album.artists }" 
                :index="track.track_number - 1" 
                :readonly="false" 
                @edit="handleEditSong"
                @delete="handleDeleteSong(track.song_id, track.title)"
              />
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

        <div class="form-group mb-4">
          <div class="flex justify-between items-center mb-2">
            <label class="mb-0">Artistas</label>
            <button type="button" class="btn btn-secondary btn-sm" @click="addArtist">+ Agregar Artista</button>
          </div>
          
          <div v-for="(artistEntry, index) in songForm.artists" :key="index" class="artist-row flex gap-2 mb-2 items-start">
              <div class="flex-1">
                  <select v-model="artistEntry.artist_id" class="form-input" required>
                    <option value="">Selecciona un Artista</option>
                    <option v-for="a in availableArtists" :key="a.id" :value="a.id">
                      {{ a.name }}
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
          <button type="button" class="btn btn-secondary" @click="isSongModalOpen = false">Cancelar</button>
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
.album-detail-view {
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
  margin-bottom: 2.5rem;
}

.hero-bg {
    position: absolute;
    top: 0; left: 0; right: 0; bottom: 0;
    background-size: cover;
    background-position: center;
    filter: blur(30px) brightness(0.3);
    z-index: 0;
}

.hero-content {
    position: relative;
    z-index: 1;
    display: flex;
    gap: 2.5rem;
    align-items: flex-end;
    width: 100%;
}

.cover-wrapper {
    width: 250px;
    height: 250px;
    border-radius: var(--radius-md);
    overflow: hidden;
    box-shadow: 0 10px 30px rgba(0,0,0,0.5);
    flex-shrink: 0;
}

.cover-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.placeholder-cover {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 5rem;
    font-weight: 700;
    background: linear-gradient(135deg, var(--bg-tertiary), var(--bg-secondary));
    color: var(--text-muted);
}

.hero-info {
    flex: 1;
    display: flex;
    flex-direction: column;
}

.type-badge {
    text-transform: uppercase;
    font-size: 0.8rem;
    font-weight: 700;
    letter-spacing: 0.05em;
    margin-bottom: 0.5rem;
    color: var(--text-secondary);
}

.album-title {
    font-size: 4rem;
    font-weight: 800;
    margin-bottom: 1rem;
    text-shadow: 2px 2px 4px rgba(0,0,0,0.5);
    line-height: 1;
}

.artist-block {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 1.5rem;
}

.artist-thumb-wrapper {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    overflow: hidden;
}

.artist-thumb {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.placeholder-thumb {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--accent-primary);
    color: var(--bg-primary);
    font-weight: bold;
    font-size: 0.8rem;
}

.artist-link {
    font-weight: 700;
    color: var(--text-primary);
    text-decoration: none;
    font-size: 1.1rem;
}

.artist-link:hover {
    text-decoration: underline;
}

.meta-dot {
    color: var(--text-muted);
    font-size: 0.8rem;
}

.release-year {
    color: var(--text-secondary);
    font-size: 1.1rem;
}

.album-dates {
  font-size: 0.85rem;
  color: var(--text-muted);
  border-top: 1px solid rgba(255,255,255,0.1);
  padding-top: 1rem;
  margin-top: 0.5rem;
}

.section-title {
    font-size: 1.5rem;
    font-weight: 600;
    margin-bottom: 1.5rem;
    color: var(--text-primary);
}

.mt-4 { margin-top: 1rem; }
.mb-4 { margin-bottom: 1rem; }
.p-4 { padding: 1.5rem; }
.text-muted { color: var(--text-muted); }

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

/* Responsive Adjustments */
@media (max-width: 768px) {
    .hero {
        align-items: flex-start;
    }
    .hero-content {
        flex-direction: column;
        align-items: center;
        text-align: center;
    }
    
    .cover-wrapper {
        width: 200px;
        height: 200px;
    }
    
    .album-title {
        font-size: 2.5rem;
    }

    .artist-block {
        justify-content: center;
    }
}
</style>
