package context

import (
	"context"
	"sync"
	"time"
)

type Context interface {
	context.Context

	// Name gets current context name
	Name() string
	// Cancel cancels current context and all childs contexts.
	// Twice call does nothing
	Cancel(reason error)
	// Finished returns finish channel for context.
	// If it closed - context finished.
	// If fully == true - waits for finish all childs tree before itself
	Finished(fully bool) (is <-chan struct{})
	// ValueSet sets key/value pair in current context
	ValueSet(key interface{}, value interface{})
	// DeadlineSet sets deadline for context and done it if deadline reaches.
	// If reason not set - creates it with description of deadline.
	// If end.IsZero() - clears deadline
	DeadlineSet(end time.Time, reason error)
	// PanicHandlerSet sets panic handler in current context.
	// If it not set in context - panic will be thrown to parent.
	// If not found any handlers - process panics
	PanicHandlerSet(handler func(ctx Context, panicVal interface{}))
	// Child runs new context and inherits all variables.
	// If current context done - child will be canceled also.
	// WARN: If name empty - sets to file:line of call
	Child(name string, worker func(childCtx Context)) (childCtx Context)
	// Childs gets all childs contexts
	Childs() []Context
	// ChildsFinished returns childs finish channel.
	// If it closed - all childs are finished.
	// If fully == true - waits for finish all tree
	ChildsFinished(fully bool) <-chan struct{}
	// ChildsCancel cancels all childs contexts
	ChildsCancel(reason error)
}

type ctx struct {
	id     uint64
	parent *ctx

	name     string
	lock     sync.RWMutex
	done     chan struct{}
	finished chan struct{}
	err      error
	value    sync.Map

	childs struct {
		nextId   uint64
		runs     sync.WaitGroup
		runsTree sync.WaitGroup
		list     map[uint64]*ctx // id:*ctx
	}
}

func newCtx(id uint64, name string, parent *ctx) *ctx {
	c := &ctx{
		id:       id,
		name:     name,
		parent:   parent,
		done:     make(chan struct{}),
		finished: make(chan struct{}),
	}

	c.childs.list = map[uint64]*ctx{}

	return c
}

var Main Context = newCtx(0, "main", nil)
