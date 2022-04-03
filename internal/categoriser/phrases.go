package categoriser

import (
	"fmt"
)

var (
	anyPhrases = []phrase{
		delimiterPhrase("in"),
		combinationPhrase{
			phrase:    "Visiting",
			delimiter: "in",
		},
	}

	weekendPhrases = []phrase{
		locationPhrase("A weekend getaway to"),
		locationPhrase("A weekend in"),
	}

	weekPhrases = []phrase{
		locationPhrase("A trip away to"),
	}

	dayPhrases = []phrase{
		locationPhrase("A day out in"),
		locationPhrase("A trip to"),
	}

	holidayPhrases = []phrase{
		locationPhrase("Holiday to"),
	}
)

// phrase is used as a generic interface that can be used for multiple types of phrases.
// Allowing only one entry point (phrase.generate).
type phrase interface {
	generate(location Location) string
}

// locationPhrase is used to create a title based on the provided string and the location of the Location
type locationPhrase string

func (l locationPhrase) generate(location Location) string {
	return fmt.Sprintf("%s %s", l, location.location)
}

// delimiterPhrase is used to create a title based on the Location.location,
// the provided string and the month of Location.startTime
type delimiterPhrase string

func (d delimiterPhrase) generate(location Location) string {
	return fmt.Sprintf("%s %s %s", location.location, d, location.startTime.Month())
}

// combinationPhrase is used to create a title from a starting combinationPhrase.phrase, the Location.location,
// combinationPhrase.delimiter, and the month of Location.startTime.
type combinationPhrase struct {
	phrase    string
	delimiter string
}

func (c combinationPhrase) generate(location Location) string {
	return fmt.Sprintf("%s %s %s %s", c.phrase, location.location, c.delimiter, location.startTime.Month())
}
