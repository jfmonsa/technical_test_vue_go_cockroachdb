<template>
  <div
    class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6"
  >
    <div class="flex">
      <svg
        class="h-5 w-5 text-red-400"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
        />
      </svg>
      <div class="ml-3">
        <h3 class="text-sm font-medium text-red-800">Error</h3>
        <p class="text-sm text-red-700 mt-1">{{ error }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ error: string | null}>();
</script>