<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import Hls from 'hls.js'
import { useCameraStore } from '../stores/camera'
import ConnectionStatus from '../components/ConnectionStatus.vue'

const cameraStore = useCameraStore()

const videoRef = ref<HTMLVideoElement | null>(null)
const containerRef = ref<HTMLElement | null>(null)
const connecting = ref(true)
const isFullscreen = ref(false)

let hls: Hls | null = null
const HLS_URL = '/live/stream.m3u8'
const MAX_RETRIES = 10
const BASE_RETRY_MS = 3000
let retryCount = 0

function initHls() {
  if (!videoRef.value) return

  if (Hls.isSupported()) {
    hls = new Hls({
      lowLatencyMode: true,
      liveSyncDurationCount: 2,
      liveMaxLatencyDurationCount: 5,
      maxBufferLength: 10,
      maxMaxBufferLength: 20,
    })

    hls.loadSource(HLS_URL)
    hls.attachMedia(videoRef.value)

    hls.on(Hls.Events.MANIFEST_PARSED, () => {
      connecting.value = false
      retryCount = 0
      videoRef.value?.play().catch(() => {
        // autoplay blocked, that's fine
      })
    })

    hls.on(Hls.Events.ERROR, (_event, data) => {
      if (data.fatal) {
        connecting.value = true
        if (data.type === Hls.ErrorTypes.NETWORK_ERROR) {
          if (retryCount < MAX_RETRIES) {
            const delay = Math.min(BASE_RETRY_MS * Math.pow(2, retryCount), 30000)
            retryCount++
            setTimeout(() => hls?.startLoad(), delay)
          }
          // After max retries, just stay in connecting state and let
          // the user see the offline overlay via store polling
        } else if (data.type === Hls.ErrorTypes.MEDIA_ERROR) {
          hls?.recoverMediaError()
        } else {
          // Unknown fatal error, destroy and reinit after a delay
          hls?.destroy()
          hls = null
          if (retryCount < MAX_RETRIES) {
            retryCount++
            setTimeout(() => initHls(), 5000)
          }
        }
      }
    })

    hls.on(Hls.Events.FRAG_LOADED, () => {
      connecting.value = false
      retryCount = 0
    })
  } else if (videoRef.value.canPlayType('application/vnd.apple.mpegurl')) {
    // Safari native HLS
    videoRef.value.src = HLS_URL
    videoRef.value.addEventListener('loadedmetadata', () => {
      connecting.value = false
      videoRef.value?.play().catch(() => {})
    })
    videoRef.value.addEventListener('error', () => {
      connecting.value = true
    })
  }
}

function toggleFullscreen() {
  if (!containerRef.value) return

  if (document.fullscreenElement) {
    document.exitFullscreen().catch(() => {})
  } else {
    containerRef.value.requestFullscreen().catch(() => {})
  }
}

function handleFullscreenChange() {
  isFullscreen.value = !!document.fullscreenElement
}

const canSnapshot = ref(true)

function takeSnapshot() {
  if (!videoRef.value) return
  // Guard against zero dimensions (stream offline)
  if (!videoRef.value.videoWidth || !videoRef.value.videoHeight) return

  const canvas = document.createElement('canvas')
  canvas.width = videoRef.value.videoWidth
  canvas.height = videoRef.value.videoHeight

  const ctx = canvas.getContext('2d')
  if (!ctx) return

  ctx.drawImage(videoRef.value, 0, 0)

  try {
    const link = document.createElement('a')
    link.download = `printer-snapshot-${new Date().toISOString().replace(/[:.]/g, '-')}.png`
    link.href = canvas.toDataURL('image/png')
    link.click()
  } catch {
    // toDataURL can throw SecurityError on tainted canvas
    canSnapshot.value = false
  }
}

onMounted(() => {
  initHls()
  cameraStore.startPolling()
  document.addEventListener('fullscreenchange', handleFullscreenChange)
})

onBeforeUnmount(() => {
  hls?.destroy()
  hls = null
  cameraStore.stopPolling()
  document.removeEventListener('fullscreenchange', handleFullscreenChange)
})
</script>

<template>
  <div ref="containerRef" class="relative bg-black" :class="isFullscreen ? 'h-screen' : 'h-[calc(100vh-4rem)]'">
    <!-- Stream offline overlay -->
    <div
      v-if="!cameraStore.online && !connecting"
      class="absolute inset-0 z-10 flex flex-col items-center justify-center bg-gray-900/90"
    >
      <svg class="h-16 w-16 text-gray-500 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
        <line x1="3" y1="3" x2="21" y2="21" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
      </svg>
      <p class="text-xl text-gray-400 font-medium">Stream Offline</p>
      <p class="text-sm text-gray-500 mt-2">The printer camera is not currently streaming</p>
    </div>

    <!-- Video -->
    <video
      ref="videoRef"
      class="w-full h-full object-contain"
      muted
      playsinline
      aria-label="3D printer live camera feed"
    ></video>

    <!-- Controls overlay -->
    <div class="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/60 to-transparent p-4">
      <div class="flex items-center justify-between">
        <ConnectionStatus :online="cameraStore.online" :connecting="connecting" />

        <div class="flex items-center space-x-3">
          <!-- Snapshot (hidden when stream is offline or snapshot failed) -->
          <button
            v-if="canSnapshot && cameraStore.online"
            @click="takeSnapshot"
            class="text-white/80 hover:text-white p-2 rounded transition"
            aria-label="Take snapshot"
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </button>

          <!-- Fullscreen -->
          <button
            @click="toggleFullscreen"
            class="text-white/80 hover:text-white p-2 rounded transition"
            aria-label="Toggle fullscreen"
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
              <path v-if="!isFullscreen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4" />
              <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 9V4m0 5H4m0 0l5-5m-5 5h5m6 5v5m0-5h5m0 0l-5 5m5-5h-5M9 15v5m0-5H4m0 0l5 5m-5-5h5m6-6V4m0 5h5m0 0l-5-5m5 5h-5" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
