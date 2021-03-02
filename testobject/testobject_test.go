package testobject

import (
	. "github.com/iostrovok/check"
	"testing"

	//"github.com/iostrovok/conveyor/faces/mmock"
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestService(t *testing.T) { TestingT(t) }

func (s *testSuite) Test_Simple(c *C) {
	c.Assert(nil, IsNil)
}
