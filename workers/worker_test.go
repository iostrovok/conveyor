package workers

import (
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"
	. "github.com/iostrovok/check"

	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/faces/mmock"
)

const (
	NameOne faces.Name = "NameOne"
	NameTwo faces.Name = "NameTwo"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestService(t *testing.T) { TestingT(t) }

func (s *testSuite) TestNeedToSkip(c *C) {

	//errCh := MockIChanEmptyMocker(c.T())
	//out := MockIChanEmptyMocker(c.T())
	//in := MockIChanEmptyMocker(c.T())
	//
	//worker := &Worker{
	//	name:  NameOne,
	//	in:    in,
	//	out:   out,
	//	errCh: errCh,
	//
	//	isLast:    false,
	//	typ:       faces.WorkerManagerType,
	//	isStarted: true,
	//}

	item := MockIItemMocker(c.T(), NameOne)
	c.Assert(needToSkip(faces.WorkerManagerType, NameOne, item), Equals, false)
	c.Assert(needToSkip(faces.FinalManagerType, NameOne, item), Equals, false)
	c.Assert(needToSkip(faces.ErrorManagerType, NameOne, item), Equals, false)

	item = MockIItemMocker(c.T(), NameTwo)
	c.Assert(needToSkip(faces.WorkerManagerType, NameOne, item), Equals, true)

	item = MockIItemMocker(c.T(), faces.EmptySkipName)
	c.Assert(needToSkip(faces.WorkerManagerType, NameOne, item), Equals, false)
	c.Assert(needToSkip(faces.FinalManagerType, NameOne, item), Equals, false)
	c.Assert(needToSkip(faces.ErrorManagerType, NameOne, item), Equals, false)
	c.Assert(needToSkip(faces.WorkerManagerType, NameOne, item), Equals, false)
}

func MockIItemMocker(t *testing.T, name faces.Name) *mmock.MockIItem {
	ctrl := gomock.NewController(t)
	m := mmock.NewMockIItem(ctrl)
	m.EXPECT().GetSkipToName().Return(name).AnyTimes()
	m.EXPECT().CleanSkipToName().AnyTimes()
	return m
}

// NewIPortfolio returns new mocker with mocked Portfolio table clients
func NewMockIChan(t *testing.T) *mmock.MockIChan {
	ctrl := gomock.NewController(t)
	return mmock.NewMockIChan(ctrl)
}

func MockIChanEmptyMocker(c *testing.T) *mmock.MockIChan {
	m := NewMockIChan(c)

	ch := make(chan int, 1)
	close(ch)

	return m
}
