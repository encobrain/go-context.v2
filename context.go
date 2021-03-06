package context

import (
	"context"
	"time"
)

type Context interface {
	context.Context

	// Name gets current context name
	Name() string
	// Cancel cancels current context and all childs contexts
	// Twice call does nothing
	Cancel(reason error)
	// Finished returns finish channel for context. If it closed - context finished.
	// If fully == true - waits for finish all childs tree before itself
	Finished(fully bool) (is <-chan struct{})
	// ValueSet sets key/value pair in current context
	ValueSet(key interface{}, value interface{})
	// DeadlineSet sets deadline for context and cancel it if deadline reaches.
	// If reason not set - creates it with description of deadline
	DeadlineSet(end time.Time, reason error)
	// PanicHandlerSet sets panic handler in current context.
	// If it not set in context - panic will be thrown to parent.
	// If not found any handlers - process panics
	PanicHandlerSet(handler func(ctx Context, panic interface{}))
	// Child runs new context and inherits all variables.
	// If current context cancel - child will be canceled also.
	// WARN: If name empty - sets to file:line of call
	Child(name string, worker func(childCtx Context)) (childCtx Context)
	// Childs gets all childs contexts
	Childs() []Context
	// ChildsFinished returns childs finish channel. If it closed - all childs are finished
	// If fully == true - waits for finish all tree
	ChildsFinished(fully bool) <-chan struct{}
	// ChildsCancel cancels all childs contexts
	ChildsCancel(reason error)
}
