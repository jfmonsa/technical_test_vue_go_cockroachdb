<template>
  <svg
    class="w-4 h-4 text-gray-400"
    :class="{ 'text-blue-600': isActive }"
    fill="none"
    stroke="currentColor"
    viewBox="0 0 24 24"
  >
    <path
      v-if="!isActive"
      stroke-linecap="round"
      stroke-linejoin="round"
      stroke-width="2"
      d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"
    />
    <path
      v-else-if="direction === 'asc'"
      stroke-linecap="round"
      stroke-linejoin="round"
      stroke-width="2"
      d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12"
    />
    <path
      v-else
      stroke-linecap="round"
      stroke-linejoin="round"
      stroke-width="2"
      d="M3 4h13M3 8h9m-9 4h9m5-4v12m0 0l-4-4m4 4l4-4"
    />
  </svg>
</template>

<script setup lang="ts">
import { computed } from "vue";

interface Props {
  field: string;
  currentField: string;
  direction: "asc" | "desc";
}

const props = defineProps<Props>();

const isActive = computed(() => props.field === props.currentField);
</script>
