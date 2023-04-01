// Functions to parse fetched explanations into stemmed comic structs

package nlp

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"xkcd/model"
)

func parseMetaData(metadata string, comic model.ExplainWikiJson) model.Comic {
	var numString, title, alt string
	for _, line := range strings.Split(metadata, "\n") {
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
		log.Fatalln(err)
	}
	title = strings.Replace(title, "title", "", 1)
	alt = strings.Replace(alt, "titletext", "", 1)

	return model.Comic{
		Num:     num,
		Title:   CleanContent(title),
		AltText: CleanContent(alt),
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
	metadataBlock := regexp.MustCompile(`(?s)\{\{.*\}\}.*==\s?Explanation`).FindString(content)
	parsedComic := parseMetaData(metadataBlock, comic)

	content = strings.Replace(content, metadataBlock, "", 1)
	transcriptBlock := regexp.MustCompile(`(?s)==\s?Transcript\s?==(.*){{comic discussion}}`).FindString(content)
	transcript, transcriptIncomplete := parseSection(transcriptBlock)

	content = strings.Replace(content, transcriptBlock, "", 1)
	explanation, explanationIncomplete := parseSection(content)

	parsedComic.Transcript = CleanContent(transcript)
	parsedComic.Explanation = CleanContent(explanation)
	parsedComic.Incomplete = transcriptIncomplete || explanationIncomplete

	return parsedComic
}
