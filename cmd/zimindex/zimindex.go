package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/tim-st/go-minsearch"
	"github.com/tim-st/go-zim"
)

func main() {

	var filename string
	var idLimit int
	var noSync bool

	flag.StringVar(&filename, "filename", "", "Filename of the ZIM file to index.")
	flag.IntVar(&idLimit, "idLimit", -1, "If idLimit>0 only the highest idLimit scores will be indexed per key.")
	flag.BoolVar(&noSync, "noSync", false, "If nosync=true indexing will be much faster but data can be lost if system crashes.")
	flag.Parse()

	if flag.NFlag() < 1 || len(filename) == 0 {
		flag.PrintDefaults()
		return
	}

	z, zimErr := zim.Open(filename)
	if zimErr != nil {
		log.Fatal(zimErr)
	}

	index, openErr := minsearch.Open(z.Filename()+".idx", noSync)

	if openErr != nil {
		log.Fatal(openErr)
	}

	var currentPos uint32

	if lastID, lastIDErr := index.LastID(); lastIDErr == nil {
		currentPos = lastID
	} else {

		var found bool
		_, currentPos, found = z.EntryWithNamespace(zim.NamespaceArticles)

		if !found {
			log.Fatal("zimindex: first entry to index not found")
		}

	}

	var batchPairs []minsearch.Pair

	for currentPos < z.ArticleCount() {
		entry, err := z.EntryAtURLPosition(currentPos)
		if err != nil {
			break
		}
		if entry.Namespace() != zim.NamespaceArticles {
			break
		}

		if entry.IsArticle() || entry.IsRedirect() {

			if currentPos%16384 == 0 {
				urlPrefix := entry.URL()
				if len(urlPrefix) > 6 {
					urlPrefix = urlPrefix[:6]
				}
				fmt.Printf("\rIndexing Directory Entry at position %d (URL-Prefix: %s)...", currentPos, urlPrefix)
			}

			batchPairs = append(batchPairs, minsearch.Pair{ID: minsearch.ID(currentPos), Text: entry.Title()})

			// TODO: support content indexing
			/*
				if entry.IsArticle() {
					reader, blobSize, err := z.BlobReader(&entry)
					if err == nil {
						var buffer = make([]byte, blobSize)
						_, err := reader.Read(buffer[:])
						if err == nil {
							batchPairs = append(batchPairs, minsearch.Pair{ID: minsearch.ID(currentPos), Text: buffer})
						}
					}
				}
			*/

		}
		currentPos++

		if len(batchPairs) >= 1000 {
			index.IndexBatch(batchPairs, idLimit)
			// use currentPos instead of ID (both uint32)
			if err := index.SetLastID(currentPos); err != nil {
				log.Fatal(err)
			}
			batchPairs = batchPairs[:0]
		}

	}

	index.SetLastID(currentPos)
	index.IndexBatch(batchPairs, idLimit)

	if updateErr := index.UpdateStatistics(); updateErr != nil {
		log.Fatal(updateErr)
	}

	fmt.Println("\rFinished!")

}
