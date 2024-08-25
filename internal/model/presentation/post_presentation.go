package presentation

import (
	"html/template"
)

type Post struct {
	Id         int
	GUID       string
	Author     string
	Title      string
	RawContent string
	Content    template.HTML
	Image      string
	Date       string
	IsPublic   int
	Rating     int
	UserRating int
	Comments   []Comment
	Session    Session
}

func (p Post) TruncatedText() string {
	chars := 0
	for i := range p.RawContent {
		chars++
		if chars > 150 {
			return p.RawContent[:i] + `...`
		}
	}
	return p.RawContent
}
