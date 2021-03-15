package workers_test

import (
	"testing"

	_ "github.com/golang/mock/mockgen/model"
	. "github.com/iostrovok/check"
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestService(t *testing.T) { TestingT(t) }

func (s *testSuite) TestNeedToSkip(c *C) {
	c.Assert(1, Equals, 1)
}
