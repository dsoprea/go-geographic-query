package main

import (
	"errors"
	"time"

	"github.com/dsoprea/go-geographic-index"
	"github.com/dsoprea/go-logging"
)

const (
	searchInterval = time.Minute * 5
)

var (
	qLogger = log.NewLogger("gq.query")
)

var (
	ErrNotFound = errors.New("not found")
)

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

var (
	timeQueryCache = make(map[time.Time]*QueryResults)
)

func runTimeQuery(ti *geoindex.TimeIndex, q TimeQuery, arguments *parameters, g *Geocoder) (results *QueryResults, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	if results, found := timeQueryCache[q.Timestamp]; found == true {
		return results, nil
	}

	ts := ti.Series()

	var matchedTimestamp time.Time
	cb := func(t time.Time) error {
		defer func() {
			if state := recover(); state != nil {
				err = log.Wrap(state.(error))
				log.Panic(err)
			}
		}()

		matchedTimestamp = t

		return nil
	}

	// TODO(dustin): !! We need to accont for sparse datasets.

	err = ts.SearchNearest(q.Timestamp, searchInterval, cb)
	log.PanicIf(err)

	if matchedTimestamp.IsZero() == true {
		qLogger.Debugf(nil, "No nearest time for: [%s] (1)", q.Timestamp)
		return nil, ErrNotFound
	}

	i := ts.Search(matchedTimestamp)

	if i >= len(ts) {
		qLogger.Debugf(nil, "No nearest time for: [%s]->[%s] (2)", q.Timestamp, matchedTimestamp)
		return nil, ErrNotFound
	}

	result := ts[i]
	if result.Time != matchedTimestamp {
		qLogger.Debugf(nil, "No nearest time for: [%s]->[%s] (3)", q.Timestamp, matchedTimestamp)
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

			// Annotate the nearest times that were found if we want to show
			// them. If the annotation is there, it'll be included in the
			// output.
			if arguments.ShowNearestImageTimes == true {
				if qr.originalTimestamp.Year() != q.Timestamp.Year() ||
					qr.originalTimestamp.Month() != q.Timestamp.Month() ||
					qr.originalTimestamp.Day() != q.Timestamp.Day() ||
					qr.originalTimestamp.Hour() != q.Timestamp.Hour() ||
					qr.originalTimestamp.Minute() != q.Timestamp.Minute() ||
					qr.originalTimestamp.Second() != q.Timestamp.Second() {
					qr.originalTimestamp = q.Timestamp
				}
			}

			qr.p = place
		}

		results.Add(qr)
	}

	timeQueryCache[q.Timestamp] = results

	return results, nil
}
