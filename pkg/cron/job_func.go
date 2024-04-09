package cron

type JobFunc func()

// Run calls j()
func (j JobFunc) Run() {
	j()
}
