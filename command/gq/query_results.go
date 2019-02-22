package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/dsoprea/go-geographic-index"
)

var (
	NoPlace = Place{}
)

type QueryResult struct {
	gr *geoindex.GeographicRecord
	p  Place

	// originalTimestamp is an optional, original timestamp before any rounding
	// was employed in the search.
	originalTimestamp time.Time
}

func (qr *QueryResult) String() string {
	timestampPhrase := ""
	if qr.originalTimestamp.IsZero() == true {
		timestampPhrase = fmt.Sprintf("%s", qr.gr.Timestamp)
	} else {
		timestampPhrase = fmt.Sprintf("%s\t[<-%s]", qr.gr.Timestamp, qr.originalTimestamp)
	}

	description := fmt.Sprintf("%s\t(%.6f,%.6f)", timestampPhrase, qr.gr.Latitude, qr.gr.Longitude)

	if qr.p != NoPlace {
		description = fmt.Sprintf("%s\t%s", description, qr.p)
	}

	return description
}

type QueryResults struct {
	results []*QueryResult
}

func NewQueryResults() *QueryResults {
	return &QueryResults{
		results: make([]*QueryResult, 0),
	}
}

func (results *QueryResults) Results() []*QueryResult {
	return results.results
}

func (results *QueryResults) Add(qr *QueryResult) {
	results.results = append(results.results, qr)
}

func (results *QueryResults) Len() int {
	return len(results.results)
}

func (results *QueryResults) Less(i, j int) bool {
	return results.results[i].gr.Timestamp.Before(results.results[j].gr.Timestamp)
}

func (results *QueryResults) Swap(i, j int) {
	results.results[j], results.results[i] = results.results[i], results.results[j]
}

func (results *QueryResults) Print() {
	sort.Sort(results)

	for _, qr := range results.results {
		fmt.Printf("%s\n", qr)
	}
}
