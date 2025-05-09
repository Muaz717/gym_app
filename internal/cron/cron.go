package cron

import (
	"context"
	"github.com/robfig/cron/v3"
	personSubService "gym_app/internal/services/person_sub"
)

type CronJobs struct {
	cronScheduler    *cron.Cron
	personSubService *personSubService.PersonSubService
}

func New(personSubService *personSubService.PersonSubService) *CronJobs {
	return &CronJobs{
		cronScheduler:    cron.New(),
		personSubService: personSubService,
	}
}

func (c *CronJobs) Start(ctx context.Context) {

	c.cronScheduler.AddFunc("@daily", func() {
		err := c.personSubService.UpdateStatuses(ctx)
		if err != nil {
			panic(err)
		}
	})

	c.cronScheduler.Start()
}

func (c *CronJobs) Stop() {
	c.cronScheduler.Stop()
}
