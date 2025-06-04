<template>
  <div class="mb-6 rounded-lg border border-gray-200 bg-white p-6 shadow-sm">
    <div class="flex flex-col items-center justify-between gap-4 sm:flex-row">
      <div class="flex items-center gap-4">
        <!-- Minimum Score Filter -->
        <div class="flex items-center gap-2">
          <label class="text-sm font-medium text-gray-700">Score mínimo:</label>
          <select
            :value="minimumScore"
            @change="handleScoreChange"
            class="rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-2 focus:ring-blue-500"
          >
            <option value="5">5+</option>
            <option value="7">7+</option>
            <option value="9">9+</option>
            <option value="12">12+</option>
          </select>
        </div>

        <!-- Limit Filter -->
        <div class="flex items-center gap-2">
          <label class="text-sm font-medium text-gray-700">Mostrar:</label>
          <select
            :value="limit"
            @change="handleLimitChange"
            class="rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-2 focus:ring-blue-500"
          >
            <option value="5">5 results</option>
            <option value="10">10 results</option>
            <option value="20">20 results</option>
            <option value="50">50 results</option>
          </select>
        </div>
      </div>

      <!-- Refresh Button -->
      <button
        @click="refreshRecommendations"
        :disabled="loading"
        class="btn-primary flex items-center gap-2 disabled:cursor-not-allowed disabled:opacity-50"
      >
        <svg
          class="h-4 w-4"
          :class="{ 'animate-spin': loading }"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
          />
        </svg>
        {{ loading ? 'Updating...' : 'Updated' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  minimumScore: number;
  limit: number;
  loading: boolean;
}>();

const emit = defineEmits(['update:minimumScore', 'update:limit', 'refresh']);

import { useRecommendationStore } from '@/stores/recommendationStore';
const { setMinimumScore, setLimit } = useRecommendationStore();

const handleScoreChange = (event: Event) => {
  const target = event.target as HTMLSelectElement;
  emit('update:minimumScore', setMinimumScore(parseInt(target.value)));
};
const handleLimitChange = (event: Event) => {
  const target = event.target as HTMLSelectElement;
  emit('update:limit', setLimit(parseInt(target.value)));
};
const refreshRecommendations = () => {
  emit('refresh');
};
</script>
