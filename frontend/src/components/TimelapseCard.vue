<script setup lang="ts">
import type { Timelapse } from '../types/timelapse'

defineProps<{
  timelapse: Timelapse
}>()

defineEmits<{
  select: [timelapse: Timelapse]
}>()

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function formatSize(bytes: number): string {
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(0)} KB`
  if (bytes < 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(2)} GB`
}
</script>

<template>
  <button
    class="bg-white dark:bg-gray-800 rounded-lg shadow-sm overflow-hidden hover:shadow-md transition cursor-pointer text-left w-full"
    @click="$emit('select', timelapse)"
  >
    <div class="aspect-video bg-gray-200 dark:bg-gray-700 relative">
      <img
        v-if="timelapse.thumbnailUrl"
        :src="timelapse.thumbnailUrl"
        :alt="timelapse.filename"
        loading="lazy"
        class="w-full h-full object-cover"
      />
      <div v-else class="w-full h-full flex items-center justify-center text-gray-400 dark:text-gray-500">
        <svg class="h-12 w-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
      </div>
      <div class="absolute bottom-2 right-2 bg-black/70 text-white text-xs px-2 py-1 rounded">
        {{ formatSize(timelapse.size) }}
      </div>
    </div>
    <div class="p-3">
      <p class="text-sm text-gray-800 dark:text-gray-200 truncate">{{ timelapse.filename }}</p>
      <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">{{ formatDate(timelapse.date) }}</p>
    </div>
  </button>
</template>
