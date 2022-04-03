package heap

import (
	"time"
)

// Photo holds the attributes to a photo's geological location, timestamp, and the addresses of the photo.
//
// An empty struct is used as the value on the hashmap, this is because an empty struct allocates no memory.
// The use of the hashmap is to remove any duplicate addresses
type Photo struct {
	Timestamp time.Time
	Latitude  float64
	Longitude float64
	Addresses map[string]struct{}
}

// An PhotoHeap is a min-heap of photos. PhotoHeap implements sort.Interface so that the heap can be ordered,
// so that heap.Heap can become a min-heap BST.
type PhotoHeap []Photo

func (d PhotoHeap) Len() int           { return len(d) }
func (d PhotoHeap) Less(i, j int) bool { return d[i].Timestamp.Unix() < d[j].Timestamp.Unix() }
func (d PhotoHeap) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

// Push is used to append a Photo to the end of PhotoHeap.
// As Push is called from heap.Push we can be certain that the type will always be Photo.
func (d *PhotoHeap) Push(x interface{}) {
	*d = append(*d, x.(Photo))
}

// Pop is used to pop the top of the stack. heap.Pop is first called before this method,
// PhotoHeap.Swap is then called to swap index 0 with len(n-1).
func (d *PhotoHeap) Pop() interface{} {
	old := *d
	n := len(old)
	x := old[n-1]
	*d = old[0 : n-1]
	return x
}
