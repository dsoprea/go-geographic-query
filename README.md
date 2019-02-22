# Overview

This tool provides an array of functionalities for efficiently querying raw GPS data from the command-line. And it's awesome.


# Features

- Provide multiple GPX files or paths and they'll be loaded and indexed.
- Provide multiple timestamp and/or coordinate queries via the command-line, and they'll be resolved against the data.
- If no queries are provided via the command-line, a console will be opened.
- The tool can also walk a given path for all JPEGs, parse EXIF for a timestamp, and lookup and print locations.
- Coordinates can be looked-up against OpenStreetMap in order to show human-readable locations.
- Most lookups are cached for the current session.


# Getting

```
$ go get -t github.com/dsoprea/go-geographic-query
```


# Examples

These examples were [obviously] run directly from "$GOPATH/github.com/dsoprea/go-geographic-query/command/gq".

Simple query with a single timestamp:

```
$ go run . --path=$HOME/gpxdata --timestamp=2018-02-18T17:02:31Z
2018-02-18 17:02:31.155 +0000 UTC   (27.048786,-82.178683)
```

The same query, with location lookup:

```
$ go run . --path=$HOME/gpxdata --timestamp=2018-02-18T17:02:31Z --reverse-geocode
2018-02-18 17:02:31.155 +0000 UTC   (27.048786,-82.178683)  3705, Parkins Terrace, North Port, Sarasota County, Florida, 34286, USA
```


The query console (no queries are passed via command-line):

```
$ go run . --path=$HOME/gpxdata
> 2018-02-18T17:02:31Z
2018-02-18 17:02:31.155 +0000 UTC   (27.048786,-82.178683)
> 27.048786,-82.178683
2018-02-18 17:02:31.155 +0000 UTC   (27.048786,-82.178683)
```


Walking and matching images:

```
$ go run . --path=$HOME/gpxdata --image-path=$HOME/images/DCIM --reverse-geocode --recursive-image-walk
10880218/DSC00579.JPG   2018-02-18 16:52:46.743 +0000 UTC   (27.048828,-82.178365)  3650, Parkins Terrace, North Port, Sarasota County, Florida, 34286, USA
10880218/DSC00636.JPG   2018-02-18 17:02:31.155 +0000 UTC   (27.048786,-82.178683)  3705, Parkins Terrace, North Port, Sarasota County, Florida, 34286, USA
10880218/DSC00625.JPG   2018-02-18 17:02:31.155 +0000 UTC   (27.048786,-82.178683)  3705, Parkins Terrace, North Port, Sarasota County, Florida, 34286, USA
10880218/DSC00578.JPG   2018-02-18 16:52:46.743 +0000 UTC   (27.048828,-82.178365)  3650, Parkins Terrace, North Port, Sarasota County, Florida, 34286, USA
10880218/DSC00647.JPG   2018-02-18 17:02:31.155 +0000 UTC   (27.048786,-82.178683)  3705, Parkins Terrace, North Port, Sarasota County, Florida, 34286, USA
10880218/DSC00461.JPG   2018-02-18 14:37:47.171 +0000 UTC   (27.048798,-82.178924)  3691, Parkins Terrace, North Port, Sarasota County, Florida, 34286, USA
10880218/DSC00452.JPG   2018-02-18 12:42:19.743 +0000 UTC   (27.048798,-82.178962)  3689, Parkins Terrace, North Port, Sarasota County, Florida, 34286, USA
10880218/DSC00471.JPG   2018-02-18 14:49:10.081 +0000 UTC   (27.048791,-82.178872)  3691, Parkins Terrace, North Port, Sarasota County, Florida, 34286, USA
10880218/DSC00649.JPG   2018-02-18 17:19:50.761 +0000 UTC   (27.103028,-82.289727)  I 75, Sarasota County, Florida, USA
10880218/DSC00659.JPG   2018-02-18 22:51:10.944 +0000 UTC   (27.140116,-82.443547)  914, Dartmoor Circle, Sarasota County, Florida, 34275, USA
10880218/DSC00641.JPG   2018-02-18 17:02:31.155 +0000 UTC   (27.048786,-82.178683)  3705, Parkins Terrace, North Port, Sarasota County, Florida, 34286, USA
10880218/DSC00634.JPG   2018-02-18 17:02:31.155 +0000 UTC   (27.048786,-82.178683)  3705, Parkins Terrace, North Port, Sarasota County, Florida, 34286, USA
10880218/DSC00627.JPG   2018-02-18 17:02:31.155 +0000 UTC   (27.048786,-82.178683)  3705, Parkins Terrace, North Port, Sarasota County, Florida, 34286, USA
10880218/DSC00449.JPG   2018-02-18 12:36:17.673 +0000 UTC   (27.048928,-82.179036)  3790, Waffle Terrace, North Port, Sarasota County, Florida, 34286, USA
10880218/DSC00621.JPG   2018-02-18 17:02:31.155 +0000 UTC   (27.048786,-82.178683)  3705, Parkins Terrace, North Port, Sarasota County, Florida, 34286, USA
10880218/DSC00639.JPG   2018-02-18 17:02:31.155 +0000 UTC   (27.048786,-82.178683)  3705, Parkins Terrace, North Port, Sarasota County, Florida, 34286, USA
10880218/DSC00486.JPG   2018-02-18 15:02:43.502 +0000 UTC   (27.048788,-82.178862)  3691, Parkins Terrace, North Port, Sarasota County, Florida, 34286, USA
10880218/DSC00477.JPG   2018-02-18 14:49:10.081 +0000 UTC   (27.048791,-82.178872)  3691, Parkins Terrace, North Port, Sarasota County, Florida, 34286, USA
10880218/DSC00662.JPG   2018-02-18 22:51:10.944 +0000 UTC   (27.140116,-82.443547)  914, Dartmoor Circle, Sarasota County, Florida, 34275, USA
10880218/DSC00582.JPG   2018-02-18 16:52:46.743 +0000 UTC   (27.048828,-82.178365)  3650, Parkins Terrace, North Port, Sarasota County, Florida, 34286, USA
10880218/DSC00586.JPG   2018-02-18 16:52:46.743 +0000 UTC   (27.048828,-82.178365)  3650, Parkins Terrace, North Port, Sarasota County, Florida, 34286, USA
...
```
