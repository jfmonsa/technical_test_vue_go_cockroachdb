<template>
  <div class="mb-6 grid grid-cols-1 gap-6 md:grid-cols-3">
    <StatCard label="Total Stocks" :value="total">
      <template #icon><StatsIcon /></template>
    </StatCard>
    <StatCard label="Search Results" :value="searchResults">
      <template #icon><SearchGreenIcon /></template>
    </StatCard>
    <StatCard label="Current Page" :value="currentPage + ' / ' + totalPages">
      <template #icon><PageIcon /></template>
    </StatCard>
  </div>
</template>

<script setup lang="ts">
import StatCard from './StatCard.vue';
import StatsIcon from './Icons/StatsIcon.vue';
import SearchGreenIcon from './Icons/SearchGreenIcon.vue';
import PageIcon from './Icons/PageIcon.vue';

defineProps<{
  total: number;
  searchResults: number;
  currentPage: number;
  totalPages: number;
}>();
</script>
