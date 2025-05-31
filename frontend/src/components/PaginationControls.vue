<template>
  <div class="bg-gray-50 px-6 py-3 border-t border-gray-200">
    <div class="flex items-center justify-between">
      <div class="text-sm text-gray-700">
        Showing {{ start }} - {{ end }} of {{ total }} results
      </div>
      <div class="flex items-center gap-2">
        <button
          @click="$emit('page', currentPage - 1)"
          :disabled="currentPage <= 1"
          class="btn-secondary disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Previous
        </button>

        <div class="flex items-center gap-1">
          <button
            v-for="page in visiblePages"
            :key="page"
            @click="$emit('page', page)"
            :class="[
              'px-3 py-1 rounded text-sm font-medium transition-colors',
              page === currentPage
                ? 'bg-blue-600 text-white'
                : 'bg-white text-gray-700 hover:bg-gray-100 border border-gray-300',
            ]"
          >
            {{ page }}
          </button>
        </div>

        <button
          @click="$emit('page', currentPage + 1)"
          :disabled="currentPage >= totalPages"
          class="btn-secondary disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Next
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  currentPage: number;
  totalPages: number;
  limit: number;
  total: number;
}>();

const start = computed(() => (props.currentPage - 1) * props.limit + 1);
const end = computed(() => Math.min(props.currentPage * props.limit, props.total));

const visiblePages = computed(() => {
  const pages = [];
  const start = Math.max(1, props.currentPage - 2);
  const end = Math.min(props.totalPages, props.currentPage + 2);

  for (let i = start; i <= end; i++) {
    pages.push(i);
  }
  return pages;
});
</script>