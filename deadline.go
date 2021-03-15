package context

import (
	"time"
)

type DeadlineReason string

func (r DeadlineReason) String() string {
	return string(r)
}

func (r DeadlineReason) Error() string {
	return string(r)
}

var deadlineKey = &struct{ bool }{}

type deadlineValue struct {
	end   time.Time
	timer *time.Timer
}

func (c *ctx) deadlineExceeded(reason error) {
	if reason == nil {
		reason = DeadlineReason(c.name + ": deadline")
	}

	c.Cancel(reason)
}

func (c *ctx) deadlineCancel() {
	dl, ok := c.value.Load(deadlineKey)

	if ok {
		dlv := dl.(*deadlineValue)

		dlv.timer.Stop()
	}
}

func (c *ctx) DeadlineSet(end time.Time, reason error) {
	dl, ok := c.value.Load(deadlineKey)

	if ok {
		dlv := dl.(*deadlineValue)

		if !dlv.timer.Stop() {
			return
		}
	}

	if end.IsZero() {
		c.value.Delete(deadlineKey)
		return
	}

	c.value.Store(deadlineKey, &deadlineValue{
		end: end,
		timer: time.AfterFunc(time.Until(end), func() {
			c.deadlineExceeded(reason)
		}),
	})
}

func (c *ctx) Deadline() (end time.Time, ok bool) {
	dv := c.Value(deadlineKey)

	if dv != nil {
		end = dv.(*deadlineValue).end
		ok = true
	}

	return
}
