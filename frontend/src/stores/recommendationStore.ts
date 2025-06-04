import type { Recommendation } from '@/models/recommendation';
import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

export const useRecommendationStore = defineStore('recommendation', () => {
  // State
  const recommendations = ref<Recommendation[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const limit = ref(10);
  const minimumScore = ref(7);

  // Getters
  const topRecommendations = computed(() => {
    return recommendations.value.sort((a, b) => b.Score - a.Score).slice(0, 5);
  });

  const averageScore = computed(() => {
    if (recommendations.value.length === 0) return 0;
    const sum = recommendations.value.reduce((acc, rec) => acc + rec.Score, 0);
    return sum / recommendations.value.length;
  });

  const scoreDistribution = computed(() => {
    const distribution = {
      excellent: 0, // Score >= 12
      good: 0, // Score >= 9
      fair: 0, // Score >= 7
      poor: 0, // Score < 7
    };

    recommendations.value.forEach((rec) => {
      if (rec.Score >= 12) distribution.excellent++;
      else if (rec.Score >= 9) distribution.good++;
      else if (rec.Score >= 7) distribution.fair++;
      else distribution.poor++;
    });

    return distribution;
  });

  // Actions
  const fetchRecommendations = async () => {
    loading.value = true;
    error.value = null;

    try {
      const response = await fetch(
        `${import.meta.env.VITE_API_BASE_URL}/recommendations?limit=${limit.value}&minimun_score=${minimumScore.value}`,
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data: Recommendation[] = await response.json();
      recommendations.value = data;
    } catch (err) {
      error.value =
        err instanceof Error ? err.message : 'Error fetching recommendations';
      console.error('Error fetching recommendations:', err);
    } finally {
      loading.value = false;
    }
  };

  const setLimit = (newLimit: number) => {
    limit.value = newLimit;
    fetchRecommendations();
  };

  const setMinimumScore = (newScore: number) => {
    minimumScore.value = newScore;
    fetchRecommendations();
  };

  return {
    // State
    recommendations,
    loading,
    error,
    limit,
    minimumScore,

    // Getters
    topRecommendations,
    averageScore,
    scoreDistribution,

    // Actions
    fetchRecommendations,
    setLimit,
    setMinimumScore,
  };
});
