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
import SearchSelect from '../../components/common/SearchSelect.vue';

const route = useRoute();
const albumId = route.params.id;

const loading = ref(true);
const error = ref(null);
const album = ref(null);
const primaryArtists = ref([]);

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

    // Fetch primary artists for the photo
    if (album.value.artists && album.value.artists.length > 0) {
        const primaryArtistRels = album.value.artists.filter(a => a.is_primary);
        const toFetch = primaryArtistRels.length > 0 ? primaryArtistRels : [album.value.artists[0]];
        
        const fetchedArtists = await Promise.all(
            toFetch.map(a => artistService.getById(a.artist_id || a.id))
        );
        primaryArtists.value = fetchedArtists.map(resp => resp.data || resp);
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
const songForm = reactive({
  id: null,
  title: '',
  duration: 0,
  artists: []
});

const formatSongDisplay = (song) => {
    const artistNames = song.artists ? song.artists.map(a => a.artist_name || a.name).join(', ') : '';
    return artistNames ? `${song.title} - ${artistNames}` : song.title;
};

const handleEditSong = async (songStub) => {
  try {
      // Fetch full song data to get accurate artist roles for this specific song
      const fullSongResp = await songService.getById(songStub.id);
      const song = fullSongResp.data || fullSongResp;
      
      const initialArtists = song.artists && song.artists.length > 0 
        ? song.artists.map(a => ({ artist_id: a.id || a.artist_id, artist_name: a.name || a.artist_name, role: a.role }))
        : [{ artist_id: '', artist_name: '', role: 'main' }];

      Object.assign(songForm, {
        id: song.id,
        title: song.title,
        duration: song.duration,
        artists: [...initialArtists]
      });
      songFormError.value = null;
      isSongModalOpen.value = true;
  } catch (err) {
      console.error("Error fetching full song details for edit:", err);
      alert('Error al obtener datos de la canción para edición.');
  }
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
    songForm.artists.push({ artist_id: '', artist_name: '', role: 'main' });
};

// ========================
// TRACK ADDITION
// ========================
const isTrackModalOpen = ref(false);
const trackFormError = ref(null);
const trackForm = reactive({
    song_id: '',
    song_title: '',
    track_number: 1
});

const openAddTrackModal = async () => {
    trackForm.song_id = '';
    trackForm.song_title = '';
    // Suggest the next track number
    trackForm.track_number = album.value.tracks ? album.value.tracks.length + 1 : 1;
    trackFormError.value = null;
    isTrackModalOpen.value = true;
};

const appendTrack = async () => {
    try {
        trackFormError.value = null;
        if (!trackForm.song_id) {
            trackFormError.value = 'Debe seleccionar una canción.';
            return;
        }

        const payload = {
            song_id: Number(trackForm.song_id),
            track_number: Number(trackForm.track_number)
        };

        await albumService.addTrack(albumId, payload);
        isTrackModalOpen.value = false;
        await fetchData(); // Refresh album details
    } catch (err) {
        console.error(err);
        trackFormError.value = 'Error al agregar la pista al álbum. Es posible que el número de pista ya exista o la canción ya esté agregada.';
    }
};

// Delete Confirmation
const isDeleteModalOpen = ref(false);
const itemToDeleteId = ref(null);
const itemToDeleteName = ref('');

const handleDeleteSong = (id, title) => {
  itemToDeleteId.value = id;
  // Specific wording requested by user
  itemToDeleteName.value = title ? `la pista "${title}" del álbum (la canción original no será eliminada)` : 'esta pista del álbum';
  isDeleteModalOpen.value = true;
};

const executeDelete = async () => {
    try {
        // Use removeTrack to purely remove the relationship, NOT delete the global song
        await albumService.removeTrack(albumId, itemToDeleteId.value);
        await fetchData(); // Refresh album data
        isDeleteModalOpen.value = false;
    } catch(err) {
        console.error(err);
        alert('Error al remover la pista');
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
                  
                  <div class="artist-block-list" v-if="primaryArtists && primaryArtists.length > 0">
                      <template v-for="(artist, index) in primaryArtists" :key="artist.id">
                          <div class="artist-item">
                              <div class="artist-thumb-wrapper">
                                  <img v-if="artist.image_url" :src="artist.image_url" :alt="artist.name" class="artist-thumb" />
                                  <div v-else class="placeholder-thumb">{{ artist.name.charAt(0) }}</div>
                              </div>
                              <router-link :to="`/artists/${artist.id}`" class="artist-link">
                                  {{ artist.name }}
                              </router-link>
                          </div>
                          <span class="meta-dot" v-if="index < primaryArtists.length - 1">,</span>
                      </template>
                      <span class="meta-dot">•</span>
                      <span class="release-year" v-if="album.release_date">{{ formatDate(album.release_date) }}</span>
                  </div>
                  <div v-else class="artist-block">
                     <span class="text-muted">Artista Desconocido</span>
                     <span class="meta-dot" v-if="album.release_date">•</span>
                     <span class="release-year" v-if="album.release_date">{{ formatDate(album.release_date) }}</span>
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
          <div class="flex justify-between items-center mb-4">
              <h2 class="section-title mb-0">Pistas</h2>
              <button class="btn btn-primary btn-sm" @click="openAddTrackModal">+ Agregar Pista</button>
          </div>
          <div v-if="!album.tracks || album.tracks.length === 0" class="empty-state">
              No hay pistas registradas en este álbum.
          </div>
          <div v-else class="tracks-list">
              <!-- SongItem handling full CRUD mappings -->
              <SongItem 
                v-for="track in album.tracks" 
                :key="track.song_id" 
                :song="{ id: track.song_id, title: track.title, duration: track.duration, artists: track.artists }" 
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
              <div class="flex-1" style="min-width: 0;">
                  <SearchSelect 
                      v-model="artistEntry.artist_id"
                      :initialName="artistEntry.artist_name"
                      :searchFn="artistService.search"
                      :formatDisplay="(a) => a.artist_name || a.name"
                      placeholder="Busca un artista..."
                      @select="(item) => artistEntry.artist_name = (item.artist_name || item.name)"
                  />
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
    
    <!-- Add Track Modal -->
    <Modal :isOpen="isTrackModalOpen" @close="isTrackModalOpen = false" title="Agregar Pista al Álbum">
      <form @submit.prevent="appendTrack">
        <div v-if="trackFormError" class="error-msg">{{ trackFormError }}</div>
        
        <div class="form-group">
            <label>Seleccionar Canción Existente</label>
            <SearchSelect 
                v-model="trackForm.song_id"
                :initialName="trackForm.song_title"
                :searchFn="songService.search"
                :formatDisplay="formatSongDisplay"
                placeholder="Busca una canción..."
                @select="(item) => trackForm.song_title = formatSongDisplay(item)"
            />
        </div>
        
        <div class="form-group">
            <label>Número de Pista</label>
            <input type="number" v-model="trackForm.track_number" class="form-input" required min="1" />
        </div>

        <div class="form-actions mt-4">
          <button type="button" class="btn btn-secondary" @click="isTrackModalOpen = false">Cancelar</button>
          <button type="submit" class="btn btn-primary">Agregar</button>
        </div>
      </form>
    </Modal>

    <!-- Confirm Delete Modal (Remove Track) -->
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

.artist-block-list {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 1.5rem;
    flex-wrap: wrap;
}

.artist-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
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
