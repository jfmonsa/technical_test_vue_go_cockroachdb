<template>
  <div
    class="overflow-hidden rounded-lg border border-gray-200 bg-white shadow-sm"
  >
    <div class="border-b border-gray-200 px-6 py-4">
      <h2 class="text-lg font-semibold text-gray-900">All Recommendations</h2>
    </div>

    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="table-header">Ranking</th>
            <th class="table-header">Ticker</th>
            <th class="table-header">Company</th>
            <th class="table-header">Score</th>
            <th class="table-header">Brokerage</th>
            <th class="table-header">Action</th>
            <th class="table-header">Rating</th>
            <th class="table-header">Target</th>
            <th class="table-header">Date</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 bg-white">
          <tr v-if="loading" class="animate-pulse">
            <td colspan="9" class="table-cell py-8 text-center">
              <div class="flex items-center justify-center">
                <svg
                  class="h-6 w-6 animate-spin text-blue-600"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    class="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    stroke-width="4"
                  ></circle>
                  <path
                    class="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  ></path>
                </svg>
                <span class="ml-2">Loading Recommendations...</span>
              </div>
            </td>
          </tr>
          <tr v-else-if="recommendations.length === 0" class="text-center">
            <td colspan="9" class="table-cell py-8 text-gray-500">
              No recommendations found with the selected minimum score
            </td>
          </tr>
          <tr
            v-else
            v-for="(rec, index) in recommendations"
            :key="rec.ticker"
            class="transition-colors hover:bg-gray-50"
          >
            <td class="table-cell">
              <div class="flex items-center justify-center">
                <span
                  class="flex h-8 w-8 items-center justify-center rounded-full bg-blue-100 text-sm font-semibold text-blue-800"
                >
                  {{ index + 1 }}
                </span>
              </div>
            </td>
            <td class="table-cell">
              <span class="font-mono font-semibold text-blue-600">{{
                rec.ticker
              }}</span>
            </td>
            <td class="table-cell">
              <div class="font-medium text-gray-900">{{ rec.company }}</div>
            </td>
            <td class="table-cell">
              <ScoreBadge :score="rec.recommendation_score" />
            </td>
            <td class="table-cell">
              <span class="text-gray-700">{{ rec.brokerage }}</span>
            </td>
            <td class="table-cell">
              <ActionBadge :action="rec.action" />
            </td>
            <td class="table-cell">
              <RatingChange :from="rec.rating_from" :to="rec.rating_to" />
            </td>
            <td class="table-cell">
              <TargetChange
                :from="rec.target_from.toString()"
                :to="rec.target_to.toString()"
              />
            </td>
            <td class="table-cell">
              <time class="text-sm text-gray-600">{{
                formatDate(rec.time)
              }}</time>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import ScoreBadge from '@/components/recommendationsView/ScoreBadge.vue';
import ActionBadge from '@/components/common/ActionBadge.vue';
import RatingChange from '@/components/common/RatingChange.vue';
import TargetChange from '@/components/common/TargetChange.vue';
import { formatDate } from '@/utils/dates';
import type { Recommendation } from '@/models/recommendation';

defineProps<{
  recommendations: Recommendation[];
  loading: boolean;
}>();
</script>
