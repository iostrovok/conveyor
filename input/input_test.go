package input_test

import (
	"testing"

	. "github.com/iostrovok/check"

	_ "github.com/iostrovok/conveyor/input"
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestService(t *testing.T) { TestingT(t) }

// simple test for syntax.
func (s *testSuite) TestNil(c *C) {
	c.Assert(1, NotNil)
}
