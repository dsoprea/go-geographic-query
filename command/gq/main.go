package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dsoprea/go-geographic-index"
	"github.com/dsoprea/go-logging"
	"github.com/jessevdk/go-flags"
)

type parameters struct {
	Paths                       []string `short:"p" long:"path" description:"Path to recursively process GPS data from (can be provided zero or more times)"`
	Filepaths                   []string `short:"f" long:"filepath" description:"File-path to read GPS data from (can be provided zero or more times)"`
	TimestampQueries            []string `short:"t" long:"timestamp" description:"Timestamp query (formatted as RFC3339); can be provided zero or more times)"`
	LocationQueries             []string `short:"l" long:"location" description:"Location query (can be provided zero or more times)"`
	DoReverseGeocode            bool     `short:"g" long:"reverse-geocode" description:"Reverse-geocode coordinates to an address (set GOOGLE_API_KEY to your Google API-key)"`
	Verbose                     bool     `short:"v" long:"verbose" description:"Print logging"`
	ImagePath                   string   `short:"i" long:"image-path" description:"Path to walk for images and match with coordinates"`
	DoWalkImagePathsRecursively bool     `short:"r" long:"recursive-image-walk" description:"Recursively walk image path"`
	ShowImageSkips              bool     `short:"s" long:"show-image-skips" description:"Show images that are skipped"`
	ShowNearestImageTimes       bool     `short:"n" long:"show-nearest-image-times" description:"Show the datapoint times that were found when searching images in addition to the image times"`
}

var (
	arguments = new(parameters)
)

var (
	commandLogger = log.NewLogger("command/gq")
)

func main() {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintError(err)
			os.Exit(1)
		}
	}()

	p := flags.NewParser(arguments, flags.Default)

	_, err := p.Parse()
	if err != nil {
		os.Exit(1)
	}

	if arguments.Verbose == true {
		scp := log.NewStaticConfigurationProvider()
		scp.SetLevelName(log.LevelNameDebug)

		log.LoadConfiguration(scp)

		cla := log.NewConsoleLogAdapter()
		log.AddAdapter("console", cla)
	}

	if len(arguments.Paths) == 0 && len(arguments.Filepaths) == 0 {
		fmt.Printf("Please provide at least one file or one path.\n")
		os.Exit(2)
	}

	ti := geoindex.NewTimeIndex()
	gi := geoindex.NewGeographicIndex()

	gc := geoindex.NewGeographicCollector(ti, gi)

	err = geoindex.RegisterDataFileProcessors(gc)
	log.PanicIf(err)

	for _, path := range arguments.Paths {
		err := gc.ReadFromPath(path)
		log.PanicIf(err)
	}

	for _, filepath := range arguments.Filepaths {
		err := gc.ReadFromFilepath(filepath)
		log.PanicIf(err)
	}

	g := NewGeocoder()

	if arguments.ImagePath != "" {
		iw := NewImageWalk(arguments, g, ti, arguments.ImagePath)
		iw.Walk()

		os.Exit(0)
	}

	results := NewQueryResults()

	if len(arguments.TimestampQueries) > 0 {
		for _, timestampRaw := range arguments.TimestampQueries {
			timestamp, err := time.Parse(time.RFC3339, timestampRaw)
			log.PanicIf(err)

			tq := NewTimeQuery(timestamp)

			currentResults, err := runTimeQuery(ti, tq, arguments, g)
			if err == nil {
				for _, qr := range currentResults.Results() {
					results.Add(qr)
				}
			} else if err != ErrNotFound {
				log.Panic(err)
			}
		}
	}

	if len(arguments.LocationQueries) > 0 {
		for _, locationQuery := range arguments.LocationQueries {
			location, err := parseLocationQuery(locationQuery)
			log.PanicIf(err)

			latitude, longitude := location[0], location[1]
			lq := NewLocationQuery(latitude, longitude)

			currentResults, err := runLocationQuery(gi, lq, arguments, g)
			if err == nil {
				for _, qr := range currentResults.Results() {
					results.Add(qr)
				}
			} else if err != ErrNotFound {
				log.Panic(err)
			}
		}
	}

	if results.Len() > 0 {
		results.Print()
		os.Exit(0)
	}

	// If no command-line queries were given, fallthrough to the interactive
	// prompt.

	err = runConsole(ti, gi, arguments)
	log.PanicIf(err)
}
