package std_test

import (
	"testing"

	_ "github.com/golang/mock/mockgen/model"
	. "github.com/iostrovok/check"

	"github.com/iostrovok/conveyor/queues/std"
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestService(t *testing.T) { TestingT(t) }

func (s *testSuite) TestNeedToSkip(c *C) {
	c.Assert(std.New(0), NotNil)
}
