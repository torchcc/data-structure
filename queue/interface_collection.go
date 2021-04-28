package queue

type Collection interface {
	// Query Operations
	/**
	 * Returns the number of elements in this collection.  If this collection
	 * contains more than <tt>Integer.MAX_VALUE</tt> elements, returns
	 * <tt>Integer.MAX_VALUE</tt>.
	 */
	Len() int // size()

	/**
	 * @return <tt>true</tt> if this collection contains no elements
	 */
	IsEmpty() bool

	/**
	 * Returns <tt>true</tt> if this collection contains the specified element.
	 * More formally, returns <tt>true</tt> if and only if this collection
	 * contains at least one element <tt>e</tt> such that
	 * <tt>(o==null&nbsp;?&nbsp;e==null&nbsp;:&nbsp;o.equals(e))</tt>.
	 *
	 * @param o element whose presence in this collection is to be tested
	 * @return <tt>true</tt> if this collection contains the specified
	 *         element
	 * @throws ClassCastException if the type of the specified element
	 *         is incompatible with this collection
	 *         (<a href="#optional-restrictions">optional</a>)
	 * @throws NullPointerException if the specified element is null and this
	 *         collection does not permit null elements
	 *         (<a href="#optional-restrictions">optional</a>)
	 */
	Contains(i interface{}) bool

	/**
	 * @Description: iterate through the collection
	 * @param f
	 */
	Range(f func(value interface{}) bool)

	/**
	 * Returns an array containing all of the elements in this collection.
	 * If this collection makes any guarantees as to what order its elements
	 * are returned by its iterator, this method must return the elements in
	 * the same order.
	 *
	 * <p>The returned array will be "safe" in that no references to it are
	 * maintained by this collection.  (In other words, this method must
	 * allocate a new array even if this collection is backed by an array).
	 * The caller is thus free to modify the returned array.
	 *
	 * <p>This method acts as bridge between array-based and collection-based
	 * APIs.
	 *
	 * @return an array containing all of the elements in this collection
	 */
	ToSlice() []interface{}

	// Modification Operations

	/**
	 * Ensures that this collection contains the specified element (optional
	 * operation).  Returns <tt>true</tt> if this collection changed as a
	 * result of the call.  (Returns <tt>false</tt> if this collection does
	 * not permit duplicates and already contains the specified element.)<p>
	 */
	Add(i interface{}) bool

	/**
	 * Removes a single instance of the specified element from this
	 * collection, if it is present (optional operation).  More formally,
	 * removes an element <tt>e</tt> such that
	 * <tt>(o==null&nbsp;?&nbsp;e==null&nbsp;:&nbsp;o.equals(e))</tt>, if
	 * this collection contains one or more such elements.  Returns
	 * <tt>true</tt> if this collection contained the specified element (or
	 * equivalently, if this collection changed as a result of the call).
	 */
	Remove(i interface{}) bool

	// Bulk Operations

	/**
	 * Returns <tt>true</tt> if this collection contains all of the elements
	 * in the specified collection.
	 *
	 * @param  c collection to be checked for containment in this collection
	 * @return <tt>true</tt> if this collection contains all of the elements
	 *         in the specified collection
	 * @throws ClassCastException if the types of one or more elements
	 *         in the specified collection are incompatible with this
	 *         collection
	 *         (<a href="#optional-restrictions">optional</a>)
	 * @throws NullPointerException if the specified collection contains one
	 *         or more null elements and this collection does not permit null
	 *         elements
	 *         (<a href="#optional-restrictions">optional</a>),
	 *         or if the specified collection is null.
	 * @see    #contains(Object)
	 */
	ContainsAll(c Collection) bool

	/**
	 * Adds all of the elements in the specified collection to this collection
	 * (optional operation).  The behavior of this operation is undefined if
	 * the specified collection is modified while the operation is in progress.
	 * (This implies that the behavior of this call is undefined if the
	 * specified collection is this collection, and this collection is
	 * nonempty.)
	 *
	 * @param c collection containing elements to be added to this collection
	 * @return <tt>true</tt> if this collection changed as a result of the call
	 * @throws UnsupportedOperationException if the <tt>addAll</tt> operation
	 *         is not supported by this collection
	 * @throws ClassCastException if the class of an element of the specified
	 *         collection prevents it from being added to this collection
	 * @throws NullPointerException if the specified collection contains a
	 *         null element and this collection does not permit null elements,
	 *         or if the specified collection is null
	 * @throws IllegalArgumentException if some property of an element of the
	 *         specified collection prevents it from being added to this
	 *         collection
	 * @throws IllegalStateException if not all the elements can be added at
	 *         this time due to insertion restrictions
	 * @see #add(Object)
	 */
	AddAll(c Collection) bool

	/**
	 * Removes all of this collection's elements that are also contained in the
	 * specified collection (optional operation).  After this call returns,
	 * this collection will contain no elements in common with the specified
	 * collection.
	 *
	 * @param c collection containing elements to be removed from this collection
	 * @return <tt>true</tt> if this collection changed as a result of the
	 *         call
	 * @throws UnsupportedOperationException if the <tt>removeAll</tt> method
	 *         is not supported by this collection
	 * @throws ClassCastException if the types of one or more elements
	 *         in this collection are incompatible with the specified
	 *         collection
	 *         (<a href="#optional-restrictions">optional</a>)
	 * @throws NullPointerException if this collection contains one or more
	 *         null elements and the specified collection does not support
	 *         null elements
	 *         (<a href="#optional-restrictions">optional</a>),
	 *         or if the specified collection is null
	 * @see #remove(Object)
	 * @see #contains(Object)
	 */
	RemoveAll(c Collection) bool

	/**
     * Removes all of the elements of this collection that satisfy the given
     * predicate.  Errors or runtime exceptions thrown during iteration or by
     * the predicate are relayed to the caller.
     *
	 */
	RemoveIf(filter func(value interface{}) bool) bool

	/**
	 * Retains only the elements in this collection that are contained in the
	 * specified collection (optional operation).  In other words, removes from
	 * this collection all of its elements that are not contained in the
	 * specified collection.
	 *
	 * @param c collection containing elements to be retained in this collection
	 * @return <tt>true</tt> if this collection changed as a result of the call
	 * @throws UnsupportedOperationException if the <tt>retainAll</tt> operation
	 *         is not supported by this collection
	 * @throws ClassCastException if the types of one or more elements
	 *         in this collection are incompatible with the specified
	 *         collection
	 *         (<a href="#optional-restrictions">optional</a>)
	 * @throws NullPointerException if this collection contains one or more
	 *         null elements and the specified collection does not permit null
	 *         elements
	 *         (<a href="#optional-restrictions">optional</a>),
	 *         or if the specified collection is null
	 * @see #remove(Object)
	 * @see #contains(Object)
	 */
	RetainAll(c Collection) bool

	/**
	 * Removes all of the elements from this collection (optional operation).
	 * The collection will be empty after this method returns.
	 *
	 * @throws UnsupportedOperationException if the <tt>clear</tt> operation
	 *         is not supported by this collection
	 */
	Clear()
}
