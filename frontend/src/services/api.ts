import type { Timelapse, StreamStatus } from '../types/timelapse'

const BASE_URL = '/api'

async function fetchJSON<T>(path: string): Promise<T> {
  const response = await fetch(`${BASE_URL}${path}`)
  if (!response.ok) {
    throw new Error('Failed to load data. Please try again later.')
  }
  return response.json() as Promise<T>
}

export async function getTimelapses(): Promise<Timelapse[]> {
  return fetchJSON<Timelapse[]>('/timelapses')
}

export async function getStreamStatus(): Promise<StreamStatus> {
  return fetchJSON<StreamStatus>('/stream/status')
}
