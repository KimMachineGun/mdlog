package mdlog

import (
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

var (
	ErrUpToDate = errors.New("local posts is up to date with remote posts")
)

type Service struct {
	blog Blog
}

func NewService(blog Blog) *Service {
	return &Service{
		blog: blog,
	}
}

func (s *Service) Plan(localPosts []*Post) (map[*Post]*Post, error) {
	remotePosts, err := s.blog.GetPosts()
	if err != nil {
		return nil, fmt.Errorf("cannot get remote posts: %v", err)
	}

	remotePostsMap := map[string]*Post{}
	for _, remotePost := range remotePosts {
		remotePostsMap[remotePost.ID] = remotePost
	}

	plan := map[*Post]*Post{}
	for _, localPost := range localPosts {
		remotePost := remotePostsMap[localPost.ID]
		if remotePost == nil {
			return nil, fmt.Errorf("unregistered local post: %s[%s]", localPost.Title, localPost.ID)
		}

		if !localPost.Equal(remotePost) {
			plan[remotePost] = localPost
		}
	}

	return plan, nil
}

func (s *Service) Sync(localPosts []*Post) (map[*Post]*Post, error) {
	plan, err := s.Plan(localPosts)
	if err != nil {
		return nil, fmt.Errorf("cannot make sync plan: %v", err)
	}

	if len(plan) == 0 {
		return nil, ErrUpToDate
	}

	g := errgroup.Group{}
	for _, localPost := range plan {
		localPost := localPost
		g.Go(func() error {
			return s.blog.Update(localPost)
		})
	}
	err = g.Wait()
	if err != nil {
		g = errgroup.Group{}
		for remotePost, _ := range plan {
			remotePost := remotePost
			g.Go(func() error {
				return s.blog.Update(remotePost)
			})
		}
		rollbackErr := g.Wait()
		if rollbackErr != nil {
			return nil, fmt.Errorf("error on roll back: %v (original err: %v)", rollbackErr, err)
		}

		return nil, fmt.Errorf("cannot sync with remote posts: %v", err)
	}

	return plan, nil
}

func (s *Service) Create() (*Post, error) {
	post, err := s.blog.CreatePost()
	if err != nil {
		return nil, fmt.Errorf("cannot create post: %v", err)
	}

	return post, nil
}
