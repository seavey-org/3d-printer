import { defineStore } from 'pinia'
import { getTimelapses } from '../services/api'
import type { Timelapse } from '../types/timelapse'

const PAGE_SIZE = 24
const MIN_FILE_SIZE = 100 * 1024 // 100KB

type SortOrder = 'newest' | 'oldest'

export const useTimelapsesStore = defineStore('timelapses', {
  state: () => ({
    allItems: [] as Timelapse[],
    sortOrder: 'newest' as SortOrder,
    currentPage: 1,
    loading: false,
    error: null as string | null
  }),

  getters: {
    filteredItems: (state): Timelapse[] => {
      // Filter out tiny files (failed captures)
      let items = state.allItems.filter(t => t.size >= MIN_FILE_SIZE)

      // Sort
      if (state.sortOrder === 'oldest') {
        items = [...items].reverse()
      }

      return items
    },

    totalPages(): number {
      return Math.max(1, Math.ceil(this.filteredItems.length / PAGE_SIZE))
    },

    paginatedItems(): Timelapse[] {
      const start = (this.currentPage - 1) * PAGE_SIZE
      return this.filteredItems.slice(start, start + PAGE_SIZE)
    },

    totalCount(): number {
      return this.filteredItems.length
    }
  },

  actions: {
    async fetchTimelapses() {
      this.loading = true
      this.error = null
      try {
        this.allItems = await getTimelapses()
      } catch (err) {
        this.error = err instanceof Error ? err.message : 'Failed to load timelapses'
      } finally {
        this.loading = false
      }
    },

    toggleSort() {
      this.sortOrder = this.sortOrder === 'newest' ? 'oldest' : 'newest'
      this.currentPage = 1
    },

    setPage(page: number) {
      if (page >= 1 && page <= this.totalPages) {
        this.currentPage = page
      }
    }
  }
})
