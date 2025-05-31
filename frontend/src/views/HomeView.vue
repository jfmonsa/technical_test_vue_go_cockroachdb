<template>
  <div class="min-h-screen bg-gray-50 py-8">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
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
        class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden"
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
import { onMounted } from "vue";
import { useStockStore } from "../stores/stockStore";
import ErrorAlert from "../components/ErrorAlert.vue";
import SearchControls from "../components/SearchControls.vue";
import PaginationControls from "@/components/PaginationControls.vue";
import HeaderSection from "../components/HeaderSection.vue";
import StatCardSection from "../components/StatCardSection.vue";
import StocksTable from "../components/StocksTable.vue";


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
  fetchStocks,
  setPage,
  setSearch,
  setSorting,
} = stockStore;

// Methods
const refreshData = () => {
  fetchStocks(currentPage);
};


// Lifecycle
onMounted(() => {
  fetchStocks();
});
</script>
