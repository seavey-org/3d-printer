<script setup lang="ts">
defineProps<{
  online: boolean
  connecting: boolean
}>()
</script>

<template>
  <div class="flex items-center space-x-2">
    <span
      aria-hidden="true"
      class="h-3 w-3 rounded-full"
      :class="{
        'bg-green-500': online && !connecting,
        'bg-yellow-500 animate-pulse': connecting,
        'bg-red-500': !online && !connecting
      }"
    ></span>
    <span class="text-sm text-gray-600 dark:text-gray-300">
      {{ connecting ? 'Connecting...' : online ? 'Live' : 'Offline' }}
    </span>
  </div>
</template>
