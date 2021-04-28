package queue

import "time"

type BlockingQueue interface {
	Queue

	/**
	 * Inserts the specified element into this queue, waiting if necessary
	 * for space to become available.
	 *
	 * @param e the element to add
	 * @throws InterruptedException if interrupted while waiting
	 * @throws NullPointerException if the specified element is null
	 * @throws IllegalArgumentException if some property of the specified
	 *         element prevents it from being added to this queue
	 */
	// 队列非满则插入, 队列满则等待.
	Put(i interface{}) error

	/**
	 * Inserts the specified element into this queue, waiting up to the
	 * specified wait time if necessary for space to become available.
	 *
	 * @param e the element to add
	 * @param timeout how long to wait before giving up, in units of
	 *        {@code unit}
	 * @param unit a {@code TimeUnit} determining how to interpret the
	 *        {@code timeout} parameter
	 * @return {@code true} if successful, or {@code false} if
	 *         the specified waiting time elapses before space is available
	 * @throws InterruptedException if interrupted while waiting
	 * @throws ClassCastException if the class of the specified element
	 *         prevents it from being added to this queue
	 * @throws NullPointerException if the specified element is null
	 * @throws IllegalArgumentException if some property of the specified
	 *         element prevents it from being added to this queue
	 */
	// 插入成功返回true, 插入失败返回false, 超时返回false
	OfferTimout(i interface{}, timeout time.Duration) bool

	/**
	 * Retrieves and removes the head of this queue, waiting if necessary
	 * until an element becomes available.
	 *
	 * @return the head of this queue
	 * @throws InterruptedException if interrupted while waiting
	 */
	// 队列非空则出列, 队列空则等待
	Take() interface{}

	PollTimeout(timeout time.Duration) interface{}

	/**
	 * Returns the number of additional elements that this queue can ideally
	 * (in the absence of memory or resource constraints) accept without
	 * blocking, or {@code Integer.MAX_VALUE} if there is no intrinsic
	 * limit.
	 *
	 * <p>Note that you <em>cannot</em> always tell if an attempt to insert
	 * an element will succeed by inspecting {@code remainingCapacity}
	 * because it may be the case that another thread is about to
	 * insert or remove an element.
	 *
	 * @return the remaining capacity
	 */
	RemainingCapacity() int
}


