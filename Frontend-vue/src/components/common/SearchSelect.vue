<script setup>
import { ref, watch, onMounted, onBeforeUnmount } from 'vue';

const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: null
  },
  initialName: {
    type: String,
    default: ''
  },
  searchFn: {
    type: Function,
    required: true
  },
  formatDisplay: {
    type: Function,
    default: (item) => item.name || item.title || item.artist_name
  },
  placeholder: {
    type: String,
    default: 'Buscar...'
  }
});

const emit = defineEmits(['update:modelValue', 'select']);

const query = ref(props.initialName);
const results = ref([]);
const isOpen = ref(false);
const isLoading = ref(false);
const wrapperRef = ref(null);
let debounceTimeout = null;

watch(() => props.initialName, (newVal) => {
    if (!isOpen.value) {
        query.value = newVal;
    }
});

const onInput = (e) => {
    query.value = e.target.value;
    emit('update:modelValue', null); // Clear selection when typing
    isOpen.value = true;
    
    if (debounceTimeout) clearTimeout(debounceTimeout);
    
    if (!query.value || query.value.trim() === '') {
        results.value = [];
        return;
    }

    debounceTimeout = setTimeout(async () => {
        isLoading.value = true;
        try {
            const resp = await props.searchFn(query.value);
            results.value = resp.data || resp || [];
        } catch (err) {
            console.error('Search error', err);
            results.value = [];
        } finally {
            isLoading.value = false;
        }
    }, 300);
};

const selectItem = (item) => {
    query.value = props.formatDisplay(item);
    isOpen.value = false;
    emit('update:modelValue', item.id);
    emit('select', item);
};

const handleClickOutside = (e) => {
    if (wrapperRef.value && !wrapperRef.value.contains(e.target)) {
        isOpen.value = false;
    }
};

onMounted(() => {
    document.addEventListener('click', handleClickOutside);
});

onBeforeUnmount(() => {
    document.removeEventListener('click', handleClickOutside);
});
</script>

<template>
  <div class="search-select-wrapper" ref="wrapperRef">
    <input 
      type="text" 
      class="form-input" 
      :placeholder="placeholder"
      :value="query"
      @input="onInput"
      @focus="() => { if(query) onInput({target: {value: query}}) }"
    />
    
    <div class="dropdown-menu" v-if="isOpen && (results.length > 0 || isLoading || query)">
        <div v-if="isLoading" class="dropdown-item text-muted">Buscando...</div>
        <div v-else-if="results.length === 0 && query && !isLoading" class="dropdown-item text-muted">No se encontraron resultados</div>
        <div 
            v-for="item in results" 
            :key="item.id" 
            class="dropdown-item selectable"
            @click="selectItem(item)"
        >
            {{ formatDisplay(item) }}
        </div>
    </div>
  </div>
</template>

<style scoped>
.search-select-wrapper {
    position: relative;
    width: 100%;
}

.dropdown-menu {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background-color: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-sm);
    margin-top: 4px;
    max-height: 250px;
    overflow-y: auto;
    z-index: 50;
    box-shadow: 0 4px 12px rgba(0,0,0,0.5);
}

.dropdown-item {
    padding: 0.75rem 1rem;
    font-size: 0.95rem;
    border-bottom: 1px solid rgba(255,255,255,0.05);
}

.dropdown-item:last-child {
    border-bottom: none;
}

.dropdown-item.selectable {
    cursor: pointer;
    transition: background-color var(--transition-fast);
}

.dropdown-item.selectable:hover {
    background-color: rgba(255,255,255,0.1);
}

.text-muted {
    color: var(--text-muted);
}
</style>
