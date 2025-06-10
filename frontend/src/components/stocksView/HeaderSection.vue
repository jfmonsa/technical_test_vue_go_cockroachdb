<script setup lang="ts">
import NavigationBtn from '../common/NavigationBtn.vue';
import PageTitleHeader from '../common/PageTitleHeader.vue';
</script>

<template>
  <div class="mb-8 flex items-center justify-between">
    <PageTitleHeader
      title="Stock Analyzer"
      subtitle="Real-time stock analysis and recommendations"
    />
    <NavigationBtn
      to="/recommendations"
      text="View Recommendations"
      class="btn-secondary flex items-center gap-2"
    />
  </div>
</template>
