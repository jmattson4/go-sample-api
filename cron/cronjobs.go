package cron

import (
	"github.com/robfig/cron/v3"
)

//InitJobs ... This function is used to setup Multiple Cronjobs for the server
//	on startup. It Returns a pointer to the cron struct to be closed on server shutdown
func InitJobs(jobs map[string]func()) *cron.Cron {
	c := cron.New()
	for timeString, job := range jobs {
		c.AddFunc(timeString, job)
	}
	c.Start()
	return c
}
