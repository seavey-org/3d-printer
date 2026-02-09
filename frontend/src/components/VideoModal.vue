<script setup lang="ts">
import { ref, watch, nextTick, onBeforeUnmount } from 'vue'
import type { Timelapse } from '../types/timelapse'

const props = defineProps<{
  timelapse: Timelapse | null
}>()

const emit = defineEmits<{
  close: []
}>()

const videoRef = ref<HTMLVideoElement | null>(null)
const modalRef = ref<HTMLElement | null>(null)

watch(() => props.timelapse, (newVal, oldVal) => {
  if (newVal && !oldVal) {
    // Opening: lock scroll and focus modal
    document.body.style.overflow = 'hidden'
    nextTick(() => {
      modalRef.value?.focus()
    })
  } else if (!newVal && oldVal) {
    // Closing: unlock scroll and pause video
    document.body.style.overflow = ''
    if (videoRef.value) {
      videoRef.value.pause()
    }
  }
})

onBeforeUnmount(() => {
  document.body.style.overflow = ''
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
  // Trap focus within the modal
  if (e.key === 'Tab') {
    const focusable = modalRef.value?.querySelectorAll<HTMLElement>(
      'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
    )
    if (!focusable || focusable.length === 0) return
    const first = focusable[0]!
    const last = focusable[focusable.length - 1]!
    if (e.shiftKey && document.activeElement === first) {
      e.preventDefault()
      last.focus()
    } else if (!e.shiftKey && document.activeElement === last) {
      e.preventDefault()
      first.focus()
    }
  }
}
</script>

<template>
  <Teleport to="body">
    <div
      v-if="timelapse"
      ref="modalRef"
      role="dialog"
      aria-modal="true"
      aria-label="Video player"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/80"
      @click="handleBackdropClick"
      @keydown="handleKeydown"
      tabindex="0"
    >
      <div class="relative w-full max-w-5xl mx-4">
        <button
          @click="emit('close')"
          class="absolute -top-10 right-0 text-white hover:text-gray-300 transition"
          aria-label="Close video player"
        >
          <svg class="h-8 w-8" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
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
