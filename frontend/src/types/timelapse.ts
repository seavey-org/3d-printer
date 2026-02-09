export interface Timelapse {
  filename: string
  url: string
  thumbnailUrl: string
  size: number
  date: string
}

export interface StreamStatus {
  online: boolean
  lastUpdated: string
}
