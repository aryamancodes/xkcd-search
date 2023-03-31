// Functions to parse fetched explanations into stemmed comic structs

package nlp

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"xkcd/model"
)

func parseMetaData(metadata string) model.Comic {
	metadata = strings.TrimSuffix(metadata, "==")
	metadata = strings.TrimSpace(metadata)
	metadataSplit := strings.Split(metadata, "|")

	numString := strings.Split(metadataSplit[1], "=")[1]
	numString = regexp.MustCompile("[0-9]+").FindString(numString)
	num, err := strconv.Atoi(numString)
	if err != nil {
		log.Fatal(err)
	}

	title := strings.Split(metadataSplit[3], "=")[1]
	title = strings.TrimSpace(title)

	alt := strings.Split(metadataSplit[7], "=")[1]
	alt = strings.TrimSuffix(alt, "}}")
	alt = strings.TrimSpace(alt)

	return model.Comic{
		Num:     num,
		Title:   CleanContent(title),
		AltText: CleanContent(alt),
	}
}

// func parseTranscript(transcript string) string {

// }

func Parse(test model.ExplainWikiJson) {
	metadataBlock := regexp.MustCompile(`(?s){{.*}}..==`).FindString(test.Parse.Wikitext.Content)
	comic := parseMetaData(metadataBlock)
	// transcriptBlock :=

	// transcript := parseTranscript("explain")
	// explanation := parseTranscript("transcript")

	// metaComic.Transcript = transcript
	fmt.Printf("%+v", comic)
}
