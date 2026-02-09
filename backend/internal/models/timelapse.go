package models

import "time"

type Timelapse struct {
	Filename     string    `json:"filename"`
	URL          string    `json:"url"`
	ThumbnailURL string    `json:"thumbnailUrl"`
	Size         int64     `json:"size"`
	Date         time.Time `json:"date"`
}

type StreamStatus struct {
	Online      bool      `json:"online"`
	LastUpdated time.Time `json:"lastUpdated"`
}
