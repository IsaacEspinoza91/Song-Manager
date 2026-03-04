<script setup>
import { ref, reactive, watch } from 'vue';
import { albumService } from '../../services/album.service';
import { artistService } from '../../services/artist.service';
import { songService } from '../../services/song.service';
import Modal from '../common/Modal.vue';
import SearchSelect from '../common/SearchSelect.vue';
import ArtistFormModal from '../artists/ArtistFormModal.vue';
import SongFormModal from '../songs/SongFormModal.vue';

const props = defineProps({
  isOpen: { type: Boolean, required: true },
  album: { type: Object, default: null }, // Null to create, Object to edit
  initialArtist: { type: Object, default: null } // { id, name } to pre-fill artist
});

const emit = defineEmits(['close', 'saved']);

const formError = ref(null);
const isEditing = ref(false);

const albumForm = reactive({
  id: null,
  title: '',
  release_date: '',
  type: 'LP',
  cover_url: '',
  artists: [],
  showTracksForm: false, // For new albums only
  tracks: []
});

const resetForm = () => {
    formError.value = null;
    if (props.album) {
        // Editing existing album
        isEditing.value = true;
        
        const initialArtists = props.album.artists && props.album.artists.length > 0
            ? props.album.artists.map(a => ({ artist_id: a.id || a.artist_id, artist_name: a.name || a.artist_name, is_primary: a.is_primary }))
            : [{ artist_id: '', artist_name: '', is_primary: true }];

        Object.assign(albumForm, {
            id: props.album.id,
            title: props.album.title,
            release_date: props.album.release_date ? props.album.release_date.split('T')[0] : '', // format for date input
            type: props.album.type || 'LP',
            cover_url: props.album.cover_url || '',
            artists: [...initialArtists],
            showTracksForm: false,
            tracks: []
        });
    } else {
        // Creating new album
        isEditing.value = false;
        let initialArtists = [];
        if (props.initialArtist && props.initialArtist.id) {
            initialArtists.push({
                artist_id: props.initialArtist.id,
                artist_name: props.initialArtist.name || '',
                is_primary: true
            });
        } else {
            initialArtists.push({ artist_id: '', artist_name: '', is_primary: true });
        }
        
        Object.assign(albumForm, {
            id: null,
            title: '',
            release_date: '',
            type: 'LP',
            cover_url: '',
            artists: initialArtists,
            showTracksForm: false,
            tracks: []
        });
    }
};

watch(() => props.isOpen, (newVal) => {
    if (newVal) {
        resetForm();
    }
});

const addArtistEntry = () => {
    albumForm.artists.push({ artist_id: '', artist_name: '', is_primary: false });
};

const removeArtistEntry = (index) => {
    albumForm.artists.splice(index, 1);
};

const toggleTracksForm = () => {
    albumForm.showTracksForm = !albumForm.showTracksForm;
    if (albumForm.showTracksForm && albumForm.tracks.length === 0) {
        albumForm.tracks.push({ song_id: '', song_title: '', track_number: 1 });
    }
};

const addTrackEntry = () => {
    const nextNumber = albumForm.tracks.length + 1;
    albumForm.tracks.push({ song_id: '', song_title: '', track_number: nextNumber });
};

const removeTrackEntry = (index) => {
    albumForm.tracks.splice(index, 1);
};

const formatSongDisplay = (song) => {
    let name = song.title;
    if (song.artists && song.artists.length > 0) {
        name += ` - ${song.artists.map(a => a.artist_name || a.name).join(', ')}`;
    }
    return name;
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

    const validArtists = albumForm.artists
        .filter(a => a.artist_id)
        .map(a => ({ artist_id: Number(a.artist_id), is_primary: a.is_primary }));
        
    if (validArtists.length > 0) {
       payload.artists = validArtists;
    } else {
        formError.value = 'Un álbum requiere de un artista.';
        return;
    }

    if (!isEditing.value && albumForm.showTracksForm) {
        payload.tracks = albumForm.tracks
            .filter(t => t.song_id)
            .map(t => ({ song_id: Number(t.song_id), track_number: Number(t.track_number) }));
    }

    let resp;
    if (isEditing.value) {
      resp = await albumService.update(albumForm.id, payload);
    } else {
      resp = await albumService.create(payload);
    }
    emit('saved', resp.data || resp);
    emit('close');
  } catch (err) {
    console.error(err);
    formError.value = 'Error al guardar el álbum.';
  }
};

const isArtistModalOpen = ref(false);
const pendingArtistIndex = ref(null);

const openCreateArtist = (index) => {
    pendingArtistIndex.value = index;
    isArtistModalOpen.value = true;
};

const handleArtistSaved = (newArtist) => {
    if (pendingArtistIndex.value !== null && albumForm.artists[pendingArtistIndex.value]) {
        albumForm.artists[pendingArtistIndex.value].artist_id = newArtist.id;
        albumForm.artists[pendingArtistIndex.value].artist_name = newArtist.name;
    }
};

const isSongModalOpen = ref(false);
const pendingTrackIndex = ref(null);

const openCreateSong = (index) => {
    pendingTrackIndex.value = index;
    isSongModalOpen.value = true;
};

const handleSongSaved = (newSong) => {
    if (newSong && pendingTrackIndex.value !== null && albumForm.tracks[pendingTrackIndex.value]) {
        albumForm.tracks[pendingTrackIndex.value].song_id = newSong.id;
        albumForm.tracks[pendingTrackIndex.value].song_title = formatSongDisplay(newSong);
    }
};
</script>

<template>
  <Modal :isOpen="isOpen" @close="$emit('close')" :title="isEditing ? 'Editar Álbum' : 'Crear Nuevo Álbum'">
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

      <!-- Multiple Artists Selection -->
      <div class="form-group mb-4">
        <div class="flex justify-between items-center mb-2">
            <label class="mb-0">Artistas</label>
            <button type="button" class="btn btn-secondary btn-sm" @click="addArtistEntry">+ Agregar Artista</button>
        </div>
        
        <div v-for="(artistEntry, index) in albumForm.artists" :key="index" class="artist-row flex gap-2 mb-2 items-start">
            <div class="flex-1 flex gap-2" style="min-width: 0;">
                <div class="flex-1">
                    <SearchSelect 
                        v-model="artistEntry.artist_id"
                        :initialName="artistEntry.artist_name"
                        :searchFn="artistService.search"
                        :formatDisplay="(a) => a.artist_name || a.name"
                        placeholder="Busca un artista..."
                        @select="(item) => artistEntry.artist_name = (item.artist_name || item.name)"
                    />
                </div>
                <button type="button" class="btn btn-secondary icon-btn" title="Crear Nuevo Artista" @click="openCreateArtist(index)">
                    ➕👤
                </button>
            </div>
            <div class="w-1/3">
                <select v-model="artistEntry.is_primary" class="form-input">
                    <option :value="true">Principal</option>
                    <option :value="false">Secundario (Colaborador)</option>
                </select>
            </div>
            <button 
                v-if="albumForm.artists.length > 1" 
                type="button" 
                class="btn btn-danger icon-btn" 
                @click="removeArtistEntry(index)">
                🗑️
            </button>
        </div>
      </div>

      <!-- Optional Initial Tracks Form (Only visible during creation) -->
      <div v-if="!isEditing" class="form-group mb-4 mt-6 border-t pt-4">
          <div class="flex items-center gap-2 mb-2">
            <input type="checkbox" id="addTracks" :checked="albumForm.showTracksForm" @change="toggleTracksForm" />
            <label for="addTracks" class="mb-0 cursor-pointer">Agregar pistas iniciales a este álbum</label>
          </div>

          <div v-if="albumForm.showTracksForm" class="tracks-list mt-4">
               <div v-for="(track, idx) in albumForm.tracks" :key="'trk-'+idx" class="flex gap-2 mb-2 items-start">
                   <div style="width: 80px;">
                       <input type="number" v-model="track.track_number" class="form-input" min="1" placeholder="N°" required />
                   </div>
                   <div class="flex-1 flex gap-2" style="min-width: 0;">
                       <div class="flex-1">
                           <SearchSelect 
                                v-model="track.song_id"
                                :initialName="track.song_title"
                                :searchFn="songService.search"
                                :formatDisplay="formatSongDisplay"
                                placeholder="Busca una canción..."
                                @select="(item) => track.song_title = formatSongDisplay(item)"
                            />
                       </div>
                       <button type="button" class="btn btn-secondary icon-btn" title="Crear Nueva Canción" @click="openCreateSong(idx)">
                           ➕🎵
                       </button>
                   </div>
                   <button 
                        v-if="albumForm.tracks.length > 1" 
                        type="button" 
                        class="btn btn-danger icon-btn" 
                        @click="removeTrackEntry(idx)">
                        🗑️
                    </button>
               </div>
               <button type="button" class="btn btn-secondary btn-sm mt-2" @click="addTrackEntry">+ Añadir otra pista</button>
          </div>
      </div>

      <div class="form-actions">
        <button type="button" class="btn btn-secondary" @click="$emit('close')">Cancelar</button>
        <button type="submit" class="btn btn-primary">Guardar</button>
      </div>
    </form>
  </Modal>

  <ArtistFormModal 
      :isOpen="isArtistModalOpen" 
      @close="isArtistModalOpen = false" 
      @saved="handleArtistSaved" 
  />

  <SongFormModal 
      :isOpen="isSongModalOpen" 
      @close="isSongModalOpen = false" 
      @saved="handleSongSaved" 
  />
</template>

<style scoped>
.error-msg {
  color: var(--danger);
  margin-bottom: 1rem;
  font-size: 0.9rem;
  background: rgba(239, 68, 68, 0.1);
  padding: 0.5rem;
  border-radius: var(--radius-sm);
}

.form-group {
    margin-bottom: 1rem;
}

.form-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    font-size: 0.9rem;
    color: var(--text-secondary);
}

.form-input {
    width: 100%;
    padding: 0.75rem 1rem;
    background: var(--bg-primary);
    border: 1px solid var(--border-color);
    color: var(--text-primary);
    border-radius: var(--radius-md);
    outline: none;
    transition: all var(--transition-fast);
}

.form-input:focus {
    border-color: var(--accent-primary);
    box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
}

.form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 2rem;
}

select.form-input option {
  background: var(--bg-primary);
  color: var(--text-primary);
}

/* Utilities */
.flex { display: flex; }
.justify-between { justify-content: space-between; }
.items-center { align-items: center; }
.items-start { align-items: flex-start; }
.flex-1 { flex: 1; }
.w-1\/3 { width: 33.333333%; }
.gap-2 { gap: 0.5rem; }
.mb-2 { margin-bottom: 0.5rem; }
.mb-4 { margin-bottom: 1rem; }
.mb-0 { margin-bottom: 0 !important; }
.mt-6 { margin-top: 1.5rem; }
.mt-4 { margin-top: 1rem; }
.mt-2 { margin-top: 0.5rem; }
.border-t { border-top: 1px solid var(--border-color); }
.pt-4 { padding-top: 1rem; }
.cursor-pointer { cursor: pointer; }

.btn-sm {
    padding: 0.25rem 0.5rem;
    font-size: 0.8rem;
}

.icon-btn {
    padding: 0.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(239, 68, 68, 0.1) !important;
}

.icon-btn:hover {
    background: rgba(239, 68, 68, 0.2) !important;
    transform: none;
}
</style>
