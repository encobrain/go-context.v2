package tests

import (
	"fmt"
	"github.com/encobrain/go-context.v2"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestCancel(t *testing.T) {
	Convey("Cancel single", t, func() {
		done := make(chan interface{})
		cancelReason := fmt.Errorf("cancel")

		ctx := context.Main.Child("cancel", func(ctx context.Context) {
			<-ctx.Done()
			done <- ctx.Err()
		})

		ctx.Cancel(cancelReason)

		select {
		case r := <-done:
			So(r, ShouldEqual, cancelReason)
		case <-time.After(time.Second):
			So(false, ShouldEqual, true)
		}
	})

	Convey("Cancel twice ignored", t, func() {
		done := make(chan interface{})
		canceled := make(chan interface{})
		cancelReason1 := fmt.Errorf("cancel1")
		cancelReason2 := fmt.Errorf("cancel2")

		ctx := context.Main.Child("cancel", func(ctx context.Context) {
			<-canceled
			<-ctx.Done()
			done <- ctx.Err()
		})

		ctx.Cancel(cancelReason1)
		ctx.Cancel(cancelReason2)
		canceled <- 1

		select {
		case r := <-done:
			So(r, ShouldEqual, cancelReason1)
		case <-time.After(time.Second):
			So(false, ShouldEqual, true)
		}
	})

	Convey("Cancel parent cancels childs also", t, func() {
		done := make(chan interface{})
		cancelReason := fmt.Errorf("cancel")

		ctx := context.Main.Child("parent", func(ctx context.Context) {
			ctx.Child("child1", func(ctx context.Context) {
				<-ctx.Done()
				done <- ctx.Err()
			})

			ctx.Child("child2", func(ctx context.Context) {
				<-ctx.Done()
				done <- ctx.Err()
			})
		})

		ctx.Cancel(cancelReason)

		select {
		case r := <-done:
			So(r, ShouldEqual, cancelReason)
		case <-time.After(time.Second):
			So(false, ShouldEqual, true)
		}

		select {
		case r := <-done:
			So(r, ShouldEqual, cancelReason)
		case <-time.After(time.Second):
			So(false, ShouldEqual, true)
		}
	})

	Convey("New child on canceled parent should canceled too", t, func() {
		done := make(chan interface{})
		cancelReason := fmt.Errorf("cancel")

		ctx := context.Main.Child("parent", func(ctx context.Context) {})

		ctx.Cancel(cancelReason)

		ctx.Child("child", func(ctx context.Context) {
			<-ctx.Done()
			done <- ctx.Err()
		})

		select {
		case r := <-done:
			So(r, ShouldEqual, cancelReason)
		case <-time.After(time.Second):
			So(false, ShouldEqual, true)
		}
	})
}
