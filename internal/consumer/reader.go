package consumer

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/JackFazackerley/photo-grouping/internal/heap"
	"github.com/araddon/dateparse"
	"github.com/sirupsen/logrus"
)

// Reader is used to read a file, parse and output rows to a channel.
type Reader struct {
	file io.ReadCloser
}

// NewReader is used to open a file from the given filePath and will return a new instance of Reader.
func NewReader(filePath string) (*Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	return &Reader{
		file: file,
	}, nil
}

// ReadCSV is used to read the contents of Reader.file and return a channel.
// Rows are parses one at a time and a heap.Photo is created, then pushed onto the channel.
//
// Before the channel is returned, a go routine is created which will read the contents of the file.
// This method owns the channel due to it knowing when the channel should be closed,
// which is once all lines have been read.
//
// ReadCSV honours context.Context, before attempting to read a line it will check to see if the context
// has been cancelled, if it hasn't the line will be read, otherwise the go routine exits
func (r *Reader) ReadCSV(ctx context.Context) <-chan heap.Photo {
	reader := csv.NewReader(r.file)
	photoChan := make(chan heap.Photo)

	go func(ctx context.Context) {
		defer func() {
			close(photoChan)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				row, err := reader.Read()
				if err != nil && errors.Is(err, io.EOF) {
					return
				}

				if len(row) == 3 {
					timestamp, err := dateparse.ParseAny(row[0])
					if err != nil {
						logrus.WithError(err).Error("parsing timestamp")
						continue
					}

					latitude, err := strconv.ParseFloat(row[1], 64)
					if err != nil {
						logrus.WithError(err).Error("parsing latitude")
						continue
					}

					longitude, err := strconv.ParseFloat(row[2], 64)
					if err != nil {
						logrus.WithError(err).Error("parsing longitude")
						continue
					}

					photoChan <- heap.Photo{
						Timestamp: timestamp,
						Longitude: longitude,
						Latitude:  latitude,
					}
				}
			}
		}
	}(ctx)

	return photoChan
}

// Close wraps the Reader.file Close, so that the file may be safely closed.
func (r *Reader) Close() error {
	return r.file.Close()
}
