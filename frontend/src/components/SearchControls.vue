<template>
  <div class="mb-6 rounded-lg border border-gray-200 bg-white p-6 shadow-sm">
    <div class="flex flex-col items-center justify-between gap-4 sm:flex-row">
      <div class="relative max-w-md flex-1">
        <SearchIcon />
        <input
          v-model="localSearchQuery"
          @input="handleSearch"
          type="text"
          placeholder="Search by ticker, company, brokerage..."
          class="block w-full rounded-lg border border-gray-300 py-2 pr-3 pl-10 focus:border-blue-500 focus:ring-2 focus:ring-blue-500"
        />
      </div>
      <RefreshButton :loading="loading" :on-click="refreshData" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';

const props = defineProps(['searchQuery', 'loading']);
const emit = defineEmits(['update:searchQuery', 'refresh']);
import SearchIcon from './Icons/SearchIcon.vue';
import RefreshButton from './RefreshDataBtn.vue';

const localSearchQuery = ref(props.searchQuery);

watch(
  () => props.searchQuery,
  (newVal) => {
    localSearchQuery.value = newVal;
  },
);

const handleSearch = (event: Event) => {
  const target = event.target as HTMLInputElement;
  emit('update:searchQuery', target.value);
};
const refreshData = () => emit('refresh');
</script>
