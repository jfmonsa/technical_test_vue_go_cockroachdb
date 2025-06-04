<template>
  <div v-if="topRecommendations.length > 0" class="mb-8">
    <h2 class="mb-4 text-xl font-semibold text-gray-900">
      🏆 Top 3 Recommendations
    </h2>
    <div class="grid grid-cols-1 gap-6 md:grid-cols-3">
      <div
        v-for="(rec, index) in topRecommendations.slice(0, 3)"
        :key="rec.ticker"
        class="rounded-lg border border-blue-200 bg-gradient-to-br from-blue-50 to-indigo-100 p-6"
      >
        <div class="mb-4 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <div
              :class="[
                'flex h-8 w-8 items-center justify-center rounded-full font-bold text-white',
                index === 0
                  ? 'bg-yellow-500'
                  : index === 1
                    ? 'bg-gray-400'
                    : 'bg-orange-500',
              ]"
            >
              {{ index + 1 }}
            </div>
            <div>
              <h3 class="text-lg font-bold text-gray-900">
                {{ rec.ticker }}
              </h3>
              <p class="text-sm text-gray-600">{{ rec.company }}</p>
            </div>
          </div>
          <ScoreBadge :score="rec.Score" size="large" />
        </div>

        <div class="space-y-2 text-sm">
          <div class="flex justify-between">
            <span class="text-gray-600">Brokerage:</span>
            <span class="font-medium">{{ rec.brokerage }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-600">Acción:</span>
            <ActionBadge :action="rec.action" />
          </div>
          <div class="flex items-center justify-between">
            <span class="text-gray-600">Target:</span>
            <TargetChange
              :from="rec.target_from.toString()"
              :to="rec.target_to.toString()"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import ScoreBadge from '@/components/recommendationsView/ScoreBadge.vue';
import ActionBadge from '@/components/common/ActionBadge.vue';
import TargetChange from '@/components/common/TargetChange.vue';
import type { Recommendation } from '@/models/recommendation';

defineProps<{
  topRecommendations: Recommendation[];
}>();
</script>
