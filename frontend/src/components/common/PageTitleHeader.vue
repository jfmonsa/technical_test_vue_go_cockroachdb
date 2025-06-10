<template>
  <div>
    <h1 class="mb-2 text-3xl font-bold text-gray-900">
      {{ title }}
    </h1>
    <p class="text-gray-600">
      {{ subtitle }}
    </p>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string;
  subtitle: string;
}>();
</script>
