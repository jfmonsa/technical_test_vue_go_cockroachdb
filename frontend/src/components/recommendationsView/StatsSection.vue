<template>
  <div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-4">
    <StatCard label="Total Recommendations" :value="numberOfRecommendations">
      <template #icon><StatsIcon /></template>
    </StatCard>

    <StatCard label="Average Score" :value="averageScore.toFixed(1)">
      <template #icon>
        <TrendUpIcon />
      </template>
    </StatCard>

    <StatCard label="Excellent (12+)" :value="scoreDistribution.excellent">
      <template #icon>
        <StarIcon />
      </template>
    </StatCard>

    <StatCard label="Good (9+)" :value="scoreDistribution.good">
      <template #icon>
        <StarIcon />
      </template>
    </StatCard>
  </div>
</template>

<script setup lang="ts">
import StatCard from '../common/StatCard.vue';
import StatsIcon from '../Icons/ChartBarIcon.vue';
import TrendUpIcon from '../Icons/TrendUpIcon.vue';
import StarIcon from '../Icons/StarIcon.vue';

defineProps<{
  numberOfRecommendations: number;
  averageScore: number;
  scoreDistribution: { excellent: number; good: number };
}>();
</script>
