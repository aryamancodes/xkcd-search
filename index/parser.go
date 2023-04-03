// Functions to parse fetched explanations into stemmed comic structs

package index

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"xkcd/model"
	"xkcd/nlp"
)

func parseMetaData(metadata []string, comic model.ExplainWikiJson) model.Comic {
	var numString, title, alt string
	for _, line := range metadata {
		if regexp.MustCompile(`\|\s?number\s*=`).Match([]byte(line)) {
			numString = strings.Split(line, "=")[1]
		}

		if strings.Contains(line, "| title ") {
			title = line
		}

		if strings.Contains(line, "| titletext ") {
			alt = line
			break
		}
	}
	numString = strings.TrimSpace(numString)
	num, err := strconv.Atoi(numString)
	if err != nil {
		log.Println(comic)
		log.Fatalln(err)
	}
	title = strings.Replace(title, "title", "", 1)
	alt = strings.Replace(alt, "titletext", "", 1)

	titleRaw, title := nlp.CleanAndStem(title)
	altRaw, alt := nlp.CleanAndStem(alt)

	return model.Comic{
		Num:        num,
		TitleRaw:   titleRaw,
		Title:      title,
		AltTextRaw: altRaw,
		AltText:    alt,
	}
}

func parseSection(section string) (string, bool) {
	//check for incomplete section and remove if needed
	incompleteRegex := regexp.MustCompile(`\{\{incomplete.* soon.\}\}`)
	isIncomplete := incompleteRegex.Match([]byte(section))
	if isIncomplete {
		section = incompleteRegex.ReplaceAllString(section, "")
	}
	//remove frequent sections found in the wiki sections such as headings, links and bullet points
	section = regexp.MustCompile(`(\{\{)|(\}\})|(\[\[)|(\]\])|(==)|:|\||\*`).ReplaceAllString(section, " ")
	return section, isIncomplete
}

func Parse(comic model.ExplainWikiJson) model.Comic {
	content := comic.Parse.Wikitext.Content
	//metadata is always within the first 10 lines
	metadataBlock := strings.Split(content, "\n")[:10]
	parsedComic := parseMetaData(metadataBlock, comic)

	transcriptBlock := regexp.MustCompile(`(?s)==\s?Transcript\s?==(.*){{comic discussion}}`).FindString(content)
	transcript, transcriptIncomplete := parseSection(transcriptBlock)

	content = strings.Replace(content, transcriptBlock, "", 1)
	for _, metadataSection := range metadataBlock {
		content = strings.Replace(content, metadataSection, "", 1)
	}
	explanation, explanationIncomplete := parseSection(content)

	parsedComic.TranscriptRaw, parsedComic.Transcript = nlp.CleanAndStem(transcript)
	parsedComic.ExplanationRaw, parsedComic.Explanation = nlp.CleanAndStem(explanation)
	parsedComic.Incomplete = transcriptIncomplete || explanationIncomplete

	return parsedComic
}
