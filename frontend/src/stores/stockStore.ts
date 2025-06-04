import type { Stock } from '@/models/stock';
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';

export interface StockResponse {
  items: Stock[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

export const useStockStore = defineStore('stock', () => {
  // State
  const stocks = ref<Stock[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);
  // For Pagination
  const currentPage = ref(1);
  const limit = ref(10);
  const total = ref(0);
  const totalPages = ref(0);
  // For Search and Sorting
  const sortField = ref<keyof Stock>('time');
  const sortDirection = ref<'asc' | 'desc'>('desc');
  const searchQuery = ref('');

  // Getters
  /**
   * Computed property that filters the list of stocks based on the current search query.
   *
   * If the search query is empty, all stocks are returned. Otherwise, it returns only the stocks
   * whose `ticker`, `company`, `brokerage`, or `action` fields include the search query (case-insensitive).
   */
  const filteredStocks = computed(() => {
    if (!searchQuery.value) return stocks.value;

    const query = searchQuery.value.toLowerCase();
    return stocks.value.filter(
      (stock) =>
        stock.ticker.toLowerCase().includes(query) ||
        stock.company.toLowerCase().includes(query) ||
        stock.brokerage.toLowerCase().includes(query) ||
        stock.action.toLowerCase().includes(query),
    );
  });

  /**
   * A computed property that returns a sorted array of stocks based on the selected sort field and direction.
   *
   * - If the sort field is "time", stocks are sorted by their date values in ascending or descending order.
   * - For other fields, stocks are sorted lexicographically using `localeCompare`.
   *
   * @remarks
   * The sorting is performed on a shallow copy of the filtered stocks to avoid mutating the original array.
   */
  const sortedStocks = computed(() => {
    const sorted = [...filteredStocks.value];

    sorted.sort((a, b) => {
      const aValue = a[sortField.value];
      const bValue = b[sortField.value];

      if (sortField.value === 'time') {
        const aTime = new Date(aValue).getTime();
        const bTime = new Date(bValue).getTime();
        return sortDirection.value === 'asc' ? aTime - bTime : bTime - aTime;
      }

      let comparison: number;
      if (typeof aValue === 'string' && typeof bValue === 'string') {
        comparison = aValue.localeCompare(bValue);
      } else if (typeof aValue === 'number' && typeof bValue === 'number') {
        comparison = aValue - bValue;
      } else {
        comparison = 0;
      }
      return sortDirection.value === 'asc' ? comparison : -comparison;
    });

    return sorted;
  });

  // Actions
  const fetchStocks = async (page: number = 1) => {
    loading.value = true;
    error.value = null;

    try {
      const response = await fetch(
        `${import.meta.env.VITE_API_BASE_URL}/stocks?page=${page}&limit=${limit.value}`,
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data: StockResponse = await response.json();

      stocks.value = data.items;
      currentPage.value = data.page;
      total.value = data.total;
      totalPages.value = data.totalPages;
    } catch (err) {
      error.value =
        err instanceof Error ? err.message : 'Error fetching stocks';
      console.error('Error fetching stocks:', err);
    } finally {
      loading.value = false;
    }
  };

  const setPage = (page: number) => {
    currentPage.value = page;
    fetchStocks(page);
  };

  const setSearch = (query: string) => {
    searchQuery.value = query;
  };

  const setSorting = (field: keyof Stock) => {
    if (sortField.value === field) {
      sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc';
    } else {
      sortField.value = field;
      sortDirection.value = 'asc';
    }
  };

  return {
    // State
    stocks,
    loading,
    error,
    currentPage,
    limit,
    total,
    totalPages,
    searchQuery,
    sortField,
    sortDirection,

    // Getters
    filteredStocks,
    sortedStocks,

    // Actions
    fetchStocks,
    setPage,
    setSearch,
    setSorting,
  };
});
