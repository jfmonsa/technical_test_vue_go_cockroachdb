<template>
  <button
    @click="onClick"
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
    {{ loading ? 'Loading...' : 'Update' }}
  </button>
</template>

<script setup lang="ts">
defineProps<{
  loading: boolean;
  onClick: () => void;
}>();
</script>
