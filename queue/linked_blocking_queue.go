package queue

import (
	"container/list"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/torchcc/data-structure/error"
)

type LinkedBlockingQueue struct {
	// The number of items in the Queue
	length int64

	// the capacity set
	capacity int

	// Lock held by take, poll, etc
	takeLock *sync.Mutex
	// Condition for waiting reads
	// Wait queue for waiting takes
	notEmpty *sync.Cond

	// Lock held by put, offer, etc
	putLock *sync.Mutex
	// Wait queue for waiting puts
	notFull *sync.Cond

	// head of linked list
	head *list.List
}

func (q *LinkedBlockingQueue) Offer(i interface{}) bool {
	if i == nil {
		panic(NilPointerError)
	}
	if q.capacity == q.Len() {
		return false
	}
	c := -1
	q.putLock.Lock()
	defer q.putLock.Unlock()
	if q.Len() < q.capacity {
		q.head.PushBack(i)
		c = q.Len()
		atomic.AddInt64(&q.length, 1)
		if c+1 < q.capacity {
			q.notFull.Signal()
		}
	}
	if c == 0 {
		q.signalNotEmpty()
	}
	return c >= 0
}

func (q *LinkedBlockingQueue) RemoveHead() (x interface{}) {
	x = q.Poll()
	if x != nil {
		return x
	}
	panic(NoSuchElementError)
}

// a helper function, we do not lock helper function
func (q *LinkedBlockingQueue) dequeue() interface{} {
	return q.head.Remove(q.head.Front())
}

/**
 * Retrieves, but does not remove, the head of this queue.  This method
 * differs from {@link #peek peek} only in that it throws an exception if
 * this queue is empty.
 *
 * <p>This implementation returns the result of <tt>peek</tt>
 * unless the queue is empty.
 *
 * @return the head of this queue
 * @throws NoSuchElementException if this queue is empty
 */
func (q *LinkedBlockingQueue) Element() interface{} {
	if x := q.Peek(); x != nil {
		return x
	} else {
		panic(NoSuchElementError)
	}

}

func (q *LinkedBlockingQueue) Peek() interface{} {
	if q.Len() == 0 {
		return nil
	}
	q.takeLock.Lock()
	defer q.takeLock.Unlock()
	if first := q.head.Front(); first == nil {
		return nil
	} else {
		return first.Value
	}
}

/**
 * Inserts the specified element at the tail of this queue, waiting if
 * necessary for space to become available.
 *
 * @throws InterruptedException {@inheritDoc}
 * @throws NullPointerException {@inheritDoc}
 */
func (q *LinkedBlockingQueue) Put(i interface{}) error {
	if i == nil {
		return NilPointerError
	}
	c := -1
	q.putLock.Lock()
	for q.Len() == q.capacity {
		q.notFull.Wait()
	}
	q.head.PushBack(i)
	c = q.Len()
	atomic.AddInt64(&q.length, 1)
	if c+1 < q.capacity {
		q.notFull.Signal()
	}
	q.putLock.Unlock()
	if c == 0 {
		q.signalNotEmpty()
	}
	return nil
}

func (q *LinkedBlockingQueue) OfferTimout(i interface{}, timeout time.Duration) bool {
	if i == nil {
		panic(NilPointerError)
	}
	c := -1
	begin := time.Now()

LOOP:
	q.putLock.Lock()
	if q.Len() < q.capacity {
		q.head.PushBack(i)
		c = q.Len()
		atomic.AddInt64(&q.length, 1)
		if c+1 < q.capacity {
			q.notFull.Signal()
		}
		if c == 0 {
			q.signalNotEmpty()
		}
		q.putLock.Unlock()
		return true
	}
	q.putLock.Unlock()

	select {
	case <-time.After(time.Microsecond * 5):
		if time.Now().Sub(begin) > timeout {
			return false
		}
		goto LOOP
	}
}

func (q *LinkedBlockingQueue) Take() interface{} {
	c := -1
	var x interface{}
	q.takeLock.Lock()
	for q.Len() == 0 {
		q.notEmpty.Wait()
	}
	x = q.dequeue()
	c = q.Len()
	atomic.AddInt64(&q.length, -1)
	if c > 1 {
		q.notEmpty.Signal()
	}
	q.takeLock.Unlock()
	if c == q.capacity {
		q.signalNotFull()
	}
	return x
}

func (q *LinkedBlockingQueue) Poll() (x interface{}) {
	if q.Len() == 0 {
		return nil
	}
	c := -1
	q.takeLock.Lock()
	defer q.takeLock.Unlock()
	if q.Len() > 0 {
		x = q.dequeue()
		c = q.Len()
		atomic.AddInt64(&q.length, -1)
		if c > 1 {
			q.notEmpty.Signal()
		}
	}
	if c == q.capacity {
		q.signalNotFull()
	}
	return x
}

func (q *LinkedBlockingQueue) PollTimeout(timeout time.Duration) (x interface{}) {
	begin := time.Now()
LOOP:
	q.takeLock.Lock()
	if q.Len() > 0 {
		c := -1
		if q.Len() > 0 {
			x = q.dequeue()
			c = q.Len()
			atomic.AddInt64(&q.length, -1)
			if c > 1 {
				q.notEmpty.Signal()
			}
		}
		if c == q.capacity {
			q.signalNotFull()
		}
		q.takeLock.Unlock()
		return
	}
	q.takeLock.Unlock()

	select {
	case <-time.After(time.Microsecond * 5):
		if time.Now().Sub(begin) > timeout {
			return
		}
		goto LOOP
	}
}

func (q *LinkedBlockingQueue) RemainingCapacity() int {
	return q.capacity - q.Len()
}

func (q *LinkedBlockingQueue) Len() int {
	return int(atomic.LoadInt64(&q.length))
}

func (q *LinkedBlockingQueue) IsEmpty() bool {
	return q.Len() == 0
}

func (q *LinkedBlockingQueue) Contains(i interface{}) bool {
	if i == nil {
		return false
	}
	q.fullyLock()
	defer q.fullyUnlock()
	for cur := q.head.Front(); cur != nil; cur = cur.Next() {
		if cur.Value == i {
			return true
		}
	}
	return false
}

func (q *LinkedBlockingQueue) Range(f func(value interface{}) bool) {
	q.fullyLock()
	defer q.fullyUnlock()
	for cur := q.head.Front(); cur != nil; cur = cur.Next() {
		if !f(cur.Value) {
			return
		}
	}
}

func (q *LinkedBlockingQueue) ToSlice() []interface{} {
	q.fullyLock()
	defer q.fullyUnlock()
	ret := make([]interface{}, 0, q.Len())
	for cur := q.head.Front(); cur != nil; cur = cur.Next() {
		ret = append(ret, cur.Value)
	}
	return ret
}

func (q *LinkedBlockingQueue) ToString() string {
	q.fullyLock()
	defer q.fullyUnlock()
	if p := q.head.Front(); p == nil {
		return "[]"
	} else {
		sb := "["
		for {
			e := p.Value
			if e == q {
				sb += "(this Collection)"
			} else {
				sb += fmt.Sprintf("%v", e)
			}
			p = p.Next()
			if p == nil {
				return sb + "]"
			}
			sb += ", "
		}
	}
}

func (q *LinkedBlockingQueue) Add(i interface{}) bool {
	if q.Offer(i) {
		return true
	}
	panic(IllegalStateError)
}

/**
 * Removes a single instance of the specified element from this queue,
 * if it is present.  More formally, removes an element {@code e} such
 * that {@code o.equals(e)}, if this queue contains one or more such
 * elements.
 * Returns {@code true} if this queue contained the specified element
 * (or equivalently, if this queue changed as a result of the call).
 *
 * @param o element to be removed from this queue, if present
 * @return {@code true} if this queue changed as a result of the call
 */
func (q *LinkedBlockingQueue) Remove(i interface{}) bool {
	if i == nil {
		return false
	}
	q.fullyLock()
	defer q.fullyUnlock()
	for cur := q.head.Front(); cur != nil; cur = cur.Next() {
		if cur.Value == i {
			q.head.Remove(cur)
			atomic.AddInt64(&q.length, -1)
			return true
		}
	}
	return false
}

func (q *LinkedBlockingQueue) ContainsAll(c Collection) bool {
	containsAll := true
	c.Range(func(value interface{}) bool {
		if !q.Contains(value) {
			containsAll = false
			return false
		}
		return true
	})
	return containsAll
}

func (q *LinkedBlockingQueue) AddAll(c Collection) bool {
	if c == nil {
		panic(NilPointerError)
	}
	if toAdd, ok := c.(*LinkedBlockingQueue); ok && toAdd == q {
		panic(IllegalArgumentError)
	}
	modified := false
	c.Range(func(value interface{}) bool {
		if q.Add(value) {
			modified = true
		}
		return true
	})
	return modified
}

func (q *LinkedBlockingQueue) RemoveAll(c Collection) bool {
	panic("implement me")
}

func (q *LinkedBlockingQueue) RemoveIf(filter func(value interface{}) bool) bool {
	panic("implement me")
}

func (q *LinkedBlockingQueue) RetainAll(c Collection) bool {
	panic("implement me")
}

/**
 * Atomically removes all of the elements from this queue.
 * The queue will be empty after this call returns.
 */
func (q *LinkedBlockingQueue) Clear() {
	panic("implement me")
	// q.fullyLock()
	// defer q.fullyUnlock()
}

func NewLinkedBlockingQueue(capacity int) *LinkedBlockingQueue {
	if capacity == 0 {
		capacity = math.MaxInt32
	}
	putLock := new(sync.Mutex)
	takeLock := new(sync.Mutex)
	return &LinkedBlockingQueue{
		capacity: capacity,
		takeLock: takeLock,
		notEmpty: sync.NewCond(takeLock),
		putLock:  putLock,
		notFull:  sync.NewCond(putLock),
		head:     list.New(),
	}
}

func FromSlice(s []interface{}) (*LinkedBlockingQueue, error) {
	q := NewLinkedBlockingQueue(0)
	q.putLock.Lock()
	defer q.putLock.Unlock()
	var n int64 = 0
	for _, item := range s {
		if item == nil {
			return nil, NilPointerError
		}
		if n == int64(q.capacity) {
			return nil, FullError
		}
		q.head.PushBack(item)
		n++
	}
	atomic.StoreInt64(&q.length, n)
	return q, nil
}

/**
 * Signals a waiting take. Called only from put/offer (which do not
 * otherwise ordinarily lock takeLock.)
 */
func (q *LinkedBlockingQueue) signalNotEmpty() {
	q.takeLock.Lock()
	defer q.takeLock.Unlock()
	q.notEmpty.Signal()
}

/**
 * Signals a waiting put. Called only from take/poll.
 */
func (q *LinkedBlockingQueue) signalNotFull() {
	q.putLock.Lock()
	defer q.putLock.Unlock()
	q.notFull.Signal()
}

/**
 * Locks to prevent both puts and takes.
 */
func (q *LinkedBlockingQueue) fullyLock() {
	q.takeLock.Lock()
	q.putLock.Lock()
}

/**
 * Unlocks to allow both puts and takes.
 */
func (q *LinkedBlockingQueue) fullyUnlock() {
	q.takeLock.Unlock()
	q.putLock.Unlock()
}
