<script setup>
import { defineProps, defineEmits } from 'vue';
import Modal from './Modal.vue';

const props = defineProps({
  isOpen: {
    type: Boolean,
    required: true
  },
  itemName: {
    type: String,
    default: 'este elemento'
  }
});

const emit = defineEmits(['close', 'confirm']);
</script>

<template>
  <Modal :isOpen="isOpen" @close="$emit('close')" title="Confirmar Eliminación">
    <div class="confirm-content">
      <div class="warning-icon">⚠️</div>
      <p>¿Estás seguro de que deseas eliminar <strong>{{ itemName }}</strong>?</p>
      <p class="text-muted">Esta acción no se puede deshacer.</p>
    </div>
    
    <div class="form-actions mt-4">
      <button class="btn btn-secondary" @click="$emit('close')">Cancelar</button>
      <button class="btn btn-solid-danger" @click="$emit('confirm')">Eliminar</button>
    </div>
  </Modal>
</template>

<style scoped>
.confirm-content {
  text-align: center;
  padding: 1rem 0;
}

.warning-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.text-muted {
  color: var(--text-muted);
  font-size: 0.9rem;
  margin-top: 0.5rem;
}

.mt-4 {
  margin-top: 1.5rem;
}

.btn-solid-danger {
  background-color: var(--danger);
  color: white;
  border-color: var(--danger);
}

.btn-solid-danger:hover {
  background-color: #b91c1c; /* darker red */
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.4);
}
</style>
