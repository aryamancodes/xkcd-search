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

// struct with the number of occurences of a term in a comic, along with the number of total comicß terms
type TermFreq struct {
	TermInComicFreq map[string]int // term -> # times term occurs in comic
	TotalTerms      int
}

// struct used to store number of comics a term occurs in, for all terms, along with the number of comics
type ComicFreq struct {
	ComicsWithTermFreq map[string]int // term -> # comics with term
	TotalComics        int
}

// struct for fetching the most recent comic
type CurrentComicJson struct {
	Number int `json:"num"`
}

// structs for nested json returned by the explain xkcd api
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
