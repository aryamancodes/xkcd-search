package model

// struct for Gorm to map struct to table name
type Tabler interface {
	TableName() string
}

// comic struct used for indexing. The raw versions of fields are
// only cleaned whereas the non-raw version are cleaned and stemmed
type Comic struct {
	Num            int
	ImageName      string
	Title          string
	TitleRaw       string
	AltText        string
	AltTextRaw     string
	Transcript     string
	TranscriptRaw  string
	Explanation    string
	ExplanationRaw string
	Incomplete     bool `gorm:"default:false"`
}

// struct to store the ranking of a query
type RankedComic struct {
	ComicNum int
	Rank     float64
}

// struct with the number of occurences of a term in a comic, along with the number of total comic terms
type TermFreq struct {
	ComicNum        int
	TermInComicFreq map[string]int    // stemmed term -> # times term occurs in comic
	StemToRawMap    map[string]string // stemmed term -> string of raw terms with same stem
	TotalTerms      int
}

// struct used to store individual terms and their term-frequencies into db
type TermFreqDTO struct {
	ComicNum int
	Term     string
	TermsRaw string
	Freq     int
}

func (TermFreqDTO) TableName() string {
	return "term_frequency"
}

// struct used to store number of comics a term occurs in, for all terms, along with the number of comics
type ComicFreq struct {
	ComicsWithTermFreq map[string]int // stemmed term -> # comics with term
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
