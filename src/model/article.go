package model

import "time"

type Article struct {
	Url          string     `json:"url"`
	Slug         string     `json:"slug"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	BodyMarkdown *string    `json:"body_markdown"`
	BodyHtml     *string    `json:"body_html"`
	Thumbnail    string     `json:"thumbnail"`
	Author       Author     `json:"author"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	Tags         []string   `json:"tags"`
}

type Author struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
