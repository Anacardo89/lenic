package presentation

import (
	"html/template"
)

type Post struct {
	Id         int
	GUID       string
	Author     User
	Title      string
	RawContent string
	Content    template.HTML
	Image      string
	Date       string
	IsPublic   bool
	Rating     int
	UserRating int
	Comments   []Comment
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
