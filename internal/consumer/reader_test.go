package consumer

import (
	"context"
	"io"
	"io/fs"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/JackFazackerley/photo-grouping/internal/heap"
	"github.com/stretchr/testify/assert"
)

func TestNewReader(t *testing.T) {
	tests := []struct {
		name        string
		createFile  bool
		expected    *Reader
		expectedErr error
	}{
		{
			name:       "opens file",
			createFile: true,
			expected: &Reader{
				file: &os.File{},
			},
			expectedErr: nil,
		},
		{
			name:        "errors opening file",
			createFile:  false,
			expected:    nil,
			expectedErr: fs.ErrNotExist,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				filePath := "temp"

				if tt.createFile {
					f, err := os.CreateTemp(".", "*.csv")
					if err != nil {
						t.Fatalf("creating temp file: %s", err)
					}
					_ = f.Close()

					defer func() {
						_ = os.Remove(f.Name())
					}()

					filePath = f.Name()
				}

				got, err := NewReader(filePath)

				assert.ErrorIs(t, err, tt.expectedErr)

				assert.IsType(t, tt.expected, got)

				if got != nil {
					assert.IsType(t, tt.expected.file, got.file)
					_ = got.Close()
				}
			},
		)
	}
}

func TestReader_ReadCSV(t *testing.T) {
	tests := []struct {
		name         string
		fileContents string
		earlyCancel  bool
		expected     []heap.Photo
	}{
		{
			name:         "parses rows without errors",
			fileContents: "2020-03-30 14:12:19,40.728808,-73.996106\n2020-03-30 14:20:10,40.728656,-73.998790\n2020-03-30 14:32:02,40.727160,-73.996044",
			expected: []heap.Photo{
				{
					Timestamp: time.Date(2020, 03, 30, 14, 12, 19, 0, time.UTC),
					Latitude:  40.728808,
					Longitude: -73.996106,
				},
				{
					Timestamp: time.Date(2020, 03, 30, 14, 20, 10, 0, time.UTC),
					Latitude:  40.728656,
					Longitude: -73.998790,
				},
				{
					Timestamp: time.Date(2020, 03, 30, 14, 32, 02, 0, time.UTC),
					Latitude:  40.727160,
					Longitude: -73.996044,
				},
			},
		},
		{
			name:         "error parsing time",
			fileContents: "not_a_time,40.728808,-73.996106",
			expected:     []heap.Photo{},
		},
		{
			name:         "error parsing latitude",
			fileContents: "2020-03-30 14:12:19,not_a_float,-73.996106",
			expected:     []heap.Photo{},
		},
		{
			name:         "error parsing longitude",
			fileContents: "2020-03-30 14:12:19,40.728808,not_a_float",
			expected:     []heap.Photo{},
		},
		{
			name:         "context cancelled",
			earlyCancel:  true,
			fileContents: "2020-03-30 14:12:19,40.728808,-73.996106",
			expected:     []heap.Photo{},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*250)
				defer cancel()

				if tt.earlyCancel {
					cancel()
				}

				r := &Reader{
					file: io.NopCloser(strings.NewReader(tt.fileContents)),
				}

				photoChan := r.ReadCSV(ctx)

				got := make([]heap.Photo, 0)

				for photo := range photoChan {
					got = append(got, photo)
				}

				assert.Equal(t, tt.expected, got)
			},
		)
	}
}
