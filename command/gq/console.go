package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/dsoprea/go-geographic-index"
	"github.com/dsoprea/go-logging"
)

var (
	ErrNotFound = errors.New("not found")
)

func runConsole(ti *geoindex.TimeIndex, gi *geoindex.GeographicIndex, arguments *parameters) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	r := bufio.NewReader(os.Stdin)
	g := NewGeocoder()

	for {
		fmt.Printf("> ")

		query, err := r.ReadString('\n')
		log.PanicIf(err)

		query = strings.TrimSpace(query)
		if len(query) == 0 {
			continue
		}

		q, err := parseQuery(query)
		if err != nil {
			if err == ErrQueryNotValid {
				fmt.Printf("Parse error.\n")
				continue
			}

			log.PanicIf(err)
		}

		switch v := q.(type) {
		case LocationQuery:
			results, err := runLocationQuery(gi, v, arguments, g)
			if err == nil {
				results.Print()
			} else if err != ErrNotFound {
				log.Panic(err)
			}
		case TimeQuery:
			results, err := runTimeQuery(ti, v, arguments, g)
			if err == nil {
				results.Print()
			} else if err != ErrNotFound {
				log.Panic(err)
			}
		default:
			log.Panicf("query-type error")
		}
	}

	return nil
}

func runLocationQuery(gi *geoindex.GeographicIndex, q LocationQuery, arguments *parameters, g *Geocoder) (results *QueryResults, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	geographicResults, err := gi.GetWithCoordinatesMetroLimited(q.Latitude, q.Longitude)
	if err != nil {
		if err == geoindex.ErrNoNearMatch {
			return nil, ErrNotFound
		}

		log.Panic(err)
	}

	results = NewQueryResults()
	for _, gr := range geographicResults {
		qr := &QueryResult{
			gr: gr,
		}

		if arguments.DoReverseGeocode == true && isGeocodeSupported() == true {
			place, err := g.getAddressForCoordinates(gr.Latitude, gr.Longitude)
			log.PanicIf(err)

			qr.p = place
		}

		results.Add(qr)
	}

	return results, nil
}

func runTimeQuery(ti *geoindex.TimeIndex, q TimeQuery, arguments *parameters, g *Geocoder) (results *QueryResults, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	ts := ti.Series()
	i := ts.Search(q.Timestamp)

	if i >= len(ts) {
		fmt.Printf("No record found for: [%s] (1)\n", q.Timestamp)
		return nil, ErrNotFound
	}

	result := ts[i]
	if result.Time != q.Timestamp {
		fmt.Printf("No record found for: [%s] (2)\n", q.Timestamp)
		return nil, ErrNotFound
	}

	results = NewQueryResults()
	for _, o := range result.Items {
		gr := o.(*geoindex.GeographicRecord)

		qr := &QueryResult{
			gr: gr,
		}

		if arguments.DoReverseGeocode == true && isGeocodeSupported() == true {
			place, err := g.getAddressForCoordinates(gr.Latitude, gr.Longitude)
			log.PanicIf(err)

			qr.p = place
		}

		results.Add(qr)
	}

	return results, nil
}
