<script setup>
import { ref, reactive, watch } from 'vue';
import { artistService } from '../../services/artist.service';
import { useToast } from '../../composables/useToast';
import Modal from '../common/Modal.vue';

const toast = useToast();

const props = defineProps({
  isOpen: { type: Boolean, required: true },
  artist: { type: Object, default: null } // Null to create, Object to edit
});

const emit = defineEmits(['close', 'saved']);

const formError = ref(null);
const isEditing = ref(false);

const artistForm = reactive({
  id: null,
  name: '',
  genre: '',
  country: '',
  bio: '',
  image_url: ''
});

const resetForm = () => {
    formError.value = null;
    if (props.artist) {
        isEditing.value = true;
        Object.assign(artistForm, props.artist);
    } else {
        isEditing.value = false;
        Object.assign(artistForm, { id: null, name: '', genre: '', country: '', bio: '', image_url: '' });
    }
};

watch(() => props.isOpen, (newVal) => {
    if (newVal) {
        resetForm();
    }
});

const saveArtist = async () => {
  try {
    formError.value = null;
    let resp;
    if (isEditing.value) {
      resp = await artistService.update(artistForm.id, artistForm);
      toast.success('Artista actualizado exitosamente');
    } else {
      resp = await artistService.create(artistForm);
      toast.success('Artista creado exitosamente');
    }
    emit('saved', resp.data || resp);
    emit('close');
  } catch (err) {
    toast.handleApiError(err, 'Error al guardar el artista');
    formError.value = 'Revisa los datos ingresados.';
  }
};
</script>

<template>
  <Modal :isOpen="isOpen" @close="$emit('close')" :title="isEditing ? 'Editar Artista' : 'Nuevo Artista'">
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
        <button type="button" class="btn btn-secondary" @click="$emit('close')">Cancelar</button>
        <button type="submit" class="btn btn-primary">Guardar</button>
      </div>
    </form>
  </Modal>
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

.form-input, .form-textarea {
    width: 100%;
    padding: 0.75rem 1rem;
    background: var(--bg-primary);
    border: 1px solid var(--border-color);
    color: var(--text-primary);
    border-radius: var(--radius-md);
    outline: none;
    transition: all var(--transition-fast);
}

.form-textarea {
    resize: vertical;
}

.form-input:focus, .form-textarea:focus {
    border-color: var(--accent-primary);
    box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
}

.form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 2rem;
}
</style>
