package workbench_test

import (
	"context"
	. "github.com/iostrovok/check"
	"github.com/iostrovok/conveyor/item"
	"github.com/iostrovok/conveyor/workbench"
	"sync"
	"testing"
)

const (
	lastID = 200
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestService(t *testing.T) { TestingT(t) }

func (s *testSuite) TestAddAndCount(c *C) {

	wb := workbench.New(2 * lastID)

	c.Assert(wb.Len(), Equals, 2*lastID)
	done := map[int]bool{}

	for i := 0; i < lastID; i++ {
		it := item.New(context.Background(), nil)
		it.SetPriority(100)
		done[wb.Add(it)] = true
	}

	//c.Logf("TestStepByStep: success: %d, total: %d\n", success, total)
	c.Assert(wb.Count(), Equals, lastID)
	c.Assert(len(done), Equals, lastID)

	for i := range done {
		c.Assert(wb.GetPriority(i), Equals, 100)
		wb.Clean(i)
	}

	c.Assert(wb.Count(), Equals, 0)

	// for empty value
	c.Assert(wb.GetPriority(lastID/2), Equals, 0)
}

func (s *testSuite) TestAddAndClean(c *C) {

	wb := workbench.New(lastID)

	c.Assert(wb.Len(), Equals, lastID)

	done := make(chan int, 1000)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10*lastID; i++ {
			done <- wb.Add(item.New(context.Background(), nil))
		}
		close(done)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case j, ok := <-done:
				if !ok {
					return
				}
				wb.Clean(j)
			}
		}
	}()

	wg.Wait()

	//c.Logf("TestStepByStep: success: %d, total: %d\n", success, total)
	c.Assert(wb.Count(), Equals, 0)
}
