package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dsoprea/go-logging"
)

var (
	ErrQueryNotValid = errors.New("query not valid")
)

type Query interface {
	QueryType() string
}

type query struct {
	queryType string
}

func (q query) QueryType() string {
	return q.queryType
}

type LocationQuery struct {
	Latitude, Longitude float64

	query
}

type TimeQuery struct {
	Timestamp time.Time

	query
}

func parseQuery(rawLine string) (q Query, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	rawLine = strings.TrimSpace(rawLine)

	var timestamp time.Time
	var latitude float64
	var longitude float64
	found := false
	if strings.HasPrefix(rawLine, "t ") == true {
		parsed, err := time.Parse(time.RFC3339, rawLine[2:])
		if err != nil {
			fmt.Printf("Could not parse annotated time query.\n")
			return nil, ErrQueryNotValid
		}

		timestamp = parsed
		found = true
	} else if strings.HasPrefix(rawLine, "l ") == true {
		locationQuery := rawLine[2:]

		location, err := parseLocationQuery(locationQuery)
		if err != nil {
			fmt.Printf("Could not parse annotated location query.\n")
			return nil, ErrQueryNotValid
		}

		latitude, longitude = location[0], location[1]
		found = true
	} else {
		parsed, err := time.Parse(time.RFC3339, rawLine)
		if err == nil {
			timestamp = parsed
			found = true
		}

		if found == false {
			location, err := parseLocationQuery(rawLine)
			if err != nil {
				return nil, ErrQueryNotValid
			}

			latitude, longitude = location[0], location[1]
			found = true
		}

		if found == false {
			fmt.Printf("Query layout not recognized.\n")
			return nil, ErrQueryNotValid
		}
	}

	// If we got to here, we successfully parse the query.

	if timestamp.IsZero() == false {
		q = NewTimeQuery(timestamp)
		return q, nil
	} else {
		q = NewLocationQuery(latitude, longitude)
		return q, nil
	}
}

func NewTimeQuery(timestamp time.Time) TimeQuery {
	return TimeQuery{
		Timestamp: timestamp,
		query: query{
			queryType: "timestamp",
		},
	}
}

func NewLocationQuery(latitude, longitude float64) LocationQuery {
	return LocationQuery{
		Latitude:  latitude,
		Longitude: longitude,
		query: query{
			queryType: "location",
		},
	}
}

func parseLocationQuery(locationQuery string) (location []float64, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	parts := strings.Split(locationQuery, ",")
	if len(parts) != 2 {
		log.Panic("could not parse location")
	}

	latitudeRaw := strings.TrimSpace(parts[0])
	longitudeRaw := strings.TrimSpace(parts[1])

	latitude, err := strconv.ParseFloat(latitudeRaw, 64)
	log.PanicIf(err)

	longitude, err := strconv.ParseFloat(longitudeRaw, 64)
	log.PanicIf(err)

	return []float64{latitude, longitude}, nil
}
