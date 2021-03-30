package context

var panicHandlerKey = &struct{ bool }{}

func (c *ctx) panicHandle(panicErr interface{}) {
	sc := c

	defer func() {
		if panicErr = recover(); panicErr != nil {
			if sc != nil {
				sc = sc.parent
			}
			if sc != nil {
				sc.panicHandle(panicErr)
			} else {
				panic(panicErr)
			}
		}
	}()

	for {
		h, ok := sc.value.Load(panicHandlerKey)

		if ok {
			h.(func(ctx Context, panic interface{}))(sc, panicErr)
			return
		}

		sc = sc.parent

		if sc == nil {
			panic(panicErr)
		}
	}
}

func (c *ctx) PanicHandlerSet(handler func(ctx Context, panic interface{})) Context {
	c.value.Store(panicHandlerKey, handler)
	return c
}
