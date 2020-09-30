package item

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/iostrovok/conveyor/faces/mmock"
	"math/rand"
	"testing"
	"time"

	. "github.com/iostrovok/check"

	"github.com/iostrovok/conveyor/faces"
	//"github.com/iostrovok/conveyor/faces/mmock"
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

func MockIWorkerMocker(t *testing.T, name faces.Name, typ faces.ManagerType, isLast bool) *mmock.MockIWorker {
	ctrl := gomock.NewController(t)
	m := mmock.NewMockIWorker(ctrl)
	m.EXPECT().GetBorderCond().Return(NameOne, typ, isLast).AnyTimes()
	return m
}

func (s *testSuite) TestNeedToSkip(c *C) {

	item := New(context.Background(), nil)
	item.SetSkipToName(NameOne)

	checkResult(c, item, MockIWorkerMocker(c.T(), NameOne, faces.WorkerManagerType, false), false, nil)
	checkResult(c, item, MockIWorkerMocker(c.T(), NameOne, faces.FinalManagerType, false), false, nil)
	checkResult(c, item, MockIWorkerMocker(c.T(), NameOne, faces.ErrorManagerType, false), false, nil)

	item.SetSkipToName(NameTwo)
	checkResult(c, item, MockIWorkerMocker(c.T(), NameOne, faces.WorkerManagerType, false), true, nil)

	item.SetSkipToName(faces.EmptySkipName)
	checkResult(c, item, MockIWorkerMocker(c.T(), NameOne, faces.WorkerManagerType, false), false, nil)
	checkResult(c, item, MockIWorkerMocker(c.T(), NameOne, faces.FinalManagerType, false), false, nil)
	checkResult(c, item, MockIWorkerMocker(c.T(), NameOne, faces.ErrorManagerType, false), false, nil)
	checkResult(c, item, MockIWorkerMocker(c.T(), NameOne, faces.WorkerManagerType, false), false, nil)
}

func (s *testSuite) TestNillContex(c *C) {
	item := New(nil, nil)
	c.Assert(item, NotNil)
}

func checkResult(c *C, item faces.IItem, w faces.IWorker, needSkip bool, errIn error) {
	skip, err := item.NeedToSkip(w)

	if errIn == nil {
		c.Assert(err, IsNil)
		c.Assert(skip, Equals, needSkip)
	} else {
		c.Assert(err, NotNil)
	}

}
