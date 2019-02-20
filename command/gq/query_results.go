package main

import (
	"fmt"
	"sort"

	"github.com/dsoprea/go-geographic-index"
)

var (
	NoPlace = Place{}
)

type QueryResult struct {
	gr *geoindex.GeographicRecord
	p  Place
}

func (qr *QueryResult) String() string {
	description := fmt.Sprintf("%s (%.6f, %.6f)", qr.gr.Timestamp, qr.gr.Latitude, qr.gr.Longitude)

	if qr.p != NoPlace {
		description = fmt.Sprintf("%s: %s", description, qr.p)
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
