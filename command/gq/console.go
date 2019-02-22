package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dsoprea/go-geographic-index"
	"github.com/dsoprea/go-logging"
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
