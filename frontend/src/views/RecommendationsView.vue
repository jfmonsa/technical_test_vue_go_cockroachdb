<template>
  <div class="min-h-screen bg-gray-50 py-8">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <!-- Header -->
      <HeaderSection />

      <!-- Controls -->
      <RecommendationsControls
        :minimumScore="minimumScore"
        :limit="limit"
        :loading="loading"
        @update:minimumScore="(val: number) => (minimumScore = val)"
        @update:limit="(val: number) => (limit = val)"
        @refresh="refreshRecommendations"
      />

      <!-- Error Message -->
      <ErrorAlert v-if="error" :error="error" />

      <!-- Stats Dashboard -->
      <StatsSection
        :numberOfRecommendations="recommendations.length"
        :averageScore="averageScore"
        :scoreDistribution="scoreDistribution"
      />

      <!-- Top 3 Recommendations Highlight -->
      <TopRecommendationsSection :topRecommendations="topRecommendations" />

      <!-- Recommendations Table -->
      <RecommendationsTable
        :recommendations="recommendations"
        :loading="loading"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useRecommendationStore } from '@/stores/recommendationStore';
import { storeToRefs } from 'pinia';
import ErrorAlert from '@/components/common/ErrorAlert.vue';
import RecommendationsControls from '@/components/recommendationsView/RecommendationsControls.vue';
import HeaderSection from '@/components/recommendationsView/HeaderSection.vue';
import StatsSection from '@/components/recommendationsView/StatsSection.vue';
import TopRecommendationsSection from '@/components/recommendationsView/TopRecommendationsSection.vue';
import RecommendationsTable from '@/components/recommendationsView/RecommendationsTable.vue';

const recommendationStore = useRecommendationStore();

const {
  recommendations,
  loading,
  error,
  limit,
  minimumScore,
  topRecommendations,
  averageScore,
  scoreDistribution,
} = storeToRefs(recommendationStore);

// Destructure actions
const { fetchRecommendations } = recommendationStore;

const refreshRecommendations = () => {
  fetchRecommendations();
};

// Lifecycle
onMounted(() => {
  fetchRecommendations();
});
</script>
