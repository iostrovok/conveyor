package faces

import (
	"github.com/iostrovok/check"
)

/*
	Each handler may have several tests start, stop and run methods.

	example:

	StartTestWithDBConnection(..)
	StartTestWithMocks(..)

	There is the same for run and stop methods:

	RunTestWithDBConnection(..)
	RunTestWithMocks(..)
	StopTestWithDBConnection(..)
	StopTestWithMocks(..)

	also the regular methods should be defined:
	Start(..)
	Run(..)
	Stop(..)

	See handler.go to get method's parameters.
*/

const (
	// StartTestHandlerPrefix is a prefix for tests start method.
	StartTestHandlerPrefix = "StartTest"
	// StopTestHandlerPrefix is a prefix for tests stop method.
	StopTestHandlerPrefix = "StopTest"
	// RunTestHandlerPrefix is a prefix for tests stop method.
	RunTestHandlerPrefix = "RunTest"
)

// ITestObject is an interface to use conveyor in setup and troubleshooting mode.
type ITestObject interface {
	// Return
	Suffix() string
	IsTestMode() bool
	TestObject() *check.C
}
