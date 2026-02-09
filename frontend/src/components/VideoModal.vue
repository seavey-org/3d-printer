<script setup lang="ts">
import { ref, watch } from 'vue'
import type { Timelapse } from '../types/timelapse'

const props = defineProps<{
  timelapse: Timelapse | null
}>()

const emit = defineEmits<{
  close: []
}>()

const videoRef = ref<HTMLVideoElement | null>(null)

watch(() => props.timelapse, (newVal) => {
  if (!newVal && videoRef.value) {
    videoRef.value.pause()
  }
})

function handleBackdropClick(e: MouseEvent) {
  if (e.target === e.currentTarget) {
    emit('close')
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    emit('close')
  }
}
</script>

<template>
  <Teleport to="body">
    <div
      v-if="timelapse"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/80"
      @click="handleBackdropClick"
      @keydown="handleKeydown"
      tabindex="0"
    >
      <div class="relative w-full max-w-5xl mx-4">
        <button
          @click="emit('close')"
          class="absolute -top-10 right-0 text-white hover:text-gray-300 transition"
        >
          <svg class="h-8 w-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>

        <video
          ref="videoRef"
          :src="timelapse.url"
          controls
          autoplay
          class="w-full rounded-lg"
        >
          Your browser does not support video playback.
        </video>

        <p class="text-center text-gray-300 mt-2 text-sm">{{ timelapse.filename }}</p>
      </div>
    </div>
  </Teleport>
</template>
