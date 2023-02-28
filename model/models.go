package model

type Comic struct {
	Number     int     `json:"num"`
	Day        int     `json:"day"`
	Month      int     `json:"month"`
	Year       int     `json:"year"`
	Title      string  `json:"title"`
	AltTitle   string  `json:"alt"`
	Transcript *string `json:"transcript"` //can be nil
	Image      string  `json:"img"`
}
