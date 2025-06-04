<template>
  <div class="flex items-center gap-2">
    <span :class="fromClass">{{ from }}</span>
    <svg
      class="h-4 w-4 text-gray-400"
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2"
        d="M9 5l7 7-7 7"
      />
    </svg>
    <span :class="toClass">{{ to }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface Props {
  from: string;
  to: string;
}

const props = defineProps<Props>();

const getRatingClass = (rating: string) => {
  const baseClass = 'px-2 py-1 rounded text-xs font-medium';

  if (
    rating.toLowerCase().includes('buy') ||
    rating.toLowerCase().includes('outperform')
  ) {
    return `${baseClass} bg-green-100 text-green-800`;
  } else if (
    rating.toLowerCase().includes('sell') ||
    rating.toLowerCase().includes('underperform')
  ) {
    return `${baseClass} bg-red-100 text-red-800`;
  } else if (
    rating.toLowerCase().includes('hold') ||
    rating.toLowerCase().includes('neutral')
  ) {
    return `${baseClass} bg-yellow-100 text-yellow-800`;
  } else {
    return `${baseClass} bg-gray-100 text-gray-800`;
  }
};

const fromClass = computed(() => getRatingClass(props.from));
const toClass = computed(() => getRatingClass(props.to));
</script>
