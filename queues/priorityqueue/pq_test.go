package priorityqueue_test

import (
	"context"
	"crypto/rand"
	"github.com/iostrovok/conveyor/faces"
	"math/big"
	"sync"
	"testing"
	"time"

	_ "github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"
	. "github.com/iostrovok/check"

	"github.com/iostrovok/conveyor/item"
	pq "github.com/iostrovok/conveyor/queues/priorityqueue"
	"github.com/iostrovok/conveyor/workbench"
)

const (
	lastID = 200
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestService(t *testing.T) { TestingT(t) }

func (s *testSuite) TestStepByStep(c *C) {

	wb := workbench.New(1000)
	st := pq.Init(context.Background(), wb, lastID+10)

	for i := 0; i < lastID; i++ {
		it := item.New(context.Background(), nil)
		it.SetID(int64(i))
		st.ChanIn() <- wb.Add(it)
	}

	success, total := readTestData(wb, st)

	c.Logf("TestStepByStep: success: %d, total: %d\n", success, total)
	c.Assert(total, Equals, lastID)
	c.Assert(float32(success) > 0.97*float32(total), Equals, true)
}

func (s *testSuite) TestInTheSameTime(c *C) {
	wb := workbench.New(1000)
	st := pq.Init(context.Background(), wb, lastID+10)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < lastID; i++ {
			it := item.New(context.Background(), nil)
			it.SetID(int64(i))
			st.ChanIn() <- wb.Add(it)
			k, _ := rand.Int(rand.Reader, big.NewInt(5))
			time.Sleep(time.Duration(k.Int64()) * time.Millisecond)
		}
	}()

	success, total := 0, 0
	go func() {
		last := -1
		defer wg.Done()
		for {
			select {
			case i, ok := <-st.ChanOut():
				if !ok {
					return
				}
				total++
				it, err := wb.Get(i)
				c.Assert(err, IsNil)

				id := int(it.GetID())
				if last == -1 {
					last = id

					continue
				}

				if id < last+2 {
					success++
				}

				last = id
				k, _ := rand.Int(rand.Reader, big.NewInt(100))
				time.Sleep(time.Duration(k.Int64()) * time.Millisecond)

			case <-time.After(5 * time.Second):
				// instead of cancel by time
				return
			}
		}
	}()

	wg.Wait()

	c.Logf("TestInTheSameTime success: %d, total: %d\n", success, total)
	c.Assert(total, Equals, lastID)
	c.Assert(float32(success) > 0.80*float32(total), Equals, true)
}

func readTestData(wb faces.IWorkBench, st *pq.PQ) (int, int) {
	success, total, lastPriority := 0, 0, 0
	for {
		select {
		case i, ok := <-st.ChanOut():
			if !ok {
				return success, total
			}
			total++

			it, err := wb.Get(i)
			if err != nil {
				continue
			}

			priority := it.GetPriority()

			if lastPriority > -1 && priority <= lastPriority {
				success++
			}
			lastPriority = priority

		case <-time.After(1 * time.Second):
			return success, total
		}
	}
}
