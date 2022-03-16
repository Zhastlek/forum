package comment

import (
	"database/sql"
	"fmt"
	"forum/internal/domain/comment"
	"forum/internal/model"
	"log"
	"time"
)

type commentStorage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) comment.CommentStorage {
	return &commentStorage{
		db: db,
	}
}

func (c *commentStorage) Create(comment *model.Comment) error {
	t := time.Now()
	fmt.Println("Start Comment storage Create method")
	result, err := c.db.Exec(`INSERT INTO comments
	(user_uuid, post_id, date, comment)
	VALUES ($1,$2,$3,$4)`,
		comment.AuthorUUID, comment.PostId, t, comment.Body)
	if err != nil {
		log.Printf("ERROR comment storage Create method db.EXEC func:--> %v\n", err)
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("ERROR comment storage Create method LastInsert func:--> %v\n", err)
		return err
	}
	fmt.Println(result.LastInsertId())
	fmt.Println(result.RowsAffected())
	fmt.Println("ID create comment", id)
	fmt.Println("END comment storage Create method")
	return nil
}

func (c *commentStorage) GetAll(postId int) ([]*model.Comment, error) {
	fmt.Println("Start Comment storage GetAll method")
	rows, err := c.db.Query(`SELECT comments.id, comments.user_uuid, comments.post_id, comments.date, comments.comment, users.login 
		FROM comments INNER JOIN users on comments.user_uuid = users.user_uuid
		WHERE comments.post_id = $1
		ORDER BY comments.id DESC`, postId)
	if err != nil {
		log.Println(err)
		// log.Fatal("Comment STORAGE  GetAll QUERY ERROR ")
		return nil, err
	}
	defer rows.Close()
	comments := []*model.Comment{}

	for rows.Next() {
		oneComment := &model.Comment{}
		err := rows.Scan(&oneComment.Id, &oneComment.AuthorUUID, &oneComment.PostId, &oneComment.Date, &oneComment.Body, &oneComment.AuthorName)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("ERROR Comment storage GetAll method :--> %v\n", err)
				return nil, err
			}
			log.Printf("ERROR Comment storage GetAll method :--> %v\n", err)
			return nil, err
		}
		comments = append(comments, oneComment)
	}
	fmt.Println("END Comment storage GetAll method")
	return comments, nil
}

func (c *commentStorage) Delete(comment *model.Comment) error {
	_, err := c.db.Exec(`DELETE FROM comments
	WHERE user_uuid = $1 and post_id = $2 and comment = $3`,
		comment.AuthorUUID, comment.PostId, comment.Body)
	if err != nil {
		log.Printf("ERROR post storage Delete method: %v\n", err)
		return err
	}
	return nil
}

// func (c *commentStorage) Update(comment *model.Comment) error {
// 	return nil
// }

func (ps *commentStorage) GetIdPosts(userUUID string) ([]int, error) {
	log.Println("Start COMMENT storage Get Id Posts method")
	rows, err := ps.db.Query(`SELECT DISTINCT post_id FROM comments
		WHERE user_uuid = $1
		ORDER BY post_id DESC`, userUUID)
	if err != nil {
		log.Println(err)
		// log.Fatal("COMMENT STORAGE Get Id Posts QUERY ERROR ")
		return nil, err
	}
	defer rows.Close()
	var listId []int

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("ERROR COMMENT storage Get Id Posts method :--> %v\n", err)
				return nil, err
			}
			log.Printf("ERROR COMMENT storage Get Id Posts method:--> %v\n", err)
			return nil, err
		}
		listId = append(listId, id)
	}
	log.Println("END COMMENT storage Get Id Posts method")
	return listId, nil
}
