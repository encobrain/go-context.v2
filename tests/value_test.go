package tests

import (
	"github.com/encobrain/go-context.v2"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
)

func TestValue(t *testing.T) {
	setV := rand.Float64()

	context.Main.ValueSet("test", setV)

	Convey("Set in main. Get in main", t, func() {
		v := context.Main.Value("test")

		So(v, ShouldEqual, setV)
	})

	Convey("Set in main. Get in child", t, func() {
		done := make(chan interface{})
		context.Main.Child("child", func(ctx context.Context) {
			done <- ctx.Value("test")
		}).Go()

		So(<-done, ShouldEqual, setV)
	})

	Convey("Set in child. Get in child", t, func() {
		setV2 := setV + 1
		done := make(chan interface{})

		context.Main.Child("child2", func(ctx context.Context) {
			ctx.ValueSet("test", setV2)

			done <- ctx.Value("test")
		}).Go()

		So(<-done, ShouldEqual, setV2)
		So(context.Main.Value("test"), ShouldEqual, setV)
	})

	Convey("Set on child. Get in child", t, func() {
		setV3 := setV + 2
		done := make(chan interface{})

		context.Main.Child("child3", func(ctx context.Context) {
			done <- ctx.Value("test")
		}).ValueSet("test", setV3).Go()

		So(<-done, ShouldEqual, setV3)
		So(context.Main.Value("test"), ShouldEqual, setV)
	})

	Convey("Set in main. Get in child. Change in main. Get in child", t, func() {
		setV4 := setV + 3
		done := make(chan interface{})
		changed := make(chan interface{})

		ctx := context.Main.Child("child4", func(ctx context.Context) {
			ctx.Child("child4.1", func(ctx context.Context) {
				done <- ctx.Value("test")
				<-changed
				done <- ctx.Value("test")
			}).Go()
		}).ValueSet("test", setV4).Go()

		So(<-done, ShouldEqual, setV4)
		setV5 := setV4 + 1
		ctx.ValueSet("test", setV5)
		close(changed)
		So(<-done, ShouldEqual, setV5)
	})
}
