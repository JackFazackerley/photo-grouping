package categoriser

import (
	"errors"
	"time"

	"github.com/JackFazackerley/photo-grouping/internal/heap"
)

var (
	oneDay   = time.Hour * 24
	fourDays = time.Hour * 96
)

// Location defines a location. startTime used to told when the first photo taken at the location and endTime
// being the last photo taken at the location
type Location struct {
	startTime time.Time
	endTime   time.Time
	location  string
}

// GenerateTitles is used to determine what type of trip (day,weekend,week,holiday)
// the Location was. Once determined it will call suggestions to generate a list of story titles and return the result.
func (l Location) GenerateTitles() []string {
	var phrases []phrase

	startTimeWeekday := l.startTime.Weekday()
	endTimeWeekday := l.endTime.Weekday()

	tripDiff := l.endTime.Sub(l.startTime)
	if tripDiff < oneDay {
		phrases = dayPhrases
	} else if tripDiff > oneDay && tripDiff <= fourDays {
		if startTimeWeekday >= time.Friday && endTimeWeekday <= time.Monday {
			phrases = weekendPhrases
		} else {
			phrases = weekPhrases
		}
	} else {
		phrases = holidayPhrases
	}

	phrases = append(phrases, anyPhrases...)

	return l.suggestions(phrases)
}

// suggestions generates titles based on the phrase(s) provided.
func (l Location) suggestions(phrases []phrase) []string {
	tripNames := make([]string, 0)

	for _, phrase := range phrases {
		tripNames = append(tripNames, phrase.generate(l))
	}

	return tripNames
}

// Group is used to group photos together based on the location of the photos.
// In order to group photos together with a degree of confidence there is an assumption that each photo is taken
// within 24 hours of the previous photo. If the photo is within 24 hours of the previous the Location endTime is set
// to the current photo.
//
// As this function uses heap.Heap we can be certain that all photos are in order based on the timestamp when popping
// from the heap.
func Group(photoHeap *heap.Heap) map[string]*Location {
	locations := make(map[string]*Location)

	for {
		photo, err := photoHeap.Pop()
		if err != nil {
			if errors.Is(err, heap.ErrEmptyHeap) {
				break
			}
		}

		for k := range photo.Addresses {
			location := &Location{
				startTime: photo.Timestamp,
				endTime:   photo.Timestamp,
				location:  k,
			}

			if lastLocation, ok := locations[k]; !ok {
				locations[k] = location
			} else {
				if photo.Timestamp.Sub(lastLocation.endTime) <= oneDay {
					lastLocation.endTime = photo.Timestamp
				}
			}
		}
	}

	return locations
}
