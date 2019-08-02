package sysdes

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewSystemDes(t *testing.T) {
	Convey("new system from string", t, func() {
		sys, err := NewSystemDes(`/Users/rg/Projects/storage-system/bff/admin`)
		So(err, ShouldBeNil)
		So(sys, ShouldNotBeNil)
		t.Log(sys.Name, sys.DirNode.Abs())
	})
}
