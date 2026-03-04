<script setup>
import { ref, reactive, watch } from 'vue';
import { songService } from '../../services/song.service';
import { artistService } from '../../services/artist.service';
import { useToast } from '../../composables/useToast';
import Modal from '../common/Modal.vue';
import Icon from '../common/Icon.vue';

const toast = useToast();
import SearchSelect from '../common/SearchSelect.vue';
import ArtistFormModal from '../artists/ArtistFormModal.vue';

const props = defineProps({
  isOpen: { type: Boolean, required: true },
  song: { type: Object, default: null }, // Null to create, Object to edit
  initialArtist: { type: Object, default: null } // { id, name } to pre-fill artist
});

const emit = defineEmits(['close', 'saved']);

const formError = ref(null);
const isEditing = ref(false);

const songForm = reactive({
  id: null,
  title: '',
  duration: 0,
  artists: []
});

const resetForm = () => {
  formError.value = null;
  if (props.song) {
    // Editing existing song
    isEditing.value = true;
    const initialArtists = props.song.artists && props.song.artists.length > 0 
      ? props.song.artists.map(a => ({ artist_id: a.id || a.artist_id, artist_name: a.name || a.artist_name, role: a.role }))
      : [{ artist_id: '', artist_name: '', role: 'main' }];

    Object.assign(songForm, {
      id: props.song.id,
      title: props.song.title,
      duration: props.song.duration,
      artists: [...initialArtists]
    });
  } else {
    // Creating new song
    isEditing.value = false;
    let initialArtists = [];
    if (props.initialArtist && props.initialArtist.id) {
        initialArtists.push({
            artist_id: props.initialArtist.id,
            artist_name: props.initialArtist.name || '',
            role: 'main'
        });
    } else {
        initialArtists.push({ artist_id: '', artist_name: '', role: 'main' });
    }

    Object.assign(songForm, {
      id: null,
      title: '',
      duration: 0,
      artists: initialArtists
    });
  }
};

watch(() => props.isOpen, (newVal) => {
    if (newVal) {
        resetForm();
    }
});

const removeArtist = (index) => {
  songForm.artists.splice(index, 1);
};

const addArtist = () => {
  songForm.artists.push({ artist_id: '', artist_name: '', role: 'main' });
};

const saveSong = async () => {
  try {
    formError.value = null;
    let payload = {
      title: songForm.title,
      duration: Number(songForm.duration)
    };
    
    const validArtists = songForm.artists
        .filter(a => a.artist_id)
        .map(a => ({ artist_id: Number(a.artist_id), role: a.role }));
        
    if (validArtists.length > 0) {
       payload.artists = validArtists;
    } else {
       // Optional: Enforce at least one artist if you want, or leave it as optional as it was
    }

    let resp;
    if (isEditing.value) {
      resp = await songService.update(songForm.id, payload);
      toast.success('Canción actualizada exitosamente');
    } else {
      resp = await songService.create(payload);
      toast.success('Canción creada exitosamente');
    }
    
    emit('saved', resp.data || resp);
    emit('close');
  } catch (err) {
    toast.handleApiError(err, 'Error al guardar la canción');
    formError.value = 'Revisa los datos ingresados.';
  }
};

const isArtistModalOpen = ref(false);
const pendingArtistIndex = ref(null);

const openCreateArtist = (index) => {
    pendingArtistIndex.value = index;
    isArtistModalOpen.value = true;
};

const handleArtistSaved = (newArtist) => {
    if (pendingArtistIndex.value !== null && songForm.artists[pendingArtistIndex.value]) {
        songForm.artists[pendingArtistIndex.value].artist_id = newArtist.id;
        songForm.artists[pendingArtistIndex.value].artist_name = newArtist.name;
    }
};
</script>

<template>
  <Modal :isOpen="isOpen" @close="$emit('close')" :title="isEditing ? 'Editar Canción' : 'Nueva Canción'">
    <form @submit.prevent="saveSong">
      <div v-if="formError" class="error-msg">{{ formError }}</div>
      
      <div class="form-group">
        <label><Icon name="music" class="label-icon" /> Título</label>
        <input type="text" v-model="songForm.title" class="form-input" required />
      </div>
      
      <div class="form-group">
        <label><Icon name="clock" class="label-icon" /> Duración (Segundos)</label>
        <input type="number" v-model="songForm.duration" class="form-input" required min="1" />
      </div>

      <div class="form-group mb-4">
        <div class="flex justify-between items-center mb-2">
          <label class="mb-0"><Icon name="users" class="label-icon" /> Artistas</label>
          <button type="button" class="btn btn-secondary btn-sm" @click="addArtist">
            <Icon name="plus" /> Agregar Artista
          </button>
        </div>
        
        <div v-for="(artistEntry, index) in songForm.artists" :key="index" class="artist-row flex gap-2 mb-2 items-start">
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
                    <Icon name="user-plus" />
                </button>
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
                <Icon name="trash" />
            </button>
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
