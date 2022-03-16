package model

type User struct {
	Id       int
	Name     string
	Password string
	Email    string
	UUID     string
}

type Post struct {
	Id         int
	Title      string
	Body       string
	AuthorId   int
	AuthorName string
	AuthorUUID string
	CategoryId []int
	Date       string
	Likes      *PostLikeAndDislike
	Comments   []*Comment
}

type PostLikeAndDislike struct {
	Id         int
	AuthorUUID string
	PostId     int
	Reaction   string
	Like       int
	DisLike    int
	MyReaction bool
	VoitStatus string
}

type CommentLikeAndDislike struct {
	Id         int
	AuthorUUID string
	CommentId  int
	Like       int
	DisLike    int
	Reaction   string
	MyReaction bool
	VoitStatus string
}

type Comment struct {
	Id         int
	PostId     int
	AuthorUUID string
	AuthorName string
	Body       string
	MyComment  bool
	Date       string
	Likes      *CommentLikeAndDislike
}

type Session struct {
	Id          int
	SessionUUID string
	Date        string
	UserUUID    string
}

type Category struct {
	Id   int
	Name string
}

type Like struct {
	PostLike    *PostLikeAndDislike
	CommentLike *CommentLikeAndDislike
}
