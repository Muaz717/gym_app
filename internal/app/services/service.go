package services

import (
	"gym_app/internal/models"
	"time"
)

func EnrichMembership(membership models.Membership) models.Membership {
	membership.FirstDay = membership.RecordingDay.Format("02-01-2006")
	membership.LastDay = membership.RecordingDay.AddDate(0, 0, 30).Format("02-01-2006")
	membership.DaysLeft = membership.RecordingDay.YearDay() + 30 - time.Now().YearDay()

	return membership
}
