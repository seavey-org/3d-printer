<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  currentPage: number
  totalPages: number
}>()

const emit = defineEmits<{
  'update:currentPage': [page: number]
}>()

const visiblePages = computed(() => {
  const pages: number[] = []
  const total = props.totalPages
  const current = props.currentPage

  // Always show first page
  pages.push(1)

  // Show pages around current
  const start = Math.max(2, current - 1)
  const end = Math.min(total - 1, current + 1)

  if (start > 2) pages.push(-1) // ellipsis
  for (let i = start; i <= end; i++) pages.push(i)
  if (end < total - 1) pages.push(-1) // ellipsis

  // Always show last page if more than 1 page
  if (total > 1) pages.push(total)

  return pages
})
</script>

<template>
  <div v-if="totalPages > 1" class="flex items-center justify-center space-x-1 mt-8">
    <button
      :disabled="currentPage <= 1"
      @click="emit('update:currentPage', currentPage - 1)"
      class="px-3 py-2 rounded text-sm text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700 disabled:opacity-40 disabled:cursor-not-allowed"
    >Prev</button>

    <template v-for="(page, i) in visiblePages" :key="i">
      <span v-if="page === -1" class="px-2 text-gray-400">...</span>
      <button
        v-else
        @click="emit('update:currentPage', page)"
        class="px-3 py-2 rounded text-sm transition"
        :class="page === currentPage
          ? 'bg-blue-600 text-white'
          : 'text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700'"
      >{{ page }}</button>
    </template>

    <button
      :disabled="currentPage >= totalPages"
      @click="emit('update:currentPage', currentPage + 1)"
      class="px-3 py-2 rounded text-sm text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700 disabled:opacity-40 disabled:cursor-not-allowed"
    >Next</button>
  </div>
</template>
