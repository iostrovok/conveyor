package conveyor_test

import (
	_ "github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"
	. "github.com/iostrovok/check"
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func (s *testSuite) Testsimple1(c *C) {
	c.Assert(1, DeepEquals, 1)
}
