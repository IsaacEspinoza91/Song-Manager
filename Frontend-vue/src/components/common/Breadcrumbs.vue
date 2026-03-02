<script setup>
import { computed } from 'vue';
import { useRoute } from 'vue-router';

const props = defineProps({
  items: {
    type: Array,
    required: true,
    // [ { label: 'Home', to: '/' }, { label: 'Artists', to: '/artists' }, { label: 'Gustavo Cerati' } ]
  }
});

const route = useRoute();
</script>

<template>
  <nav class="breadcrumbs" aria-label="Breadcrumb">
    <ol>
      <li v-for="(item, index) in items" :key="index" :class="{ 'active': !item.to }">
        <router-link v-if="item.to" :to="item.to" class="link">
          {{ item.label }}
        </router-link>
        <span v-else class="current">
          {{ item.label }}
        </span>
        <span v-if="index < items.length - 1" class="separator">/</span>
      </li>
    </ol>
  </nav>
</template>

<style scoped>
.breadcrumbs {
  margin-bottom: 1.5rem;
  font-size: 0.9rem;
}

.breadcrumbs ol {
  display: flex;
  flex-wrap: wrap;
  list-style: none;
  padding: 0;
  margin: 0;
}

.breadcrumbs li {
  display: flex;
  align-items: center;
}

.link {
  color: var(--text-muted);
  text-decoration: none;
  transition: color var(--transition-fast);
}

.link:hover {
  color: var(--accent-primary);
}

.current {
  color: var(--text-primary);
  font-weight: 500;
}

.separator {
  margin: 0 0.5rem;
  color: var(--text-muted);
  opacity: 0.5;
}
</style>
