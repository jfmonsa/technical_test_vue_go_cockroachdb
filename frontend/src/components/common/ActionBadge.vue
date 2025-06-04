<template>
  <span :class="badgeClass">
    {{ action }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface Props {
  action: string;
}

const props = defineProps<Props>();

const badgeClass = computed(() => {
  const baseClass =
    'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium';

  if (props.action.includes('upgraded')) {
    return `${baseClass} bg-green-100 text-green-800`;
  } else if (props.action.includes('downgraded')) {
    return `${baseClass} bg-red-100 text-red-800`;
  } else if (props.action.includes('initiated')) {
    return `${baseClass} bg-blue-100 text-blue-800`;
  } else if (props.action.includes('reiterated')) {
    return `${baseClass} bg-gray-100 text-gray-800`;
  } else if (props.action.includes('lowered')) {
    return `${baseClass} bg-orange-100 text-orange-800`;
  } else {
    return `${baseClass} bg-purple-100 text-purple-800`;
  }
});
</script>
