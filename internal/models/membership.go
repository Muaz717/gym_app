package models

import "time"

type Membership struct {
	Id           int       `json:"-"`
	Number       int       `json:"number" required:"true"`
	RecordingDay time.Time `json:"-"`
	FirstDay     string    `json:"firstDay"`
	LastDay      string    `json:"lastDay"`
	DaysLeft     int       `json:"daysLeft"`
	Person       Person    `json:"person,omitempty"`
}
