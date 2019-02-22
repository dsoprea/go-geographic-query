package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"io/ioutil"

	"github.com/dsoprea/go-exif"
	"github.com/dsoprea/go-geographic-index"
	"github.com/dsoprea/go-jpeg-image-structure"
	"github.com/dsoprea/go-logging"
	"github.com/randomingenuity/go-utility/filesystem"
)

const (
	JpegImageExtension = ".jpg"
)

var (
	iwLogger = log.NewLogger("gq.image_walk")
)

type ImageWalk struct {
	arguments *parameters
	g         *Geocoder
	ti        *geoindex.TimeIndex
	imagePath string
}

func NewImageWalk(arguments *parameters, g *Geocoder, ti *geoindex.TimeIndex, imagePath string) *ImageWalk {
	return &ImageWalk{
		arguments: arguments,
		g:         g,
		ti:        ti,
		imagePath: imagePath,
	}
}

func (iw *ImageWalk) Walk() (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	if iw.arguments.DoWalkImagePathsRecursively == true {
		err := iw.recursiveReadFromPath(iw.imagePath)
		log.PanicIf(err)
	} else {
		entries, err := ioutil.ReadDir(iw.imagePath)
		log.PanicIf(err)

		for _, fi := range entries {
			filename := fi.Name()
			filepath := path.Join(iw.imagePath, filename)

			extension := path.Ext(fi.Name())
			extension = strings.ToLower(extension)

			if extension != JpegImageExtension {
				continue
			}

			err := iw.processImage(filepath)
			log.PanicIf(err)
		}
	}

	return nil
}

func (iw *ImageWalk) recursiveReadFromPath(rootPath string) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	// Allow all directories and any file whose extension is associated with a
	// processor.
	filter := func(parent string, child os.FileInfo) (bool, error) {
		if child.IsDir() == true {
			return true, nil
		}

		extension := path.Ext(child.Name())
		extension = strings.ToLower(extension)

		if extension == JpegImageExtension {
			return true, nil
		}

		return false, nil
	}

	filesC, errC := rifs.ListFiles(rootPath, filter)

FilesRead:

	for {
		select {
		case err, ok := <-errC:
			if ok == true {
				// TODO(dustin): Can we close these on the other side after sending and still get our data?
				close(filesC)
				close(errC)
			}

			log.PanicIf(err)

		case vf, ok := <-filesC:
			// We have finished reading. `vf` has an empty value.
			if ok == false {
				// The goroutine finished.
				break FilesRead
			}

			if vf.Info.IsDir() == false {
				err := iw.processImage(vf.Filepath)
				log.PanicIf(err)
			}
		}
	}

	return nil
}

func (iw *ImageWalk) processImage(filepath string) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	relFilepath := filepath[len(iw.imagePath)+1:]

	timestamp, err := getTimestampFromJpeg(filepath)
	if err != nil {
		iwLogger.Warningf(nil, "Skipping unreadable/unparseable image: [%s] [%s]", relFilepath, err)

		if iw.arguments.ShowImageSkips == true {
			fmt.Printf("! %s: Unreadable/unparseable\n", relFilepath)
		}

		return nil
	}

	tq := NewTimeQuery(timestamp)
	results, err := runTimeQuery(iw.ti, tq, iw.arguments, iw.g)
	if err != nil {
		iwLogger.Warningf(nil, "Could not match image to time: [%s] [%s]", relFilepath, err)

		if iw.arguments.ShowImageSkips == true {
			fmt.Printf("! %s: No match\n", relFilepath)
		}

		return nil
	}

	innerResults := results.Results()
	qr := innerResults[0]

	fmt.Printf("%s\t%s\n", relFilepath, qr)

	return nil
}

func getFirstExifTagStringValue(rootIfd *exif.Ifd, tagName string) (value string, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	results, err := rootIfd.FindTagWithName(tagName)
	if err != nil {
		if log.Is(err, exif.ErrTagNotFound) == true {
			results = nil
		} else {
			log.Panic(err)
		}
	} else {
		if len(results) == 0 {
			results = nil
		}
	}

	if results != nil {
		ite := results[0]

		valueRaw, err := rootIfd.TagValue(ite)
		log.PanicIf(err)

		value = valueRaw.(string)
	}

	return value, nil
}

func getTimestampFromJpeg(filepath string) (timestamp time.Time, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	jmp := jpegstructure.NewJpegMediaParser()

	data, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	sl, err := jmp.ParseBytes(data)
	log.PanicIf(err)

	rootIfd, _, err := sl.Exif()
	if err != nil {
		// Skip if it doesn't have EXIF data.
		if log.Is(err, jpegstructure.ErrNoExif) == true {
			return time.Time{}, nil
		}

		log.Panic(err)
	}

	// Get the picture timestamp as stored in the EXIF.

	tagName := "DateTime"

	timestampPhrase, err := getFirstExifTagStringValue(rootIfd, tagName)
	log.PanicIf(err)

	if timestampPhrase == "" {
		iwLogger.Warningf(nil, "Image has an empty timestamp: [%s]", filepath)
		return time.Time{}, nil
	}

	timestamp, err = exif.ParseExifFullTimestamp(timestampPhrase)
	if err != nil {
		iwLogger.Warningf(nil, "Image's timestamp is unparseable: [%s] [%s]", filepath, timestampPhrase)
		return time.Time{}, nil
	}

	return timestamp, nil
}
