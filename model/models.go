package model

type Comic struct {
	Number     int     `json:"num,"`
	Day        int     `json:"day,string"`
	Month      int     `json:"month,string"`
	Year       int     `json:"year,string"`
	Title      string  `json:"title"`
	AltTitle   string  `json:"alt"`
	Transcript *string `json:"transcript"` //can be nil
	Image      string  `json:"img"`
}
