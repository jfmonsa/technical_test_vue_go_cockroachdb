<template>
  <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
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
import StatCard from "./StatCard.vue";
import StatsIcon from "./Icons/StatsIcon.vue";
import SearchGreenIcon from "./Icons/SearchGreenIcon.vue";
import PageIcon from "./Icons/PageIcon.vue";

defineProps<{
  total: number;
  searchResults: number;
  currentPage: number;
  totalPages: number;
}>();
</script>