package timer

import "time"

type Timer struct {
}

func (p *Timer) Init() error {
	return nil
}

func (p *Timer) Now() time.Time {
	// TODO improve me
	return time.Now()
}
