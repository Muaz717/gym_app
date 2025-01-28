package models

import "time"

type Membership struct {
	Id           int
	Number       int `json:"number"`
	RecordingDay time.Duration
	Person       string `json:"person"`
}
