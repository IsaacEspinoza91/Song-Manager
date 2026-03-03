<script setup>
import { computed } from 'vue';

const props = defineProps({
  currentPage: {
    type: Number,
    required: true
  },
  totalPages: {
    type: Number,
    required: true
  },
  totalItems: {
    type: Number,
    required: true
  }
});

const emit = defineEmits(['page-change']);

const visiblePages = computed(() => {
    const pages = [];
    const maxVisible = 5;
    
    if (props.totalPages <= maxVisible) {
        for (let i = 1; i <= props.totalPages; i++) pages.push(i);
        return pages;
    }
    
    // Logic for pages with ellipses
    if (props.currentPage <= 3) {
        return [1, 2, 3, 4, '...', props.totalPages];
    } else if (props.currentPage >= props.totalPages - 2) {
        return [1, '...', props.totalPages - 3, props.totalPages - 2, props.totalPages - 1, props.totalPages];
    } else {
        return [1, '...', props.currentPage - 1, props.currentPage, props.currentPage + 1, '...', props.totalPages];
    }
});

const changePage = (page) => {
    if (page === '...' || page === props.currentPage || page < 1 || page > props.totalPages) return;
    emit('page-change', page);
};

const goToNext = () => changePage(props.currentPage + 1);
const goToPrev = () => changePage(props.currentPage - 1);
</script>

<template>
  <div class="pagination-wrapper" v-if="totalPages > 0">
      <div class="pagination-controls glass-panel">
          <button 
              class="page-btn nav-btn" 
              :disabled="currentPage === 1" 
              @click="goToPrev"
              aria-label="Página anterior"
          >
              &laquo;
          </button>
          
          <button 
              v-for="(page, index) in visiblePages" 
              :key="`page-${index}-${page}`"
              class="page-btn"
              :class="{ 'active': page === currentPage, 'dots': page === '...' }"
              :disabled="page === '...'"
              @click="changePage(page)"
          >
              {{ page }}
          </button>

          <button 
              class="page-btn nav-btn" 
              :disabled="currentPage === totalPages" 
              @click="goToNext"
              aria-label="Página siguiente"
          >
              &raquo;
          </button>
      </div>
      <div class="page-info">
          ({{ totalItems }} elementos)
      </div>
  </div>
</template>

<style scoped>
.pagination-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  margin-top: 3rem;
  padding-bottom: 2rem;
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.5rem;
  border-radius: var(--radius-full);
}

.page-btn {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: none;
  background: transparent;
  color: var(--text-primary);
  font-weight: 600;
  font-size: 0.95rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition-fast);
}

.page-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.1);
}

.page-btn.active {
  background: var(--accent-primary);
  color: var(--bg-primary);
  box-shadow: 0 0 15px rgba(30, 215, 96, 0.4);
}

.page-btn.dots {
  cursor: default;
  color: var(--text-muted);
  background: transparent;
}

.page-btn.nav-btn {
  font-size: 1.2rem;
  color: var(--accent-primary);
}

.page-btn.nav-btn:hover:not(:disabled) {
  background: rgba(30, 215, 96, 0.1);
}

.page-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.page-info {
  color: var(--text-muted);
  font-size: 0.85rem;
}
</style>
