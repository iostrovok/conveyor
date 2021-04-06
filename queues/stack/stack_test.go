package stack_test

// 	Test package for stack.

import (
	"context"
	"crypto/rand"
	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/workbench"
	"math/big"
	"sync"
	"testing"
	"time"

	. "github.com/iostrovok/check"

	_ "github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"

	"github.com/iostrovok/conveyor/item"
	"github.com/iostrovok/conveyor/queues/stack"
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestService(t *testing.T) { TestingT(t) }

const (
	lastID = 100
)

func (s *testSuite) TestStepBystep(c *C) {
	wb := workbench.New(1000)
	st := stack.Init(context.Background(), lastID+10)
	for i := 0; i < lastID; i++ {
		it := item.New(context.Background(), nil)
		it.SetID(int64(i))
		st.ChanIn() <- wb.Add(it)
	}

	success, total := readTestData(wb, st)

	c.Assert(total, Equals, lastID)
	c.Logf("success: %d, total: %d\n", success, total)
	c.Assert(float32(success) > 0.95*float32(total), Equals, true)
}

func (s *testSuite) TestInTheSameTime(c *C) {
	st := stack.Init(context.Background(), lastID+10)
	wb := workbench.New(1000)
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
				if err != nil {
					continue
				}

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
				return
			}
		}
	}()

	wg.Wait()

	c.Assert(total, Equals, lastID)
	c.Logf("success: %d, total: %d\n", success, total)
	c.Assert(float32(success) >= 0.90*float32(total), Equals, true)
}

func readTestData(wb faces.IWorkBench, st *stack.Stack) (int, int) {
	success := 0
	last := -1
	total := 0
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
}
