package mdlog

import (
	"errors"
	"fmt"

	"github.com/KimMachineGun/mdlog/internal/md2html"
)

const (
	PostStatusDraft PostStatus = iota
	PostStatusPublished
)

var (
	postStatusText = map[PostStatus]string{
		PostStatusDraft:     "DRAFT",
		PostStatusPublished: "PUBLISHED",
	}
)

type PostStatus int

func PostStatusText(status PostStatus) string {
	return postStatusText[status]
}

type Post struct {
	ID      string
	Title   string
	Content string
	Tags    []string
	Status  PostStatus
}

func (p *Post) Validate() error {
	if p.ID == "" {
		return errors.New("id is required field")
	}
	if _, ok := postStatusText[p.Status]; !ok {
		return errors.New("invalid post status")
	}

	return nil
}

func (p *Post) Equal(p2 *Post) bool {
	if p == p2 {
		return true
	}

	if p.ID != p2.ID {
		return false
	}
	if p.Title != p2.Title {
		return false
	}
	if p.Content != p2.Content {
		return false
	}
	if len(p.Tags) != len(p2.Tags) {
		return false
	}
	for idx, tag := range p.Tags {
		if tag != p2.Tags[idx] {
			return false
		}
	}
	if p.Status != p2.Status {
		return false
	}

	return true
}

func ParsePost(md []byte) (*Post, error) {
	html, err := md2html.Convert(md)
	if err != nil {
		return nil, err
	}

	post := Post{
		Content: html.Content,
	}
	if id, ok := html.Meta["id"]; ok {
		post.ID = fmt.Sprint(id)
	}
	if title, ok := html.Meta["title"]; ok {
		post.Title = fmt.Sprint(title)
	}
	if tags, ok := html.Meta["tags"].([]interface{}); ok && len(tags) > 0 {
		post.Tags = make([]string, 0, len(tags))
		for _, tag := range tags {
			post.Tags = append(post.Tags, fmt.Sprint(tag))
		}
	}
	if status, ok := html.Meta["status"].(int); ok {
		post.Status = PostStatus(status)
	}

	err = post.Validate()
	if err != nil {
		return nil, err
	}

	return &post, nil
}
