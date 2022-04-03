package consumer

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/JackFazackerley/photo-grouping/internal/heap"
	"github.com/stretchr/testify/assert"
	"googlemaps.github.io/maps"
)

const (
	apiKey = "some_api_key"
)

type mockGeocodingClient struct {
	result []maps.GeocodingResult
	err    error
}

func (m mockGeocodingClient) ReverseGeocode(ctx context.Context, r *maps.GeocodingRequest) ([]maps.GeocodingResult, error) {
	return m.result, m.err
}

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		apiKey      string
		expected    *Consumer
		expectedErr string
	}{
		{
			name:        "creates new consumer",
			apiKey:      apiKey,
			expectedErr: "",
		},
		{
			name:        "errors creating new maps client",
			apiKey:      "",
			expectedErr: "creating new maps client",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				_, err := NewConsumer(tt.apiKey)

				if tt.expectedErr == "" {
					assert.NoError(t, err)
				} else {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
			},
		)
	}
}

func TestConsumer_Run(t *testing.T) {
	tests := []struct {
		name              string
		photo             heap.Photo
		client            geocodeClient
		earlyChannelClose bool
		expected          heap.Photo
	}{
		{
			name: "adds photo to the heap",
			photo: heap.Photo{
				Timestamp: time.Date(2022, 01, 03, 10, 11, 12, 0, time.UTC),
				Latitude:  51.5072,
				Longitude: 0.1276,
			},
			client: mockGeocodingClient{
				result: []maps.GeocodingResult{
					{
						AddressComponents: []maps.AddressComponent{
							{
								LongName:  "London",
								ShortName: "London",
								Types: []string{
									"locality",
								},
							},
							{
								LongName:  "Greater London",
								ShortName: "Greater London",
								Types: []string{
									"administrative_area_level_2",
								},
							},
							{
								LongName:  "England",
								ShortName: "England",
								Types: []string{
									"administrative_area_level_1",
								},
							},
							{
								LongName:  "United Kingdom",
								ShortName: "GB",
								Types: []string{
									"country",
								},
							},
						},
					},
				},
				err: nil,
			},
			expected: heap.Photo{
				Timestamp: time.Date(2022, 01, 03, 10, 11, 12, 0, time.UTC),
				Latitude:  51.5072,
				Longitude: 0.1276,
				Addresses: map[string]struct{}{
					"London":         {},
					"Greater London": {},
					"England":        {},
					"United Kingdom": {},
				},
			},
		},
		{
			name: "nothing added to the heap on error",
			photo: heap.Photo{
				Timestamp: time.Date(2022, 01, 03, 10, 11, 12, 0, time.UTC),
				Latitude:  51.5072,
				Longitude: 0.1276,
			},
			client: mockGeocodingClient{
				err: errors.New("client error"),
			},
			expected: heap.Photo{},
		},
		{
			name: "exits on channel close",
			photo: heap.Photo{
				Timestamp: time.Date(2022, 01, 03, 10, 11, 12, 0, time.UTC),
				Latitude:  51.5072,
				Longitude: 0.1276,
			},
			client:            mockGeocodingClient{},
			earlyChannelClose: true,
			expected:          heap.Photo{},
		},
		{
			name: "deduplicates on known addresses",
			photo: heap.Photo{
				Timestamp: time.Date(2022, 01, 03, 10, 11, 12, 0, time.UTC),
				Latitude:  51.5072,
				Longitude: 0.1276,
			},
			client: mockGeocodingClient{
				result: []maps.GeocodingResult{
					{
						AddressComponents: []maps.AddressComponent{
							{
								LongName:  "London",
								ShortName: "London",
								Types: []string{
									"locality",
								},
							},
							{
								LongName:  "London",
								ShortName: "London",
								Types: []string{
									"locality",
								},
							},
							{
								LongName:  "Greater London",
								ShortName: "Greater London",
								Types: []string{
									"administrative_area_level_2",
								},
							},
							{
								LongName:  "England",
								ShortName: "England",
								Types: []string{
									"administrative_area_level_1",
								},
							},
							{
								LongName:  "United Kingdom",
								ShortName: "GB",
								Types: []string{
									"country",
								},
							},
						},
					},
				},
				err: nil,
			},
			expected: heap.Photo{
				Timestamp: time.Date(2022, 01, 03, 10, 11, 12, 0, time.UTC),
				Latitude:  51.5072,
				Longitude: 0.1276,
				Addresses: map[string]struct{}{
					"London":         {},
					"Greater London": {},
					"England":        {},
					"United Kingdom": {},
				},
			},
		},
		{
			name: "ignores unaccepted address types",
			photo: heap.Photo{
				Timestamp: time.Date(2022, 01, 03, 10, 11, 12, 0, time.UTC),
				Latitude:  51.5072,
				Longitude: 0.1276,
			},
			client: mockGeocodingClient{
				result: []maps.GeocodingResult{
					{
						AddressComponents: []maps.AddressComponent{
							{
								LongName:  "London",
								ShortName: "London",
								Types: []string{
									"locality",
								},
							},
							{
								LongName:  "E1 6AN",
								ShortName: "",
								Types: []string{
									"postal_code",
								},
							},
						},
					},
				},
				err: nil,
			},
			expected: heap.Photo{
				Timestamp: time.Date(2022, 01, 03, 10, 11, 12, 0, time.UTC),
				Latitude:  51.5072,
				Longitude: 0.1276,
				Addresses: map[string]struct{}{
					"London": {},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*250)
				defer cancel()

				c := &Consumer{
					client: tt.client,
				}

				wg := &sync.WaitGroup{}
				photosChan := make(chan heap.Photo, 1)
				photoHeap := heap.New()

				photosChan <- tt.photo

				if tt.earlyChannelClose {
					close(photosChan)
				}

				wg.Add(1)
				go c.Run(ctx, photoHeap, photosChan, wg)

				wg.Wait()

				got, _ := photoHeap.Pop()
				assert.Equal(t, tt.expected, got)
			},
		)
	}
}
