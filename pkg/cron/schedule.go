package cron

import (
	"errors"
	"time"
)

type Schedule interface {
	Next(t time.Time) time.Time
}

type schedule struct {
	periodTime time.Duration
	hh         int
	mm         int
}

func NewSchedule(periodTime time.Duration, opts ...interface{}) (*schedule, error) {
	var err error
	hh, mm := 0, 0
	if len(opts) > 0 {
		hhmm := opts[0].(string)
		hh, mm, err = parse(hhmm)
		if err != nil {
			return nil, err
		}
	}

	return &schedule{
		periodTime: periodTime,
		hh:         hh,
		mm:         mm,
	}, nil
}

// format hh:mm
func parse(hhmm string) (int, int, error) {
	hh := int(hhmm[0]-'0')*10 + int(hhmm[1]-'0')
	mm := int(hhmm[3]-'0')*10 + int(hhmm[4]-'0')

	if hh < 0 && hh >= 24 {
		return hh, mm, errors.New("invalid hour format")
	} else if mm < 0 && mm >= 60 {
		return hh, mm, errors.New("invalid minute format")
	}

	return hh, mm, nil
}

func (s *schedule) Next(t time.Time) time.Time {
	return t.Truncate(time.Second).Add(s.periodTime)
}
