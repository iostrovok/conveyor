package item_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/iostrovok/check"
	"github.com/pkg/errors"

	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/faces/mmock"
	"github.com/iostrovok/conveyor/item"
)

const (
	NameOne faces.Name = "NameOne"
	NameTwo faces.Name = "NameTwo"
)

var err = errors.New(item.LastHandlerErrorNote)

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
	it := item.New(context.Background(), nil)
	it.SetSkipToName(NameOne)

	checkResult(c, it, MockIWorkerMocker(c.T(), NameOne, faces.WorkerManagerType, false), false)
	checkResult(c, it, MockIWorkerMocker(c.T(), NameOne, faces.FinalManagerType, false), false)
	checkResult(c, it, MockIWorkerMocker(c.T(), NameOne, faces.ErrorManagerType, false), false)

	it.SetSkipToName(NameTwo)
	checkResult(c, it, MockIWorkerMocker(c.T(), NameOne, faces.WorkerManagerType, false), true)

	it.SetSkipToName(faces.EmptySkipName)
	checkResult(c, it, MockIWorkerMocker(c.T(), NameOne, faces.WorkerManagerType, false), false)
	checkResult(c, it, MockIWorkerMocker(c.T(), NameOne, faces.FinalManagerType, false), false)
	checkResult(c, it, MockIWorkerMocker(c.T(), NameOne, faces.ErrorManagerType, false), false)
	checkResult(c, it, MockIWorkerMocker(c.T(), NameOne, faces.WorkerManagerType, false), false)

	checkResult(c, it, MockIWorkerMocker(c.T(), NameOne, faces.WorkerManagerType, true), false)
}

func (s *testSuite) TestNillContex(c *C) {
	it := item.New(context.Background(), nil)
	c.Assert(it, NotNil)
}

func checkResult(c *C, item faces.IItem, w faces.IWorker, needSkip bool) {
	skip, err := item.NeedToSkip(w)

	c.Assert(err, IsNil)
	c.Assert(skip, Equals, needSkip)
}
