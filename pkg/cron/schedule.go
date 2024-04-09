package cron

import (
	"errors"
	"time"
)

// type Schedule interface {
// 	Next(t time.Time) time.Time
// }

// type AtSchedule interface {
// 	At(t string) Schedule
// 	Schedule
// }

// type periodicSchedule struct {
// 	period time.Duration
// }

// func Every(p time.Duration) AtSchedule {
// 	if p < time.Second {
// 		p = time.Second
// 	}

// 	p = p - time.Duration(p.Nanoseconds())%time.Second // truncate up to second

// 	return &periodicSchedule{
// 		period: p,
// 	}
// }

// func (ps periodicSchedule) Next(t time.Time) time.Time {
// 	return t.Truncate(time.Second).Add(ps.period)
// }

// func (ps periodicSchedule) At(t string) Schedule {
// 	if ps.period < time.Hour*24 {
// 		panic("period must be at least in days")
// 	}

// 	// parse t naively
// 	h, m, err := parse(t)

// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	return &atSchedule{
// 		period: ps.period,
// 		hh:     h,
// 		mm:     m,
// 	}
// }

// type atSchedule struct {
// 	period time.Duration
// 	hh     int
// 	mm     int
// }

// func (as atSchedule) reset(t time.Time) time.Time {
// 	return time.Date(t.Year(), t.Month(), t.Day(), as.hh, as.mm, 0, 0, time.UTC)
// }

// func (as atSchedule) Next(t time.Time) time.Time {
// 	next := as.reset(t)
// 	if t.After(next) {
// 		return next.Add(as.period)
// 	}
// 	return next
// }

// func parse(hhmm string) (hh int, mm int, err error) {

// 	hh = int(hhmm[0]-'0')*10 + int(hhmm[1]-'0')
// 	mm = int(hhmm[3]-'0')*10 + int(hhmm[4]-'0')

// 	if hh < 0 || hh > 24 {
// 		hh, mm = 0, 0
// 		err = errors.New("invalid hh format")
// 	}
// 	if mm < 0 || mm > 59 {
// 		hh, mm = 0, 0
// 		err = errors.New("invalid mm format")
// 	}

// 	return
// }

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
