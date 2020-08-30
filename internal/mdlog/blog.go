package mdlog

type Blog interface {
	GetPosts() ([]*Post, error)
	Update(post *Post) error
	CreatePost() (*Post, error)
}
