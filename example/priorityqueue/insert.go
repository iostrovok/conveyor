package main

import (
	"context"
	"fmt"
	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/item"
	"math/rand"
	"strings"
	"time"

	_ "github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"

	pq "github.com/iostrovok/conveyor/queues/priorityqueue"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {

	ctx := context.Background()

	prs := []int{7, 15, 17, 19, 8, 19, 4, 15, 19, 3, 18, 18, 4, 3, 8, 13, 10, 11, 4, 5}

	array := make([]faces.IItem, 30)
	for i, p := range prs {
		a0 := item.NewItem(ctx, nil).SetPriority(p).SetID(int64(i))
		pq.InsertToBody(&array, a0, i)

		fmt.Printf("%s%s", PrintBody(array, i), "\n=======================\n")
	}

}

func PrintBody(body []faces.IItem, last int) string {
	out := make([]string, 0)
	for i := 0; i < last; i++ {
		item := body[i]
		out = append(out, fmt.Sprintf("%d] %d => %d", i, item.GetID(), item.GetPriority()))
	}

	return strings.Join(out, "\n")
}
