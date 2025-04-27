package models

import "time"

// PersonSubscription представляет подписку клиента на абонемент
type PersonSubscription struct {
	Number         string    `json:"number"`          // Номер абонемента
	ClientID       int64     `json:"client_id"`       // ID клиента
	SubscriptionID string    `json:"subscription_id"` // ID абонемента
	StartDate      time.Time `json:"start_date"`      // Дата начала
	EndDate        time.Time `json:"end_date"`        // Дата окончания
	Status         string    `json:"status"`          // Статус абонемента (active/frozen/completed)
}
