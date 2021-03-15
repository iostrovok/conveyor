package faces

import (
	"github.com/iostrovok/check"
)

// File describes the test object.

type ITestObject interface {
	// Return
	Suffix() string
	IsTestMode() bool
	TestObject() *check.C
}
