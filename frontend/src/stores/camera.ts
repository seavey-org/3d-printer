import { defineStore } from 'pinia'
import { getStreamStatus } from '../services/api'

export const useCameraStore = defineStore('camera', {
  state: () => ({
    online: false,
    lastUpdated: '',
    loading: false,
    pollTimer: null as ReturnType<typeof setInterval> | null
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
      this.pollTimer = setInterval(() => this.fetchStatus(), 10_000)
    },

    stopPolling() {
      if (this.pollTimer) {
        clearInterval(this.pollTimer)
        this.pollTimer = null
      }
    }
  }
})
