package main

import (
	"fmt"

	"github.com/codingsince1985/geo-golang/openstreetmap"

	"github.com/dsoprea/go-logging"
)

type Place struct {
	description string
}

func (p Place) String() string {
	return p.description
}

type Geocoder struct {
	cachedPlaces map[string]Place
}

func NewGeocoder() *Geocoder {
	return &Geocoder{
		cachedPlaces: make(map[string]Place),
	}
}

func (g *Geocoder) getAddressForCoordinates(latitude, longitude float64) (p Place, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	// Avoids epsilon issues.
	key := fmt.Sprintf("%.6s,%.6s", latitude, longitude)

	if p, found := g.cachedPlaces[key]; found == false {
		osg := openstreetmap.Geocoder()

		info, err := osg.ReverseGeocode(latitude, longitude)
		log.PanicIf(err)

		p = Place{
			description: info.FormattedAddress,
		}

		g.cachedPlaces[key] = p
		return p, nil
	} else {
		return p, nil
	}
}

func isGeocodeSupported() bool {
	return true
}
