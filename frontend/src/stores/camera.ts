import { defineStore } from 'pinia'
import { getStreamStatus } from '../services/api'

let pollTimer: ReturnType<typeof setInterval> | null = null
let visibilityHandler: (() => void) | null = null

export const useCameraStore = defineStore('camera', {
  state: () => ({
    online: false,
    lastUpdated: '',
    loading: false,
  }),

  actions: {
    async fetchStatus() {
      this.loading = true
      try {
        const status = await getStreamStatus()
        this.online = status.online
        this.lastUpdated = status.lastUpdated
      } catch {
        this.online = false
      } finally {
        this.loading = false
      }
    },

    startPolling() {
      this.stopPolling()
      this.fetchStatus()
      pollTimer = setInterval(() => this.fetchStatus(), 10_000)

      // Pause polling when tab is hidden, resume when visible
      visibilityHandler = () => {
        if (document.hidden) {
          if (pollTimer) {
            clearInterval(pollTimer)
            pollTimer = null
          }
        } else {
          if (pollTimer) clearInterval(pollTimer)
          this.fetchStatus()
          pollTimer = setInterval(() => this.fetchStatus(), 10_000)
        }
      }
      document.addEventListener('visibilitychange', visibilityHandler)
    },

    stopPolling() {
      if (pollTimer) {
        clearInterval(pollTimer)
        pollTimer = null
      }
      if (visibilityHandler) {
        document.removeEventListener('visibilitychange', visibilityHandler)
        visibilityHandler = null
      }
    }
  }
})
