<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useTimelapsesStore } from '../stores/timelapses'
import TimelapseCard from '../components/TimelapseCard.vue'
import VideoModal from '../components/VideoModal.vue'
import Pagination from '../components/Pagination.vue'
import type { Timelapse } from '../types/timelapse'

const store = useTimelapsesStore()
const route = useRoute()
const router = useRouter()

// Route is the single source of truth for which timelapse is selected.
// The computed derives the modal state from the URL, no watchers needed.
const selectedTimelapse = computed<Timelapse | null>(() => {
  const param = route.params.filename
  const filename = typeof param === 'string' ? param : ''
  if (!filename) return null
  // Search filteredItems so we don't open modals for items hidden by the size filter
  return store.filteredItems.find(t => t.filename === filename) || null
})

function selectTimelapse(timelapse: Timelapse) {
  router.replace({
    name: 'Timelapses',
    params: { filename: timelapse.filename }
  })
}

function closeModal() {
  router.replace({ name: 'Timelapses' })
}

async function refreshTimelapses() {
  await store.fetchTimelapses()
}

onMounted(() => {
  if (store.allItems.length === 0) {
    refreshTimelapses()
  }
})
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 py-8">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-800 dark:text-white">Timelapses</h1>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
          {{ store.totalCount }} videos
        </p>
      </div>
      <button
        @click="store.toggleSort()"
        class="flex items-center space-x-2 px-3 py-2 rounded-lg bg-white dark:bg-gray-800 shadow-sm text-sm text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 transition"
        aria-label="Toggle sort order"
      >
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12" />
        </svg>
        <span>{{ store.sortOrder === 'newest' ? 'Newest first' : 'Oldest first' }}</span>
      </button>
    </div>

    <!-- Loading -->
    <div v-if="store.loading" class="flex justify-center py-20" role="status" aria-label="Loading timelapses">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <span class="sr-only">Loading...</span>
    </div>

    <!-- Error -->
    <div v-else-if="store.error" class="text-center py-20">
      <p class="text-red-500 dark:text-red-400">{{ store.error }}</p>
      <button
        @click="refreshTimelapses()"
        class="mt-4 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
      >Retry</button>
    </div>

    <!-- Empty state -->
    <div v-else-if="store.totalCount === 0" class="text-center py-20">
      <svg class="h-16 w-16 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
      </svg>
      <p class="text-gray-500 dark:text-gray-400">No timelapses found</p>
    </div>

    <!-- Grid -->
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      <TimelapseCard
        v-for="t in store.paginatedItems"
        :key="t.filename"
        :timelapse="t"
        @select="selectTimelapse($event)"
      />
    </div>

    <!-- Pagination -->
    <Pagination
      :current-page="store.currentPage"
      :total-pages="store.totalPages"
      @update:current-page="store.setPage($event)"
    />

    <!-- Video Modal -->
    <VideoModal
      :timelapse="selectedTimelapse"
      @close="closeModal"
    />
  </div>
</template>
