<template>
  <router-link :to="to" class="btn-secondary flex items-center gap-2">
    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2"
        d="M10 19l-7-7m0 0l7-7m0 0v14"
      />
    </svg>
    {{ text }}
  </router-link>
</template>

<script setup lang="ts">
defineProps<{
  to: string;
  text: string;
}>();
</script>
