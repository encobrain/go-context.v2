package tests

import (
	"fmt"
	"github.com/encobrain/go-context.v2"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestDeadline(t *testing.T) {
	Convey("Set is ok on get", t, func() {
		ctx := context.Main.Child("name", func(ctx context.Context) {}).Go()

		dlv := time.Now().Add(time.Millisecond * 100)

		ctx.DeadlineSet(dlv, nil)

		dlt, ok := ctx.Deadline()

		So(ok, ShouldEqual, true)
		So(dlt, ShouldEqual, dlv)
	})

	Convey("Set is fired", t, func() {
		done := make(chan interface{})
		dlReason := fmt.Errorf("deadline")

		ctx := context.Main.Child("name", func(ctx context.Context) {
			<-ctx.Done()
			done <- ctx.Err()
		}).Go()

		ctx.DeadlineSet(time.Now().Add(time.Millisecond*100), dlReason)

		select {
		case r := <-done:
			So(r, ShouldEqual, dlReason)
		case <-time.After(time.Second):
			So(false, ShouldEqual, true)

		}
	})

	Convey("Set zero - cancel deadline", t, func() {
		done := make(chan interface{})

		ctx := context.Main.Child("name", func(ctx context.Context) {
			select {
			case <-ctx.Done():
				done <- ctx.Err()
			case <-time.After(time.Millisecond * 150):
				done <- "ok"
			}
		}).DeadlineSet(time.Now().Add(time.Millisecond*100), nil).Go()

		<-time.After(time.Millisecond * 50)
		ctx.DeadlineSet(time.Time{}, nil)

		select {
		case r := <-done:
			So(r, ShouldEqual, "ok")
		case <-time.After(time.Second):
			So(false, ShouldEqual, true)
		}
	})

	Convey("Set childs deadline longer than parent fires with parent reason", t, func() {
		done := make(chan interface{})
		dlReason := fmt.Errorf("deadline")

		ctx := context.Main.Child("parent", func(ctx context.Context) {
			done <- ctx.Child("child", func(ctx context.Context) {
				<-ctx.Done()
				done <- ctx.Err()
			}).Go()
		}).Go()

		childCtx := (<-done).(context.Context)

		dlv := time.Now().Add(time.Millisecond * 100)
		dlvc := time.Now().Add(time.Millisecond * 200)

		childCtx.DeadlineSet(dlvc, nil)
		ctx.DeadlineSet(dlv, dlReason)

		select {
		case r := <-done:
			So(time.Now().Before(dlvc), ShouldEqual, true)
			So(r, ShouldEqual, dlReason)
		case <-time.After(time.Second):
			So(false, ShouldEqual, true)
		}
	})
}
