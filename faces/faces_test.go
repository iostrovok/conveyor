package faces_test

import (
	"testing"

	_ "github.com/golang/mock/mockgen/model"
	. "github.com/iostrovok/check"

	"github.com/iostrovok/conveyor/faces"
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestService(t *testing.T) { TestingT(t) }

func (s *testSuite) TestNeedToSkip(c *C) {
	h, err := faces.MakeEmptyHandler("")
	c.Assert(err, IsNil)
	c.Assert(h, NotNil)
}
