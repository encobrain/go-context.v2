package context

func (c *ctx) Go() (cc Context) {
	cc = c

	defer func() {
		recover()
	}()

	close(c.started)

	return
}

func (c *ctx) Done() (done <-chan struct{}) {
	return c.done
}

func (c *ctx) Err() (err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.err
}

func (c *ctx) ValueSet(key interface{}, value interface{}) Context {
	c.value.Store(key, value)
	return c
}

func (c *ctx) Value(key interface{}) (value interface{}) {
	f := c
	var ok bool

	for {
		value, ok = f.value.Load(key)

		if ok {
			return
		}

		if f = f.parent; f == nil {
			return
		}
	}
}

func (c *ctx) Name() string {
	return c.name
}

func (c *ctx) Finished(fully bool) (is <-chan struct{}) {
	if !fully {
		return c.finished
	}

	ch := make(chan struct{})

	go func() {
		<-c.finished
		c.childs.runsTree.Wait()
		close(ch)
	}()

	return ch
}
