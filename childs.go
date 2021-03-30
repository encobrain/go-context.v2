package context

import (
	"runtime"
	"strconv"
)

func (c *ctx) Child(name string, worker func(childCtx Context)) (childCtx Context) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if name == "" {
		_, file, line, _ := runtime.Caller(1)
		name = file + ":" + strconv.Itoa(line)
	}

	id := c.childs.nextId
	c.childs.nextId++

	newCtx := newCtx(id, name, c)

	if c.err != nil {
		newCtx.Cancel(c.err)
	}

	c.childs.list[id] = newCtx

	c.childs.runs.Add(1)
	c.childs.runsTree.Add(1)
	go newCtx.run(worker)

	return newCtx
}

func (c *ctx) run(worker func(childCtx Context)) {
	<-c.started

	defer func() {
		if panicErr := recover(); panicErr != nil {
			c.panicHandle(panicErr)
		}

		close(c.finished)

		if par := c.parent; par != nil {
			par.childFinished(c)
			c.childs.runsTree.Wait()
			par.childTreeFinished(c)
		}
	}()

	worker(c)
}

func (c *ctx) childFinished(childCtx *ctx) {
	c.childs.runs.Done()
}

func (c *ctx) childTreeFinished(childCtx *ctx) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.childs.list, childCtx.id)
	c.childs.runsTree.Done()
}

func (c *ctx) Childs() []Context {
	c.lock.Lock()
	defer c.lock.Unlock()

	ret := make([]Context, len(c.childs.list))

	i := 0

	for _, cc := range c.childs.list {
		ret[i] = cc
		i++
	}

	return ret
}

func (c *ctx) ChildsFinished(fully bool) (is <-chan struct{}) {
	ch := make(chan struct{})

	go func() {
		w := &c.childs.runs

		if fully {
			w = &c.childs.runsTree
		}

		w.Wait()
		close(ch)
	}()

	return ch
}

func (c *ctx) ChildsCancel(reason error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if reason == nil {
		reason = CancelReason(c.name + ": child cancel")
	}

	for _, cc := range c.childs.list {
		cc.Cancel(reason)
	}
}
