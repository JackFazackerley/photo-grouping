package categoriser

import (
	"testing"
	"time"

	"github.com/JackFazackerley/photo-grouping/internal/heap"
	"github.com/stretchr/testify/assert"
)

func TestLocation_GenerateTitles(t *testing.T) {
	tests := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
		location  string
		expected  []string
	}{
		{
			name:      "generates weekend phrases",
			startTime: time.Date(2022, 04, 02, 10, 10, 10, 0, time.UTC),
			endTime:   time.Date(2022, 04, 04, 14, 10, 10, 0, time.UTC),
			location:  "London",
			expected: []string{
				"A weekend getaway to London",
				"A weekend in London",
				"London in April",
				"Visiting London in April",
			},
		},
		{
			name:      "generates week phrases",
			startTime: time.Date(2022, 03, 28, 10, 10, 10, 0, time.UTC),
			endTime:   time.Date(2022, 03, 30, 14, 10, 10, 0, time.UTC),
			location:  "London",
			expected: []string{
				"A trip away to London",
				"London in March",
				"Visiting London in March",
			},
		},
		{
			name:      "generates day phrases",
			startTime: time.Date(2022, 03, 28, 10, 10, 10, 0, time.UTC),
			endTime:   time.Date(2022, 03, 28, 14, 10, 10, 0, time.UTC),
			location:  "London",
			expected: []string{
				"A day out in London",
				"A trip to London",
				"London in March",
				"Visiting London in March",
			},
		},
		{
			name:      "generates holiday phrases",
			startTime: time.Date(2022, 03, 28, 10, 10, 10, 0, time.UTC),
			endTime:   time.Date(2022, 04, 05, 14, 10, 10, 0, time.UTC),
			location:  "London",
			expected: []string{
				"Holiday to London",
				"London in March",
				"Visiting London in March",
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				l := Location{
					startTime: tt.startTime,
					endTime:   tt.endTime,
					location:  tt.location,
				}
				got := l.GenerateTitles()

				assert.Equal(t, tt.expected, got)
			},
		)
	}
}

func TestGroup(t *testing.T) {
	tests := []struct {
		name     string
		photos   []heap.Photo
		expected map[string]*Location
	}{
		{
			name: "grouping by same location",
			photos: []heap.Photo{
				{
					Timestamp: time.Date(2022, 03, 28, 10, 10, 10, 0, time.UTC),
					Addresses: map[string]struct{}{
						"London": {},
					},
				},
				{
					Timestamp: time.Date(2022, 03, 28, 12, 10, 10, 0, time.UTC),
					Addresses: map[string]struct{}{
						"London": {},
					},
				},
			},
			expected: map[string]*Location{
				"London": {
					startTime: time.Date(2022, 03, 28, 10, 10, 10, 0, time.UTC),
					endTime:   time.Date(2022, 03, 28, 12, 10, 10, 0, time.UTC),
					location:  "London",
				},
			},
		},
		{
			name: "grouping by separate locations",
			photos: []heap.Photo{
				{
					Timestamp: time.Date(2022, 03, 28, 10, 10, 10, 0, time.UTC),
					Addresses: map[string]struct{}{
						"London":         {},
						"United Kingdom": {},
					},
				},
				{
					Timestamp: time.Date(2022, 03, 28, 12, 10, 10, 0, time.UTC),
					Addresses: map[string]struct{}{
						"London":         {},
						"United Kingdom": {},
					},
				},
			},
			expected: map[string]*Location{
				"London": {
					startTime: time.Date(2022, 03, 28, 10, 10, 10, 0, time.UTC),
					endTime:   time.Date(2022, 03, 28, 12, 10, 10, 0, time.UTC),
					location:  "London",
				},
				"United Kingdom": {
					startTime: time.Date(2022, 03, 28, 10, 10, 10, 0, time.UTC),
					endTime:   time.Date(2022, 03, 28, 12, 10, 10, 0, time.UTC),
					location:  "United Kingdom",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				photoHeap := heap.New()

				for _, photo := range tt.photos {
					photoHeap.Push(photo)
				}

				got := Group(photoHeap)

				assert.Equal(t, tt.expected, got)
			},
		)
	}
}
