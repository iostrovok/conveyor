package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	_ "github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"

	"github.com/iostrovok/conveyor/item"
	"github.com/iostrovok/conveyor/queues/stack"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {

	lastId := 100
	ctx, cancel := context.WithCancel(context.Background())
	st := stack.Init(lastId+10, ctx)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i <= lastId; i++ {
			st.ChanIn() <- item.NewItem(context.Background(), nil).SetID(int64(i))
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		}
	}()

	success, total := 0, 0
	go func() {
		last := -1
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case it, ok := <-st.ChanOut():
				if !ok {
					return
				}
				total++

				fmt.Printf("%d] id: %d\n", total, it.GetID())

				id := int(it.GetID())
				if last > -1 && id < last {
					success++
				}
				last = id

				time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
			}
		}
	}()

	<-time.After(3 * time.Second)
	cancel()

	wg.Wait()

	fmt.Printf("success: %d, total: %d\n", success, total)
}
