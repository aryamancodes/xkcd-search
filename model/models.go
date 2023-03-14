package model

// comic struct used for indexing
type Comic struct {
	Number      int
	Date        string
	Title       string
	Transcript  string
	Explanation string
	Image       string
}

// struct for fetching the most recent comic
type CurrentComicJson struct {
	Number int `json:"num"`
}

// struct for nested json returned by the explain xkcd api
type ExplainWikiJson struct {
	Parse ParseStruct `json:"parse"`
}

type ParseStruct struct {
	Title    string         `json:"title"`
	Wikitext WikitextStruct `json:"wikitext"`
}

type WikitextStruct struct {
	Content string `json:"*"`
}
