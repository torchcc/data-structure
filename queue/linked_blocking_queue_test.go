package queue

import (
	"fmt"
	"testing"
	"time"
)

// we use only Offer and PollTimeout
// Take Poll Offer Put, PollTimeout  PutTimeout
func TestLinkedBlockingQueue_Put(t *testing.T) {
	queue := NewLinkedBlockingQueue(10)

	go func() {
		time.Sleep(10 * time.Second)
		for i := 0; i < 10; i++ {
			fmt.Printf("putting %d", i)
			 queue.OfferTimout(i, time.Second * 5)
		}
		queue.Range(func(value interface{}) bool {
			i := value.(int)
			fmt.Println(i)
			if i == 8 {
				return false
			}
			return true
		})
		time.Sleep(8 * time.Minute)
	}()

	for {
		fmt.Println("begin ....")
		x := queue.PollTimeout(2 * time.Second)
		fmt.Println("x: ", x)
		// time.Sleep(100 * time.Minute)


	}
}
