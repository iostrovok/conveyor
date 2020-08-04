package priorityqueue

import (
	"context"
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"

	_ "github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"
	. "github.com/iostrovok/check"

	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/item"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestService(t *testing.T) { TestingT(t) }

func (s *testSuite) TestFindPosition(c *C) {

	ctx := context.Background()

	a := make([]faces.IItem, 10)
	a[0] = item.NewItem(ctx, nil).SetPriority(10)
	a[1] = item.NewItem(ctx, nil).SetPriority(15)
	a[2] = item.NewItem(ctx, nil).SetPriority(16)
	a[3] = item.NewItem(ctx, nil).SetPriority(17)
	a[4] = item.NewItem(ctx, nil).SetPriority(20)

	c.Assert(findPosition(a, 10, 5), Equals, 0)
	c.Assert(findPosition(a, 11, 5), Equals, 1)
	c.Assert(findPosition(a, 15, 5), Equals, 1)
	c.Assert(findPosition(a, 16, 5), Equals, 2)
	c.Assert(findPosition(a, 19, 5), Equals, 4)
	c.Assert(findPosition(a, 20, 5), Equals, 5)
	c.Assert(findPosition(a, 100, 5), Equals, 5)
}

func (s *testSuite) TestInsertToBody(c *C) {

	ctx := context.Background()

	array := make([]faces.IItem, 10)
	a0 := item.NewItem(ctx, nil).SetPriority(10)
	a1 := item.NewItem(ctx, nil).SetPriority(13)
	a2 := item.NewItem(ctx, nil).SetPriority(15)
	a3 := item.NewItem(ctx, nil).SetPriority(15)
	a4 := item.NewItem(ctx, nil).SetPriority(13)
	a5 := item.NewItem(ctx, nil).SetPriority(12)

	InsertToBody(&array, a0, 0)
	InsertToBody(&array, a1, 1)
	InsertToBody(&array, a2, 2)
	InsertToBody(&array, a3, 3)
	InsertToBody(&array, a4, 4)
	InsertToBody(&array, a5, 5)

	c.Assert(array[0].GetPriority(), Equals, 10)
	c.Assert(array[1].GetPriority(), Equals, 12)
	c.Assert(array[2].GetPriority(), Equals, 13)
	c.Assert(array[3].GetPriority(), Equals, 13)
	c.Assert(array[4].GetPriority(), Equals, 15)
	c.Assert(array[5].GetPriority(), Equals, 15)
}

func (s *testSuite) TestInsertToBody2(c *C) {

	ctx := context.Background()

	prs := []int{7, 15, 17, 19, 8, 19, 4, 15, 19, 3, 18, 18, 4, 3, 8, 13, 10, 11, 4, 5}

	c.Logf("%+v\n", prs)

	array := make([]faces.IItem, 100)
	for i, p := range prs {
		a0 := item.NewItem(ctx, nil).SetPriority(p)
		InsertToBody(&array, a0, i)
	}

	sort.Ints(prs)
	c.Logf("%+v\n", prs)
	for i, p := range prs {
		c.Logf("%d] %d == %d\n", i, array[i].GetPriority(), p)
		c.Assert(array[i].GetPriority(), Equals, p)
	}
}

func (s *testSuite) TestInsertToBodyLong(c *C) {
	ctx := context.Background()

	count := 1000
	prs := make([]int, count)
	array := make([]faces.IItem, count)
	for i := 0; i < count; i++ {
		p := rand.Intn(20)
		prs[i] = p
		a0 := item.NewItem(ctx, nil).SetPriority(p)
		InsertToBody(&array, a0, i)
	}

	sort.Ints(prs)
	c.Logf("%+v\n", prs)
	for i, p := range prs {
		c.Logf("%d] %d == %d\n", i, array[i].GetPriority(), p)
		c.Assert(array[i].GetPriority(), Equals, p)
	}

}

func (s *testSuite) TestStepBystep(c *C) {

	lastId := 200
	st := Init(lastId+10, context.Background())
	for i := 0; i < lastId; i++ {
		st.ChanIn() <- item.NewItem(context.Background(), nil).SetID(int64(i))
	}

	success, total := readTestData(st)

	c.Logf("TestStepBystepTestStepBystep: success: %d, total: %d\n", success, total)
	c.Assert(total, Equals, lastId)
	c.Assert(float32(success) > 0.97*float32(total), Equals, true)
}

func (s *testSuite) TestInTheSameTime(c *C) {

	lastId := 200
	st := Init(lastId+10, context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < lastId; i++ {
			st.ChanIn() <- item.NewItem(context.Background(), nil).SetID(int64(i))
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
				// instead of cancel by time
				return
			}
		}
	}()

	wg.Wait()

	c.Logf("TestInTheSameTime success: %d, total: %d\n", success, total)
	c.Assert(total, Equals, lastId)
	c.Assert(float32(success) > 0.80*float32(total), Equals, true)
}

func readTestData(st *PQ) (int, int) {
	success, total, lastPriority := 0, 0, 0
	for {
		select {
		case item, ok := <-st.ChanOut():
			if !ok {
				return success, total
			}
			total++

			priority := item.GetPriority()

			if lastPriority > -1 && priority <= lastPriority {
				success++
			}
			lastPriority = priority


		case <-time.After(1 * time.Second):
			return success, total
		}
	}

	return success, total
}
