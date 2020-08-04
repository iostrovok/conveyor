package main



import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"

	"github.com/iostrovok/conveyor/item"
	pq "github.com/iostrovok/conveyor/queues/priorityqueue"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {

	lastId := 100
	ctx, cancel := context.WithCancel(context.Background())
	st := pq.Init(lastId+10, ctx)

	insertedList := make([]string, lastId)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < lastId; i++ {
			p := rand.Intn(20)
			insertedList[i] = strconv.Itoa(p)
			st.ChanIn() <- item.NewItem(context.Background(), nil).SetPriority(p).SetID(int64(i))
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		}
		fmt.Printf("All done\n")
		fmt.Printf("\n\n%s\n\n", st.Print())

		fmt.Printf("\n\n%s\n\n", strings.Join(insertedList, " - "))
	}()

	success, total := 0, 0
	go func() {
		//time.Sleep(5 * time.Second)
		lastPriority := -1
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-st.ChanOut():
				if !ok {
					return
				}
				total++

				priority := item.GetPriority()

				fmt.Printf("%d / %d] priority:	%d\n", total, item.GetID(), item.GetPriority())

				if lastPriority > -1 && priority <= lastPriority {
					success++
				}
				lastPriority = priority

				time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
			}
		}
	}()

	<-time.After(10 * time.Second)
	cancel()

	wg.Wait()

	fmt.Printf("success: %d, total: %d from %d\n", success, total, lastId)
}
