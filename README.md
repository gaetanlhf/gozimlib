# go-zim

Package `zim` implements reading support for the ZIM File Format.

Documentation at <https://godoc.org/github.com/tim-st/go-zim>.

Download and install package `zim` with `go get -u github.com/tim-st/go-zim/...`

You can download a ZIM file for testing [here](https://download.kiwix.org/zim/).

## Commands

The command above installs the tools of this package to the `$GOPATH/bin/`.

### zimserver

Tool for browsing a ZIM file in your webbrowser via an HTTP interface.

### zimindex

Tool for creating a full text index of a given ZIM file (currently only titles are indexed).

### zimsearch

Tool that lists search results for a given ZIM file and text query.
If no index file created by `zimindex` is found, a builtin prefix search is used. Otherwise the index file is used to retrieve search results sorted by score, where the search result can be calculated by union or intersection operation.

### zimtext

Tool to extract clean texts from a Wikipedia ZIM file.
Each clean HTML paragraph is written on a single line in a text file.
