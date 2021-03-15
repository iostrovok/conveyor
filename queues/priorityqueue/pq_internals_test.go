package priorityqueue

import (
	"context"
	"testing"

	_ "github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"
	. "github.com/iostrovok/check"

	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/item"
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestService(t *testing.T) { TestingT(t) }

func (s *testSuite) TestFindPosition(c *C) {
	ctx := context.Background()

	a := make([]faces.IItem, 10)
	a[0] = item.New(ctx, nil)
	a[0].SetPriority(10)
	a[1] = item.New(ctx, nil)
	a[1].SetPriority(15)
	a[2] = item.New(ctx, nil)
	a[2].SetPriority(16)
	a[3] = item.New(ctx, nil)
	a[3].SetPriority(17)
	a[4] = item.New(ctx, nil)
	a[4].SetPriority(20)

	c.Assert(findPosition(a, 10, 5), Equals, 0)
	c.Assert(findPosition(a, 11, 5), Equals, 1)
	c.Assert(findPosition(a, 15, 5), Equals, 1)
	c.Assert(findPosition(a, 16, 5), Equals, 2)
	c.Assert(findPosition(a, 19, 5), Equals, 4)
	c.Assert(findPosition(a, 20, 5), Equals, 5)
	c.Assert(findPosition(a, 100, 5), Equals, 5)
}
