<template>
  <div class="flex items-center gap-2">
    <span class="text-sm font-mono text-gray-700">{{ from }}</span>
    <svg
      class="w-4 h-4 text-gray-400"
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
    <span class="text-sm font-mono" :class="changeClass">{{ to }}</span>
    <span
      v-if="changePercentage !== 0"
      :class="percentageClass"
      class="text-xs font-medium"
    >
      {{ changePercentage > 0 ? "+" : "" }}{{ changePercentage.toFixed(1) }}%
    </span>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";

interface Props {
  from: string;
  to: string;
}

const props = defineProps<Props>();

const fromValue = computed(() => {
  const match = props.from.match(/[\d.]+/);
  return match ? parseFloat(match[0]) : 0;
});

const toValue = computed(() => {
  const match = props.to.match(/[\d.]+/);
  return match ? parseFloat(match[0]) : 0;
});

const changePercentage = computed(() => {
  if (fromValue.value === 0) return 0;
  return ((toValue.value - fromValue.value) / fromValue.value) * 100;
});

const changeClass = computed(() => {
  if (toValue.value > fromValue.value) {
    return "text-green-600";
  } else if (toValue.value < fromValue.value) {
    return "text-red-600";
  } else {
    return "text-gray-700";
  }
});

const percentageClass = computed(() => {
  if (changePercentage.value > 0) {
    return "text-green-600";
  } else if (changePercentage.value < 0) {
    return "text-red-600";
  } else {
    return "text-gray-600";
  }
});
</script>
