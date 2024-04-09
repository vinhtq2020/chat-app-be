package cron

import (
	"sort"
	"time"
)

type Cron struct {
	running bool
	entries []*Entry
	stop    chan struct{}
	add     chan *Entry
}

func NewCron() *Cron {
	return &Cron{
		stop: make(chan struct{}),
		add:  make(chan *Entry),
	}
}

func (c *Cron) AddJob(schedule Schedule, job Job) {
	c.entries = append(c.entries, &Entry{
		Schedule: schedule,
		Job:      job,
	})
}

func (c *Cron) Start() {
	c.running = true
	c.run()
}

func (c *Cron) Add(entry *Entry) {
	if !c.running {
		c.entries = append(c.entries, entry)
		return
	}
	c.add <- entry

}

func (c *Cron) Stop() {
	c.stop <- struct{}{}
}

func (c *Cron) run() {
	now := time.Now().Local()
	for _, entry := range c.entries {
		entry.Next = entry.Schedule.Next(now)
	}

	var effective time.Time

	for {
		sort.Sort(byTime(c.entries))
		if len(c.entries) > 0 {
			effective = c.entries[0].Next
		} else {
			effective = now.AddDate(30, 0, 0)
		}
		select {
		case now = <-time.After(effective.Sub(now)):
			for _, entry := range c.entries {

				if effective != entry.Next {
					break
				}
				entry.Prev = now
				entry.Next = entry.Schedule.Next(now)
				go entry.Job.Run()
			}
		case e := <-c.add:
			e.Next = e.Schedule.Next(now)
			c.entries = append(c.entries, e)

		case <-c.stop:
			return
		}
	}
}
