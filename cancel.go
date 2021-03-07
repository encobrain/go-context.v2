package context

type CancelReason string

func (r CancelReason) String() string {
	return string(r)
}

func (r CancelReason) Error() string {
	return string(r)
}

func (c *ctx) Cancel(reason error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.err != nil {
		return
	}

	c.deadlineCancel()

	if reason == nil {
		reason = CancelReason(c.name + ": canceled")
	}

	c.err = reason

	for _, cc := range c.childs.list {
		cc.Cancel(reason)
	}

	close(c.done)
}
