package heap

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHeap_Push(t *testing.T) {
	tests := []struct {
		name     string
		photos   []Photo
		expected *PhotoHeap
	}{
		{
			name: "pushes to the heap and exists",
			photos: []Photo{
				{
					Timestamp: time.Date(2022, 01, 02, 10, 11, 12, 0, time.UTC),
					Latitude:  -45,
					Longitude: 10,
					Addresses: map[string]struct{}{
						"London":         {},
						"United Kingdom": {},
					},
				},
			},
			expected: &PhotoHeap{
				{
					Timestamp: time.Date(2022, 01, 02, 10, 11, 12, 0, time.UTC),
					Latitude:  -45,
					Longitude: 10,
					Addresses: map[string]struct{}{
						"London":         {},
						"United Kingdom": {},
					},
				},
			},
		},
		{
			name: "pushed unsorted",
			photos: []Photo{
				{
					Timestamp: time.Date(2022, 01, 03, 10, 11, 12, 0, time.UTC),
					Latitude:  -45,
					Longitude: 10,
					Addresses: map[string]struct{}{
						"London":         {},
						"United Kingdom": {},
					},
				},
				{
					Timestamp: time.Date(2022, 01, 02, 10, 11, 12, 0, time.UTC),
					Latitude:  -45,
					Longitude: 10,
					Addresses: map[string]struct{}{
						"London":         {},
						"United Kingdom": {},
					},
				},
			},
			expected: &PhotoHeap{
				{
					Timestamp: time.Date(2022, 01, 02, 10, 11, 12, 0, time.UTC),
					Latitude:  -45,
					Longitude: 10,
					Addresses: map[string]struct{}{
						"London":         {},
						"United Kingdom": {},
					},
				},
				{
					Timestamp: time.Date(2022, 01, 03, 10, 11, 12, 0, time.UTC),
					Latitude:  -45,
					Longitude: 10,
					Addresses: map[string]struct{}{
						"London":         {},
						"United Kingdom": {},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				h := &Heap{
					photoHeap: new(PhotoHeap),
					mu:        &sync.RWMutex{},
				}

				for _, photo := range tt.photos {
					h.Push(photo)
				}

				assert.Equal(t, tt.expected, h.photoHeap)
			},
		)
	}
}

func TestHeap_Pop(t *testing.T) {

	tests := []struct {
		name        string
		photoHeap   *PhotoHeap
		expected    Photo
		expectedErr error
	}{
		{
			name: "Pops item from heap",
			photoHeap: &PhotoHeap{
				{
					Timestamp: time.Date(2022, 01, 02, 10, 11, 12, 0, time.UTC),
					Latitude:  -45,
					Longitude: 10,
					Addresses: map[string]struct{}{
						"London":         {},
						"United Kingdom": {},
					},
				},
			},
			expected: Photo{
				Timestamp: time.Date(2022, 01, 02, 10, 11, 12, 0, time.UTC),
				Latitude:  -45,
				Longitude: 10,
				Addresses: map[string]struct{}{
					"London":         {},
					"United Kingdom": {},
				},
			},
			expectedErr: nil,
		},
		{
			name:        "Errors on empty heap",
			photoHeap:   &PhotoHeap{},
			expected:    Photo{},
			expectedErr: ErrEmptyHeap,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				h := &Heap{
					photoHeap: tt.photoHeap,
					mu:        &sync.RWMutex{},
				}
				got, err := h.Pop()
				assert.Equal(t, tt.expectedErr, err)
				assert.Equal(t, tt.expected, got)
			},
		)
	}
}
