package user

import "forum/internal/model"

// Storage ...
type UserStorage interface {
	GetOne(uuid string) (*model.User, error)
	Create(user *model.User) error
	//GetAll(limit, offset int) []*model.User
	//Delete(user *model.User) error
}
