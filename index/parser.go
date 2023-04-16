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

func parseMetaData(metadata string, comic model.ExplainWikiJson) model.Comic {
	var numString, title, alt, image string
	metadata = strings.Replace(metadata, `\n`, "\n", -1)
	for _, line := range strings.Split(metadata, "\n") {
		if cleanNum := regexp.MustCompile(`\|\s?number\s*=\s*[0-9]*`).FindString(line); cleanNum != "" {
			numString = strings.Split(cleanNum, "=")[1]
		}

		if strings.Contains(line, "| title ") {
			title = line
		}

		if strings.Contains(line, "| image ") {
			image = strings.Split(line, "=")[1]
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
	image = strings.Replace(image, "image", "", 1)

	titleRaw, title := nlp.CleanAndStem(title)
	altRaw, alt := nlp.CleanAndStem(alt)
	image = strings.Replace(image, " ", "", -1)

	return model.Comic{
		Num:        num,
		ImageName:  image,
		TitleRaw:   titleRaw,
		Title:      title,
		AltTextRaw: altRaw,
		AltText:    alt,
	}

}

func parseSection(section string) (string, bool) {
	//check for incomplete section and remove if needed
	incompleteRegex := regexp.MustCompile(`\{\{(?s)incomplete.*?\}\}`)
	isIncomplete := incompleteRegex.Match([]byte(section))
	if isIncomplete {
		section = incompleteRegex.ReplaceAllString(section, "")
	}
	//remove frequent sections found in the wiki sections such as headings, links and bullet points
	section = regexp.MustCompile(`(\[http[\S]+)(\{\{)|(\}\})|(\[\[)|(\]\])|(==)|:|\||\*`).ReplaceAllString(section, " ")
	return section, isIncomplete
}

func Parse(comic model.ExplainWikiJson) model.Comic {
	content := comic.Parse.Wikitext.Content
	//metadata is always before the first (== prefixed) heading
	metadataBlock := strings.Split(content, "==")[0]
	parsedComic := parseMetaData(metadataBlock, comic)

	transcriptBlock := regexp.MustCompile(`(?s)==\s?Transcript\s?==(.*){{comic discussion}}`).FindString(content)
	transcript, transcriptIncomplete := parseSection(transcriptBlock)

	content = strings.Replace(content, transcriptBlock, "", 1)
	content = strings.Replace(content, metadataBlock, "", 1)
	explanation, explanationIncomplete := parseSection(content)

	parsedComic.TranscriptRaw, parsedComic.Transcript = nlp.CleanAndStem(transcript)
	parsedComic.ExplanationRaw, parsedComic.Explanation = nlp.CleanAndStem(explanation)
	parsedComic.Incomplete = transcriptIncomplete || explanationIncomplete

	return parsedComic
}
