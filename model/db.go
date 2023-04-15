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

// struct used to store individual terms and their comic-frequencies into db
type ComicFreqDTO struct {
	Term string
	Freq int
}

func (ComicFreqDTO) TableName() string {
	return "comic_frequency"
}
