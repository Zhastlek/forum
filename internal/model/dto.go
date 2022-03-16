package model

type CreatePostDto struct {
	Title      string
	Body       string
	AuthorUUID string
	CategoryId []int
}

type UpdatePostDTO struct {
	Id         int
	Title      string
	Body       string
	CategoryId []int
	AuthorId   int
}

type UserDto struct {
	Name     string
	Email    string
	Password string
	UUID     string
}

type UpdateUUIDDTO struct {
	UUID string
	Name string
}

type CreateCommentDto struct {
	PostId            int
	CommentAuthorUUID string
	CommentBody       string
}

type SessionDto struct {
	MyUUID string
	Value  string
}
