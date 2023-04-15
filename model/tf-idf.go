package model

// struct to store the ranking of a query, along with an array of term sections
type RankedComic struct {
	ComicNum     int
	Rank         float64
	TermSections []TermSection
}

// struct storing a term and how many times in each section the stemmed version of the term appears
type TermSection struct {
	Term             string
	TitleCount       int
	AltCount         int
	TranscriptCount  int
	ExplanationCount int
}

// struct with the number of occurences of a term in a comic, along with the number of total comic terms
type TermFreq struct {
	ComicNum        int
	TermInComicFreq map[string]int    // stemmed term -> # times term occurs in comic
	StemToRawMap    map[string]string // stemmed term -> string of raw terms with same stem
}

// struct used to store number of comics a term occurs in, for all terms, along with the number of comics
type ComicFreq struct {
	ComicsWithTermFreq map[string]int // stemmed term -> # comics with term
	TotalComics        int
}
