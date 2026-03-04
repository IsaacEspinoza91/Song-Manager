<script setup>
import { useToast } from '../../composables/useToast';

const { toasts, removeToast } = useToast();
</script>

<template>
  <div class="toast-container">
    <TransitionGroup name="toast">
      <div 
        v-for="toast in toasts" 
        :key="toast.id" 
        class="toast"
        :class="`toast-${toast.type}`"
      >
        <div class="toast-icon">
          <span v-if="toast.type === 'success'">✅</span>
          <span v-else-if="toast.type === 'error'">❌</span>
          <span v-else>ℹ️</span>
        </div>
        <div class="toast-content">
          {{ toast.message }}
        </div>
        <button class="toast-close" @click="removeToast(toast.id)">✖</button>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.toast-container {
  position: fixed;
  bottom: 2rem;
  right: 2rem;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  pointer-events: none; /* Let clicks pass through empty areas */
}

.toast {
  pointer-events: auto;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  border-radius: var(--radius-md);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  background: var(--bg-secondary);
  color: var(--text-primary);
  min-width: 300px;
  max-width: 400px;
  border-left: 4px solid transparent;
}

.toast-success {
  border-left-color: var(--accent-primary);
}

.toast-error {
  border-left-color: var(--danger);
  background: rgba(239, 68, 68, 0.1); /* Slight red tint backing */
}

.toast-icon {
  font-size: 1.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.toast-content {
  flex: 1;
  font-size: 0.9rem;
  line-height: 1.4;
}

.toast-close {
  background: none;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  padding: 0;
  font-size: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0.7;
  transition: opacity 0.2s ease;
}

.toast-close:hover {
  opacity: 1;
  color: var(--text-primary);
}

/* Animations */
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(50px);
}

.toast-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>
