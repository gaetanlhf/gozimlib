# go-zim

Package `zim` implements reading support for the ZIM File Format.

Documentation at <https://godoc.org/github.com/tim-st/go-zim>.

Download and install package `zim` with `go get -u github.com/tim-st/go-zim/...`

You can download a ZIM file for testing [here](https://download.kiwix.org/zim/).

## Commands

If you want to try the `zimserver` tool, install it with `go install github.com/tim-st/go-zim/cmd/zimserver`

If you want to extract sentences or texts from a Wikipedia ZIM file use `zimtext` tool, install it with `go install github.com/tim-st/go-zim/cmd/zimtext`

`zimindex` can create an index file for a ZIM file (currently only titles are indexed), where the file size of the index file can be controlled quite good when appropriate parameters are set.

You can do full text searches (union and intersection) on the ZIM file with index file using `zimsearch`. When no index file is found, the tool searches by prefix matches.
