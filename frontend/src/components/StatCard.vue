<template>
    <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <div class="flex items-center">
        <slot name="icon" />
        <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">{{ label }}</p>
            <p class="text-2xl font-semibold text-gray-900">
            {{ value.toLocaleString() }}
            </p>
        </div>
        </div>
    </div>
</template>

<script setup lang="ts">
defineProps<{
    label: string;
  value: string | number;
}>();
</script>