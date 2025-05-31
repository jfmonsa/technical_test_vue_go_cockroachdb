<template>
  <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 mb-6">
    <div class="flex flex-col sm:flex-row gap-4 items-center justify-between">
      <div class="relative flex-1 max-w-md">
        <SearchIcon />
        <input
          v-model="localSearchQuery"
          @input="handleSearch"
          type="text"
          placeholder="Buscar por ticker, compañía, brokerage..."
          class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
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
import SearchIcon from "./Icons/SearchIcon.vue";
import RefreshButton from "./RefreshDataBtn.vue";

const localSearchQuery = ref(props.searchQuery);

watch(() => props.searchQuery, (newVal) => {
  localSearchQuery.value = newVal;
});

const handleSearch = (event: Event) => {
  const target = event.target as HTMLInputElement;
  emit('update:searchQuery', target.value);
};
const refreshData = () => emit('refresh');
</script>