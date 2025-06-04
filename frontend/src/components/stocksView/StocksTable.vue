<template>
  <div class="overflow-x-auto">
    <table class="min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-50">
        <tr>
          <th @click="setSorting('ticker')" class="table-header">
            <div class="flex items-center gap-2">
              Ticker
              <SortIcon
                :field="'ticker'"
                :current-field="sortField"
                :direction="sortDirection"
              />
            </div>
          </th>
          <th @click="setSorting('company')" class="table-header">
            <div class="flex items-center gap-2">
              Company
              <SortIcon
                :field="'company'"
                :current-field="sortField"
                :direction="sortDirection"
              />
            </div>
          </th>
          <th @click="setSorting('brokerage')" class="table-header">
            <div class="flex items-center gap-2">
              Brokerage
              <SortIcon
                :field="'brokerage'"
                :current-field="sortField"
                :direction="sortDirection"
              />
            </div>
          </th>
          <th @click="setSorting('action')" class="table-header">
            <div class="flex items-center gap-2">
              Action
              <SortIcon
                :field="'action'"
                :current-field="sortField"
                :direction="sortDirection"
              />
            </div>
          </th>
          <th class="table-header">Rating</th>
          <th class="table-header">Target</th>
          <th @click="setSorting('time')" class="table-header">
            <div class="flex items-center gap-2">
              Date
              <SortIcon
                :field="'time'"
                :current-field="sortField"
                :direction="sortDirection"
              />
            </div>
          </th>
        </tr>
      </thead>
      <tbody class="divide-y divide-gray-200 bg-white">
        <tr v-if="loading" class="animate-pulse">
          <td colspan="7" class="table-cell py-8 text-center">
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
              <span class="ml-2">Cargando datos...</span>
            </div>
          </td>
        </tr>
        <tr v-else-if="stocks.length === 0" class="text-center">
          <td colspan="7" class="table-cell py-8 text-gray-500">
            No results found
          </td>
        </tr>
        <tr
          v-else
          v-for="stock in stocks"
          :key="stock.ticker"
          class="transition-colors hover:bg-gray-50"
        >
          <td class="table-cell">
            <span class="font-mono font-semibold text-blue-600">{{
              stock.ticker
            }}</span>
          </td>
          <td class="table-cell">
            <div class="font-medium text-gray-900">
              {{ stock.company }}
            </div>
          </td>
          <td class="table-cell">
            <span class="text-gray-700">{{ stock.brokerage }}</span>
          </td>
          <td class="table-cell">
            <ActionBadge :action="stock.action" />
          </td>
          <td class="table-cell">
            <RatingChange :from="stock.rating_from" :to="stock.rating_to" />
          </td>
          <td class="table-cell">
            <TargetChange
              :from="String(stock.target_from)"
              :to="String(stock.target_to)"
            />
          </td>
          <td class="table-cell">
            <time class="text-sm text-gray-600">{{
              formatDate(stock.time)
            }}</time>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import SortIcon from '../SortIcon.vue';
import ActionBadge from '../common/ActionBadge.vue';
import RatingChange from '../common/RatingChange.vue';
import TargetChange from '../common/TargetChange.vue';
import { formatDate } from '../../utils/dates';
import type { Stock } from '@/models/stock';

defineProps<{
  stocks: Stock[];
  loading: boolean;
  sortField: string;
  sortDirection: 'asc' | 'desc';
  setSorting: (field: keyof Stock) => void;
}>();
</script>
