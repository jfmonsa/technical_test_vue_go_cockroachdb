<script setup lang="ts">
import NavigationBtn from '../common/NavigationBtn.vue';
import PageTitleHeader from '../common/PageTitleHeader.vue';
</script>

<!-- TODO: hacer un componente reutilizable para el layout principal de todas las cards -->
<template>
  <div class="mb-8 flex items-center justify-between">
    <PageTitleHeader
      title="Investment Recommendations"
      subtitle="The best investment opportunities based on expert analysis"
    />
    <NavigationBtn
      to="/"
      text="View All Stocks"
      class="btn-secondary flex items-center gap-2"
    />
  </div>
</template>
