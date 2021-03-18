package testobject_test

import (
	"testing"

	. "github.com/iostrovok/check"

	_ "github.com/iostrovok/conveyor/testobject"
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestService(t *testing.T) { TestingT(t) }

func (s *testSuite) Test_Simple(c *C) {
	c.Assert(nil, IsNil)
}
