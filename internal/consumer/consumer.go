package consumer

import (
	"context"
	"fmt"
	"sync"

	"github.com/JackFazackerley/photo-grouping/internal/heap"
	log "github.com/sirupsen/logrus"
	"googlemaps.github.io/maps"
)

var (
	// acceptedLocationTypes is used to only select specific maps.
	// AddressComponent types as the API doesn't support a filter for "locality"
	// see more here: https://developers.google.com/maps/documentation/places/web-service/supported_types#table2
	acceptedLocationTypes = []string{
		"country",
		"locality",
		"sublocality",
		"administrative_area_level_3",
		"administrative_area_level_2",
		"administrative_area_level_1",
	}
)

// geocodeClient is used so that unit tests can be easily written for the consumer
type geocodeClient interface {
	ReverseGeocode(ctx context.Context, r *maps.GeocodingRequest) ([]maps.GeocodingResult, error)
}

// Consumer is used to hold the maps.Client so that concurrently running Consumer.
// Run methods don't need the client passed in each time.
type Consumer struct {
	client geocodeClient
}

// NewConsumer is used to create the maps.Client and return an instance of Consumer.
// An API key is required in order to connect to the API.
// Options is a variadic argument allowing this package to be used with other maps.ClientOption(s).
func NewConsumer(APIKey string, options ...maps.ClientOption) (*Consumer, error) {
	options = append(options, maps.WithAPIKey(APIKey))

	client, err := maps.NewClient(options...)
	if err != nil {
		return nil, fmt.Errorf("creating new maps client: %w", err)
	}

	return &Consumer{
		client: client,
	}, nil
}

// Run is used to consume parsed heap.Photo(s) from the channel and add to the heap.Heap.
// If the context or channel is closed this method will return early.
func (c *Consumer) Run(ctx context.Context, photoHeap *heap.Heap, photos <-chan heap.Photo, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case photo, ok := <-photos:
			if !ok {
				return
			}

			if err := c.getGeocoding(ctx, photoHeap, photo); err != nil {
				log.WithError(err).Error("getting location")
			}
		}
	}
}

// getGeocoding is used to communicate with Google's ReverseGeocode API,
// the Latitude and Longitude are used to return the approximate location of the photo.
// The results from the API are then processed and duplicated results are guaranteed to not occur on the heap.
// Photo from the use of a hashmap.
//
// Once each heap.Photo's addresses have been stored it will then be pushed onto the heap.Heap and sorted.
//
// if the request to the API fails, an error is returned.
func (c *Consumer) getGeocoding(ctx context.Context, photoHeap *heap.Heap, photo heap.Photo) error {
	results, err := c.client.ReverseGeocode(
		ctx, &maps.GeocodingRequest{
			LatLng: &maps.LatLng{
				Lat: photo.Latitude,
				Lng: photo.Longitude,
			},
			ResultType: []string{
				"locality",
			},
		},
	)
	if err != nil {
		return fmt.Errorf("getting location: %w", err)
	}

	if len(results) > 0 {
		photo.Addresses = make(map[string]struct{})

		for _, result := range results {
			for _, address := range result.AddressComponents {
				if acceptedTypes(address.Types) {
					if _, ok := photo.Addresses[address.LongName]; !ok {
						photo.Addresses[address.LongName] = struct{}{}
					}
				}
			}
		}

		photoHeap.Push(photo)
	}

	return nil
}

// acceptedTypes is used to determine if any of the address types from locationTypes are present within
// acceptedLocationTypes. If there is a match we end early and return true for a match, otherwise we return false.
func acceptedTypes(locationTypes []string) bool {
	for _, acceptedLocation := range acceptedLocationTypes {
		for _, locationType := range locationTypes {
			if acceptedLocation == locationType {
				return true
			}
		}
	}
	return false
}
