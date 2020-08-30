package blogger

import (
	"context"

	"github.com/KimMachineGun/mdlog/internal/mdlog"

	"golang.org/x/oauth2"
	bl "google.golang.org/api/blogger/v3"
	"google.golang.org/api/option"
)

type Blog struct {
	service *bl.Service
	blogID  string
}

func New(ctx context.Context, url string, tokenSource oauth2.TokenSource) (*Blog, error) {
	service, err := bl.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, err
	}

	b, err := service.Blogs.GetByUrl(url).Do()
	if err != nil {
		return nil, err
	}

	return &Blog{
		service: service,
		blogID:  b.Id,
	}, nil
}

func (b *Blog) GetPosts() ([]*mdlog.Post, error) {
	pl, err := b.service.Posts.List(b.blogID).Status("LIVE", "DRAFT").Do()
	if err != nil {
		return nil, err
	}

	posts := make([]*mdlog.Post, 0, len(pl.Items))
	for _, post := range pl.Items {
		status := mdlog.PostStatusDraft
		if post.Status == "LIVE" {
			status = mdlog.PostStatusPublished
		}
		posts = append(posts, &mdlog.Post{
			ID:      post.Id,
			Title:   post.Title,
			Content: post.Content,
			Tags:    post.Labels,
			Status:  status,
		})
	}

	return posts, nil
}

func (b *Blog) Update(post *mdlog.Post) error {
	err := post.Validate()
	if err != nil {
		return err
	}

	p, err := b.service.Posts.Update(b.blogID, post.ID, &bl.Post{
		Title:   post.Title,
		Content: post.Content,
		Labels:  post.Tags,
	}).Do()
	if err != nil {
		return err
	}

	status := "DRAFT"
	if post.Status == mdlog.PostStatusPublished {
		status = "LIVE"
	}

	if p.Status != status {
		if status == "LIVE" {
			_, err = b.service.Posts.Publish(b.blogID, post.ID).Do()
		} else {
			_, err = b.service.Posts.Revert(b.blogID, post.ID).Do()
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (b Blog) CreatePost() (*mdlog.Post, error) {
	p, err := b.service.Posts.Insert(b.blogID, &bl.Post{}).Do()
	if err != nil {
		return nil, err
	}

	_, err = b.service.Posts.Revert(b.blogID, p.Id).Do()
	if err != nil {
		return nil, err
	}

	return &mdlog.Post{
		ID: p.Id,
	}, nil
}
