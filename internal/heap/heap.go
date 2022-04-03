package heap

import (
	"container/heap"
	"errors"
	"sync"
)

var (
	// ErrEmptyHeap is the error returned by Pop when there are no more results on the heap.
	// This doesn't necessarily mean that the heap will always be empty, due to the nature of a BST,
	// additional data can be added at a later date.
	ErrEmptyHeap = errors.New("empty heap")
)

// Heap is an implementation of a BST by using heap.Interface.
type Heap struct {
	photoHeap *PhotoHeap
	mu        *sync.RWMutex
}

// New initialises the heap and returns
func New() *Heap {
	photoHeap := new(PhotoHeap)

	heap.Init(photoHeap)

	return &Heap{photoHeap: photoHeap, mu: &sync.RWMutex{}}
}

// Push is used to push a Photo onto the heap. Push can be safely called concurrently due to the use of a sync.RWMutex.
func (h *Heap) Push(photo Photo) {
	h.mu.Lock()
	defer h.mu.Unlock()
	heap.Push(h.photoHeap, photo)
}

// Pop is used to pop the first item from the heap and returns. As heap.
// Pop has no internal checks on if the heap is empty or not, the check is done within this function.
// If the length is 0 we return ErrEmptyHeap.
//
// As all items that are pushed onto the queue are strongly typed, type casting from v to Photo is guaranteed to work.
func (h *Heap) Pop() (Photo, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.photoHeap.Len() == 0 {
		return Photo{}, ErrEmptyHeap
	}

	v := heap.Pop(h.photoHeap)
	photo := v.(Photo)

	return photo, nil
}
