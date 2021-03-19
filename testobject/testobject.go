/*
Package testobject realizes the ITestObject interface.
*/
package testobject

import (
	"sync"

	"github.com/iostrovok/check"
	"github.com/iostrovok/conveyor/faces"
)

// TestObject is an implementation of faces.ITestObject Interface .
type TestObject struct {
	sync.RWMutex

	mode   bool
	object *check.C
	suffix string
}

// Empty is a constructor of empty test object.
func Empty() faces.ITestObject {
	ob := &TestObject{
		mode: false,
	}

	return ob
}

// New is a constructor.
func New(mode bool, object *check.C, suffix string) faces.ITestObject {
	ob := &TestObject{
		mode:   mode,
		object: object,
		suffix: suffix,
	}

	return ob
}

// IsTestMode is a simple getter.
func (ob *TestObject) IsTestMode() bool {
	return ob.mode
}

// TestObject is a simple getter.
func (ob *TestObject) TestObject() *check.C {
	return ob.object
}

// Suffix is a simple getter.
func (ob *TestObject) Suffix() string {
	return ob.suffix
}
