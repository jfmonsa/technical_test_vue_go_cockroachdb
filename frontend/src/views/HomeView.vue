<template>
  <div class="min-h-screen bg-gray-50 py-8">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <HeaderSection />

      <SearchControls
        :searchQuery="searchQuery"
        :loading="loading"
        @update:searchQuery="setSearch"
        @refresh="refreshData"
      />

      <ErrorAlert v-if="error" :error="error" />

      <StatCardSection
        :total="total"
        :searchResults="sortedStocks.length"
        :currentPage="currentPage"
        :totalPages="totalPages"
      />

      <!-- Table -->
      <div
        class="overflow-hidden rounded-lg border border-gray-200 bg-white shadow-sm"
      >
        <StocksTable
          :stocks="sortedStocks"
          :loading="loading"
          :sortField="sortField"
          :sortDirection="sortDirection"
          :setSorting="setSorting"
        />
        <!-- Pagination -->
        <PaginationControls
          :currentPage="currentPage"
          :totalPages="totalPages"
          :limit="limit"
          :total="total"
          @page="setPage"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useStockStore } from '../stores/stockStore';
import ErrorAlert from '../components/common/ErrorAlert.vue';
import SearchControls from '../components/stocksView/SearchControls.vue';
import PaginationControls from '@/components/stocksView/PaginationControls.vue';
import HeaderSection from '../components/stocksView/HeaderSection.vue';
import StatCardSection from '../components/stocksView/StatCardSection.vue';
import StocksTable from '../components/stocksView/StocksTable.vue';
import { storeToRefs } from 'pinia';

const stockStore = useStockStore();

// Destructure reactive properties
const {
  loading,
  error,
  currentPage,
  limit,
  total,
  totalPages,
  searchQuery,
  sortField,
  sortDirection,
  sortedStocks,
} = storeToRefs(stockStore);

const { fetchStocks, setPage, setSearch, setSorting } = stockStore;
// Methods
const refreshData = () => {
  fetchStocks(currentPage.value);
};

// Lifecycle
onMounted(() => {
  fetchStocks();
});
</script>
