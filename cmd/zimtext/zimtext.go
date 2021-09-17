package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/tim-st/go-zim"
)

func main() {

	var filenameZim string
	var filenameText string
	var limit int
	var singleSentences bool
	var regexFilter string

	flag.StringVar(&filenameZim, "zim", "", "Path to the ZIM file to read from.")
	flag.StringVar(&filenameText, "txt", "", "Path to the text file, that is created or truncated if exists.")
	flag.IntVar(&limit, "limit", -1, "Stop after N lines were written (where N >= limit).")
	flag.BoolVar(&singleSentences, "sentences", false, "Only write paragraphs which are likely a single sentence.")
	flag.StringVar(&regexFilter, "regexFilter", "", "Optional Regex to define which text should be used for your language. The input text is already clean (without HTML etc). If the string is empty, all texts are used.")
	flag.Parse()

	if flag.NFlag() < 2 || len(filenameZim) == 0 || len(filenameText) == 0 {
		flag.PrintDefaults()
		return
	}

	var funcWriteText func(htmlSrc io.Reader, target *bufio.Writer, limit int) int

	if len(regexFilter) > 0 {
		if regex, errRegexCompilation := regexp.Compile(regexFilter); errRegexCompilation != nil {
			log.Fatal(errRegexCompilation)
		} else {
			funcWriteText = func(htmlSrc io.Reader, target *bufio.Writer, limit int) int {
				return WriteParagraphs(htmlSrc, target, func(p *Paragraph) bool {
					return p.IsUsableText() && regex.MatchString(p.Text)
				}, limit)
			}
		}
	} else if singleSentences {
		funcWriteText = WriteCleanSentences
	} else {
		funcWriteText = WriteCleanText
	}

	z, zimOpenErr := zim.Open(filenameZim)

	if zimOpenErr != nil {
		log.Fatal(zimOpenErr)
	}

	var txtFile, txtFileErr = os.Create(filenameText)

	if txtFileErr != nil {
		log.Fatal(txtFileErr)
	}

	var bufWriter = bufio.NewWriterSize(txtFile, 1<<22) // 4mb buffer

	var paragraphsWritten = 0

	var printProgress func(int)

	articleCount := z.ArticleCount()

	if limit > 0 {
		printProgress = func(int) {
			if paragraphsWritten%32 == 0 {
				fmt.Printf("\r%.1f%%", (float32(paragraphsWritten)/float32(limit))*100)
			}
		}
	} else {
		limit = int((^uint(0)) >> 1)
		printProgress = func(idx int) {
			if idx%32 == 0 {
				fmt.Printf("\r%.1f%%", (float32(idx)/float32(articleCount))*100)
			}
		}
	}

	for idx := uint32(0); idx < articleCount; idx++ {
		printProgress(int(idx))
		var requiredParagraphs = limit - paragraphsWritten
		if requiredParagraphs <= 0 || paragraphsWritten >= limit {
			break
		}
		entry, err := z.EntryAtURLPosition(idx)
		if err != nil {
			log.Fatal(err)
		}
		if entry.Namespace() != zim.NamespaceArticles {
			continue
		}
		blobReader, blobSize, err := z.BlobReaderAt(entry.ClusterNumber(), entry.BlobNumber())
		_ = blobSize
		if err != nil {
			continue
		}
		paragraphsWritten += funcWriteText(blobReader, bufWriter, requiredParagraphs)
	}

	bufWriter.Flush()
	txtFile.Close()
	fmt.Print("\r100.0%")
}
