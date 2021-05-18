package tests

import (
	"github.com/encobrain/go-context.v2"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestPanic(t *testing.T) {
	Convey("Catch in self context", t, func() {
		done := make(chan interface{})
		panicVal := &struct{}{}

		ctx := context.Main.Child("self", func(ctx context.Context) {
			ctx.PanicHandlerSet(func(ctx context.Context, panicVal interface{}) {
				done <- ctx
				done <- panicVal
			})

			panic(panicVal)
		}).Go()

		select {
		case v := <-done:
			So(v, ShouldEqual, ctx)

			select {
			case v := <-done:
				So(v, ShouldEqual, panicVal)
			case <-time.After(time.Second):
				So(false, ShouldEqual, true)
			}
		case <-time.After(time.Second):
			So(false, ShouldEqual, true)
		}
	})

	Convey("Catch in parent context", t, func() {
		done := make(chan interface{})
		panicVal := &struct{}{}

		ctx := context.Main.Child("parent", func(ctx context.Context) {
			ctx.PanicHandlerSet(func(ctx context.Context, panicVal interface{}) {
				done <- ctx
				done <- panicVal
			})

			ctx.Child("child", func(ctx context.Context) {
				panic(panicVal)
			}).Go()
		}).Go()

		select {
		case v := <-done:
			So(v, ShouldEqual, ctx)

			select {
			case v := <-done:
				So(v, ShouldEqual, panicVal)
			case <-time.After(time.Second):
				So(false, ShouldEqual, true)
			}
		case <-time.After(time.Second):
			So(false, ShouldEqual, true)
		}
	})

	Convey("Wait childs finished in panic handler success", t, func(c C) {
		ctx := context.Main.Child("parent", func(ctx context.Context) {
			ctx.PanicHandlerSet(func(ctx context.Context, panicVal interface{}) {
				<-ctx.ChildsFinished(true)
			})

			ctx.Child("child", func(ctx context.Context) {
				panic("fail")
			}).Go()
		}).Go()

		select {
		case <-ctx.Finished(true):
		case <-time.After(time.Second):
			So(false, ShouldEqual, true)

		}
	})
}
