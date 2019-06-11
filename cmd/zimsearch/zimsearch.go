package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tim-st/go-minsearch"
	"github.com/tim-st/go-zim"
)

func main() {

	var filename string
	var query string
	var limit int
	var intersection bool

	flag.StringVar(&filename, "filename", "", "Filename of the ZIM file to use.")
	flag.StringVar(&query, "query", "", "The text to search in the ZIM index file.")
	flag.IntVar(&limit, "limit", -1, "Limit the output of the result to the given number.")
	flag.BoolVar(&intersection, "intersection", false, "true = intersection set; false = union set")
	flag.Parse()

	if flag.NFlag() < 2 || len(filename) == 0 || len(query) == 0 {
		flag.PrintDefaults()
		return
	}

	z, zimErr := zim.Open(filename)
	if zimErr != nil {
		log.Fatal(zimErr)
	}

	idxPath := z.Filename() + ".idx"
	if f, fErr := os.Open(idxPath); fErr != nil {
		// index file doesn't exist; use builtin prefix search
		if limit < 1 {
			limit = 100
		}
		suggestions := z.EntriesWithSimilarity(zim.NamespaceArticles, []byte(query), limit)
		for idx, suggestion := range suggestions {
			title := string(suggestion.Title())
			url := string(suggestion.URL())
			fmt.Println(idx, title, url)
		}

	} else {
		f.Close()
		if index, openErr := minsearch.Open(idxPath, true); openErr == nil {
			start := time.Now()
			var setOp = minsearch.Union
			if intersection {
				setOp = minsearch.Intersection
			}
			queryResults, queryErr := index.Search([]byte(query), setOp, 0)
			fmt.Printf("Took: %s\n", time.Since(start))

			if queryErr != nil {
				log.Fatal(queryErr)
			}

			for idx, result := range queryResults {
				if limit > 0 && idx == limit {
					break
				}
				entity, entryErr := z.EntryAtURLPosition(result.ID)
				if entryErr != nil {
					continue
				}
				title := string(entity.Title())
				url := string(entity.URL())
				score := result.Score
				fmt.Println(idx, title, url, score)
			}

		}
	}

}
