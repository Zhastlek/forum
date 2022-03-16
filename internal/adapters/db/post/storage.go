package post

import (
	"database/sql"
	"fmt"
	"forum/internal/domain/post"
	"forum/internal/model"
	"log"
	"time"
)

type postStorage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) post.PostStorage {
	return &postStorage{
		db: db,
	}
}

func (ps *postStorage) GetOne(id int) ([]*model.Post, error) {
	fmt.Println("Start Post storage GetOne method")
	rows, err := ps.db.Query(`SELECT posts.id, posts.title, posts.body, posts.date, posts.user_uuid, users.login
		FROM posts INNER JOIN users on posts.user_uuid = users.user_uuid
		WHERE posts.id=$1`, id)
	if err != nil {
		log.Println(err)
		// log.Fatal("POST STORAGE QUERY ERROR ")
		return nil, err
	}
	defer rows.Close()
	posts := []*model.Post{}

	for rows.Next() {
		onePost := &model.Post{}
		err := rows.Scan(&onePost.Id, &onePost.Title, &onePost.Body, &onePost.Date, &onePost.AuthorUUID, &onePost.AuthorName)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("ERROR post storage GetOne method :--> %v\n", err)
				return nil, err
			}
			log.Printf("ERROR post storage GetOne method :--> %v\n", err)
			return nil, err
		}
		rowsCategory, err := ps.db.Query(`SELECT category_id FROM post_category 
		WHERE post_id = $1
		ORDER BY id DESC`, onePost.Id)
		if err != nil {
			log.Println(err)
			// log.Fatal("POST STORAGE QUERY ERROR ")
			return nil, err
		}
		defer rows.Close()
		categories := []int{}

		for rowsCategory.Next() {
			var numCategory int
			err := rowsCategory.Scan(&numCategory)
			if err != nil {
				if err == sql.ErrNoRows {
					log.Printf("ERROR post storage GetAll method  ErrNoRows:--> %v\n", err)
					return nil, err
				}
				log.Printf("ERROR post storage GetAll method :--> %v\n", err)
				return nil, err
			}
			categories = append(categories, numCategory)
		}
		onePost.CategoryId = categories
		posts = append(posts, onePost)
	}
	fmt.Println("END Post storage GetOne method")
	return posts, nil
}

func (ps *postStorage) GetAll() ([]*model.Post, error) {
	fmt.Println("Start Post storage GetAll method")
	rows, err := ps.db.Query(`SELECT posts.id, posts.title, posts.body, posts.date, posts.user_uuid, users.login 
		FROM posts INNER JOIN users on posts.user_uuid = users.user_uuid
		ORDER BY posts.id DESC`)
	if err != nil {
		log.Println(err)
		// log.Fatal("POST STORAGE QUERY ERROR ")
		return nil, err
	}
	defer rows.Close()
	posts := []*model.Post{}

	for rows.Next() {
		onePost := &model.Post{}
		err := rows.Scan(&onePost.Id, &onePost.Title, &onePost.Body, &onePost.Date, &onePost.AuthorUUID, &onePost.AuthorName)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("ERROR post storage GetAll method  ErrNoRows:--> %v\n", err)
				return nil, err
			}
			log.Printf("ERROR post storage GetAll method :--> %v\n", err)
			return nil, err
		}
		rowsCategory, err := ps.db.Query(`SELECT category_id FROM post_category 
		WHERE post_id = $1
		ORDER BY id DESC`, onePost.Id)
		if err != nil {
			log.Println(err)
			// log.Fatal("POST STORAGE QUERY ERROR ")
			return nil, err
		}
		defer rows.Close()
		categories := []int{}

		for rowsCategory.Next() {
			var numCategory int
			err := rowsCategory.Scan(&numCategory)
			if err != nil {
				if err == sql.ErrNoRows {
					log.Printf("ERROR post storage GetAll method  ErrNoRows:--> %v\n", err)
					return nil, err
				}
				log.Printf("ERROR post storage GetAll method :--> %v\n", err)
				return nil, err
			}
			categories = append(categories, numCategory)
		}
		onePost.CategoryId = categories
		posts = append(posts, onePost)
	}

	fmt.Println("END Post storage GetAll method")
	return posts, nil
}

func (ps *postStorage) Create(post *model.Post) (int, error) {
	t := time.Now()
	fmt.Println("Start Post storage Create method")
	result, err := ps.db.Exec(`INSERT INTO posts
	(title,body,date, user_uuid)
	VALUES ($1,$2,$3,$4)`,
		post.Title, post.Body, t, post.AuthorUUID)
	if err != nil {
		log.Printf("ERROR post storage Create method db.EXEC func:--> %v\n", err)
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("ERROR post storage Create method LastInsert func:--> %v\n", err)
		return 0, err
	}
	fmt.Println(result.LastInsertId())
	fmt.Println(result.RowsAffected())
	fmt.Println("END Post storage Create method")
	return int(id), nil
}

func (ps *postStorage) CreatePostCategory(postId int, categoryId int) error {
	fmt.Println("Start Post storage CreatePostCategory method")
	_, err := ps.db.Exec(`INSERT INTO post_category
	(post_id,category_id)
	VALUES ($1,$2)`,
		postId, categoryId)
	if err != nil {
		log.Printf("ERROR post storage CreatePostCategory method db.EXEC func:--> %v\n", err)
		return err
	}

	fmt.Println("END Post storage CreatePostCategory method")
	return nil
}

func (ps *postStorage) DeletePostCategory(postId int) error {
	_, err := ps.db.Exec(`DELETE FROM post_category
		WHERE post_id=$1`, postId)
	if err != nil {
		log.Printf("ERROR post storage DeletePostCategory method: %v\n", err)
		return err
	}
	return nil
}

func (ps *postStorage) Delete(id int) error {
	_, err := ps.db.Exec(`DELETE FROM posts
		WHERE id=$1`, id)
	if err != nil {
		log.Printf("ERROR post storage Delete method: %v\n", err)
		return err
	}
	return nil
}

func (ps *postStorage) Update(post *model.Post) error {
	_, err := ps.db.Exec(`UPDATE posts
	SET title=$1,
	body=$2
	WHERE id = $3`, post.Title, post.Body, post.Id)
	if err != nil {
		log.Printf("ERROR post storage update func :--->%v\n", err)
		return err
	}
	return nil
}

func (ps *postStorage) GetIdPostsFromCategory(categoryId int) ([]int, error) {
	log.Println("Start Post storage  Get Id Posts From Category  method")
	rows, err := ps.db.Query(`SELECT post_id FROM post_category
		WHERE category_id = $1`, categoryId)
	if err != nil {
		log.Println(err)
		// log.Fatal("POST STORAGE Get Id Posts From Category  QUERY ERROR ")
		return nil, err
	}
	defer rows.Close()
	var listId []int

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("ERROR post storage Get Id Posts From Category method :--> %v\n", err)
				return nil, err
			}
			log.Printf("ERROR post storage Get Id Posts From Category method:--> %v\n", err)
			return nil, err
		}
		listId = append(listId, id)
	}
	fmt.Println("END Post storage Get Id Posts From Category  method")
	return listId, nil
}

func (ps *postStorage) SortedByCategory(categoryId int) ([]*model.Post, error) {
	fmt.Println("Start Post storage SortedByCategory method")
	rows, err := ps.db.Query(`SELECT posts.id, posts.title, posts.body, posts.date, posts.user_uuid, users.login
		FROM post_category INNER JOIN posts on post_category.post_id = posts.id
		INNER JOIN users on posts.user_uuid = users.user_uuid
		WHERE post_category.category_id = $1
		ORDER BY posts.id DESC`, categoryId)
	if err != nil {
		log.Println(err)
		// log.Fatal("POST STORAGE SortedByCategory QUERY ERROR ")
		return nil, err
	}
	defer rows.Close()
	posts := []*model.Post{}

	for rows.Next() {
		onePost := &model.Post{}
		err := rows.Scan(&onePost.Id, &onePost.Title, &onePost.Body, &onePost.Date, &onePost.AuthorUUID, &onePost.AuthorName)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("ERROR post storage SortedByCategory method :--> %v\n", err)
				return nil, err
			}
			log.Printf("ERROR post storage SortedByCategory method :--> %v\n", err)
			return nil, err
		}
		posts = append(posts, onePost)
	}
	fmt.Println("END Post storage SortedByCategory method")
	return posts, nil
}

func (ps *postStorage) GetIdPosts(userUUID string) ([]int, error) {
	log.Println("Start Post storage Get Id Posts method")
	rows, err := ps.db.Query(`SELECT id FROM posts
		WHERE user_uuid = $1
		ORDER BY id DESC`, userUUID)
	if err != nil {
		log.Println(err)
		// log.Fatal("POST STORAGE Get Id Posts QUERY ERROR ")
		return nil, err
	}
	defer rows.Close()
	var listId []int

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("ERROR post storage Get Id Posts method :--> %v\n", err)
				return nil, err
			}
			log.Printf("ERROR post storage Get Id Posts method:--> %v\n", err)
			return nil, err
		}
		listId = append(listId, id)
	}
	fmt.Println("END Post storage Get Id Posts method")
	return listId, nil
}
