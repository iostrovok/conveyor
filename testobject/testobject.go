/*
	Package realizes the IItem interface.
*/
package testobject

import (
	"github.com/iostrovok/check"
	"github.com/iostrovok/conveyor/faces"
	"sync"
)

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
