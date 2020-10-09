package stack

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	. "github.com/iostrovok/check"

	_ "github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"

	"github.com/iostrovok/conveyor/item"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestService(t *testing.T) { TestingT(t) }

func (s *testSuite) TestStepBystep(c *C) {

	lastId := 100
	st := Init(lastId+10, context.Background())
	for i := 0; i < lastId; i++ {
		st.ChanIn() <- item.New(context.Background(), nil).SetID(int64(i))
	}

	success, total := readTestData(st)

	c.Assert(total, Equals, lastId)
	c.Logf("success: %d, total: %d\n", success, total)
	c.Assert(float32(success) > 0.95*float32(total), Equals, true)
}

func (s *testSuite) TestInTheSameTime(c *C) {

	lastId := 100
	st := Init(lastId+10, context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < lastId; i++ {
			st.ChanIn() <- item.New(context.Background(), nil).SetID(int64(i))
			time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
		}

	}()

	success, total := 0, 0
	go func() {
		last := -1
		defer wg.Done()
		for {
			select {
			case it, ok := <-st.ChanOut():
				if !ok {
					return
				}
				total++

				id := int(it.GetID())
				if last == -1 {
					last = id
					continue
				}

				if id < last+2 {
					success++
				}

				last = id
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

			case <-time.After(5 * time.Second):
				return
			}
		}
	}()

	wg.Wait()

	c.Assert(total, Equals, lastId)
	c.Logf("success: %d, total: %d\n", success, total)
	c.Assert(float32(success) >= 0.90*float32(total), Equals, true)
}

func readTestData(st *Stack) (int, int) {
	success := 0
	last := -1
	total := 0
	for {
		select {
		case it, ok := <-st.ChanOut():
			if !ok {
				return success, total
			}
			total++

			id := int(it.GetID())
			if last == -1 {
				last = id
				continue
			}

			if id < last {
				success++
			}

			last = id

		case <-time.After(1 * time.Second):
			return success, total
		}
	}

	return success, total
}
