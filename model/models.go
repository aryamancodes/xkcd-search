package model

// struct for Gorm to map struct to table name
type Tabler interface {
	TableName() string
}

// comic struct used for indexing
type Comic struct {
	Num         int
	Title       string
	AltText     string
	Transcript  string
	Explanation string
	Incomplete  bool `gorm:"default:false"`
}

// struct to store the ranking of a query
type RankedComic struct {
	ComicNum int
	Rank     float64
}

// struct with the number of occurences of a term in a comic, along with the number of total comic terms
type TermFreq struct {
	ComicNum        int
	TermInComicFreq map[string]int // term -> # times term occurs in comic
	TotalTerms      int
}

// struct used to store individual terms and their term-frequencies into db
type TermFreqDTO struct {
	ComicNum int
	Term     string
	Freq     int
}

func (TermFreqDTO) TableName() string {
	return "term_frequency"
}

// struct used to store number of comics a term occurs in, for all terms, along with the number of comics
type ComicFreq struct {
	ComicsWithTermFreq map[string]int // term -> # comics with term
	TotalComics        int
}

// struct used to store individual terms and their comic-frequencies into db
type ComicFreqDTO struct {
	Id   int `gorm:"autoIncrement"`
	Term string
	Freq int
}

func (ComicFreqDTO) TableName() string {
	return "comic_frequency"
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
