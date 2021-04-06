package priorityqueue

import (
	"testing"

	"github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"
	. "github.com/iostrovok/check"

	"github.com/iostrovok/conveyor/faces/mmock"
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestService(t *testing.T) { TestingT(t) }

func (s *testSuite) TestFindPosition(c *C) {
	ctrl := gomock.NewController(c.T())
	defer ctrl.Finish()

	wb := mmock.NewMockIWorkBench(ctrl)
	wb.EXPECT().GetPriority(0).AnyTimes().Return(10)
	wb.EXPECT().GetPriority(1).AnyTimes().Return(15)
	wb.EXPECT().GetPriority(2).AnyTimes().Return(16)
	wb.EXPECT().GetPriority(3).AnyTimes().Return(17)
	wb.EXPECT().GetPriority(4).AnyTimes().Return(20)

	c.Assert(findPosition(wb, 10, 5), Equals, 0)
	c.Assert(findPosition(wb, 11, 5), Equals, 1)
	c.Assert(findPosition(wb, 15, 5), Equals, 1)
	c.Assert(findPosition(wb, 16, 5), Equals, 2)
	c.Assert(findPosition(wb, 19, 5), Equals, 4)
	c.Assert(findPosition(wb, 20, 5), Equals, 5)
	c.Assert(findPosition(wb, 100, 5), Equals, 5)
}
