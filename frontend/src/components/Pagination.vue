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
  const pages: { key: string; value: number }[] = []
  const total = props.totalPages
  const current = props.currentPage

  // Always show first page
  pages.push({ key: 'first', value: 1 })

  // Show pages around current
  const start = Math.max(2, current - 1)
  const end = Math.min(total - 1, current + 1)

  if (start > 2) pages.push({ key: 'ellipsis-start', value: -1 })
  for (let i = start; i <= end; i++) pages.push({ key: `page-${i}`, value: i })
  if (end < total - 1) pages.push({ key: 'ellipsis-end', value: -1 })

  // Always show last page if more than 1 page
  if (total > 1) pages.push({ key: 'last', value: total })

  return pages
})
</script>

<template>
  <nav v-if="totalPages > 1" class="flex items-center justify-center space-x-1 mt-8" aria-label="Pagination">
    <button
      :disabled="currentPage <= 1"
      @click="emit('update:currentPage', currentPage - 1)"
      class="px-3 py-2 rounded text-sm text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700 disabled:opacity-40 disabled:cursor-not-allowed"
    >Prev</button>

    <template v-for="page in visiblePages" :key="page.key">
      <span v-if="page.value === -1" class="px-2 text-gray-400" aria-hidden="true">...</span>
      <button
        v-else
        @click="emit('update:currentPage', page.value)"
        class="px-3 py-2 rounded text-sm transition"
        :class="page.value === currentPage
          ? 'bg-blue-600 text-white'
          : 'text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700'"
        :aria-current="page.value === currentPage ? 'page' : undefined"
      >{{ page.value }}</button>
    </template>

    <button
      :disabled="currentPage >= totalPages"
      @click="emit('update:currentPage', currentPage + 1)"
      class="px-3 py-2 rounded text-sm text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700 disabled:opacity-40 disabled:cursor-not-allowed"
    >Next</button>
  </nav>
</template>
