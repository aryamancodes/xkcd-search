package model

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
