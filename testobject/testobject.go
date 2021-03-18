package testobject

import (
	"sync"

	"github.com/iostrovok/check"
	"github.com/iostrovok/conveyor/faces"
)

/*
	Package realizes the ITestObject interface.
*/

type TestObject struct {
	sync.RWMutex

	mode   bool
	object *check.C
	suffix string
}

func Empty() faces.ITestObject {
	ob := &TestObject{
		mode: false,
	}

	return ob
}

func New(mode bool, object *check.C, suffix string) faces.ITestObject {
	ob := &TestObject{
		mode:   mode,
		object: object,
		suffix: suffix,
	}

	return ob
}

func (ob *TestObject) IsTestMode() bool {
	return ob.mode
}

func (ob *TestObject) TestObject() *check.C {
	return ob.object
}

func (ob *TestObject) Suffix() string {
	return ob.suffix
}
