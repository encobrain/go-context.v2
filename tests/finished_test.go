package tests

import (
	"github.com/encobrain/go-context.v2"
	. "github.com/smartystreets/goconvey/convey"
	"sync/atomic"
	"testing"
	"time"
)

func TestFinished(t *testing.T) {
	Convey("Only child", t, func() {
		done := int32(0)
		ctx := context.Main.Child("test", func(ctx context.Context) {
			atomic.AddInt32(&done, 1)
		}).Go()

		select {
		case <-ctx.Finished(false):
			So(atomic.LoadInt32(&done), ShouldEqual, 1)
		case <-time.After(time.Second):
			So(false, ShouldEqual, true)
		}
	})

	Convey("Only childs", t, func() {
		done := int32(0)
		ctx := context.Main.Child("test2", func(ctx context.Context) {
			ctx.Child("test2.1", func(ctx context.Context) {
				atomic.AddInt32(&done, 1)
			}).Go()
			ctx.Child("test2.2", func(ctx context.Context) {
				atomic.AddInt32(&done, 2)
			}).Go()
		}).Go()

		<-ctx.Finished(false)

		select {
		case <-time.After(time.Second):
			So(false, ShouldEqual, true)

		case <-ctx.ChildsFinished(false):
			So(atomic.LoadInt32(&done), ShouldEqual, 3)
		}
	})

	Convey("Fully", t, func() {
		done := int32(0)
		ctx := context.Main.Child("test2", func(ctx context.Context) {
			ctx.Child("test2.1", func(ctx context.Context) {
				atomic.AddInt32(&done, 1)
			}).Go()
			ctx.Child("test2.2", func(ctx context.Context) {
				atomic.AddInt32(&done, 2)
			}).Go()

			atomic.AddInt32(&done, 4)
		}).Go()

		select {
		case <-ctx.Finished(true):
			So(atomic.LoadInt32(&done), ShouldEqual, 7)
		case <-time.After(time.Second):
			So(false, ShouldEqual, true)
		}
	})

	Convey("Childs fully", t, func() {
		done := int32(0)
		ctx := context.Main.Child("test2", func(ctx context.Context) {
			ctx.Child("test2.1", func(ctx context.Context) {
				ctx.Child("test.2.1.1", func(ctx context.Context) {
					atomic.AddInt32(&done, 4)
				}).Go()
				atomic.AddInt32(&done, 1)
			}).Go()
			ctx.Child("test2.2", func(ctx context.Context) {
				atomic.AddInt32(&done, 2)
			}).Go()
		}).Go()

		<-ctx.Finished(false)

		select {
		case <-ctx.ChildsFinished(true):
			So(atomic.LoadInt32(&done), ShouldEqual, 7)
		case <-time.After(time.Second):
			So(false, ShouldEqual, true)
		}
	})
}
