package user

import (
	"database/sql"
	"forum/internal/domain/user"
	"forum/internal/model"
	"log"
)

type userStorage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) user.UserStorage {
	return &userStorage{
		db: db,
	}
}

func (us *userStorage) GetOne(element string) (*model.User, error) {
	row := us.db.QueryRow(`select * from users
		WHERE user_uuid = $1 or email = $2`, element, element)
	u := &model.User{}
	err := row.Scan(&u.Id, &u.Name, &u.Password, &u.Email, &u.UUID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			log.Println("NO ROWS")
			return nil, nil
		}
		log.Println("Error storage in GetOne user", err)
		return nil, err
	}
	return u, nil
}

//func (us *userStorage) GetAll(limit, offset int) []*model.User {
//	return nil
//}

func (us *userStorage) Create(user *model.User) error {
	_, err := us.db.Exec(`INSERT INTO users
	(login,password,email, user_uuid) 
		VALUES($1,$2,$3,$4)`, user.Name, user.Password, user.Email, user.UUID)
	if err != nil {
		log.Printf("%v\n", err)
		log.Println("Error storage in user creation")
		return err
	}
	return nil
}

//func (us *userStorage) Delete(user *model.User) error {
//	return nil
//}
