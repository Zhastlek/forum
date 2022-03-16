package category

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/domain/category"
	"forum/internal/model"
	"log"
)

type categoryStorage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) category.CategoryStorage {
	return &categoryStorage{
		db: db,
	}
}

func (c *categoryStorage) GetAll(ctx context.Context) ([]*model.Category, error) {
	fmt.Println("Start Category storage GetAll method")
	rows, err := c.db.Query(`SELECT * FROM categories`)
	if err != nil {
		log.Println(err)
		// log.Fatal("CATEGORY STORAGE QUERY ERROR ")
		return nil, err
	}
	defer rows.Close()
	categories := []*model.Category{}

	log.Println("sign in to rows scan")
	for rows.Next() {
		category := &model.Category{}
		err := rows.Scan(&category.Id, &category.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("ERROR ErrNoRows category storage GetAll method: %v\n", err)
				break
			}
			fmt.Printf("ERROR category storage GetAll method: %v\n", err)
			return nil, err
		}
		categories = append(categories, category)
		log.Println(categories)
	}
	fmt.Println("END Category storage GetAll method")

	return categories, nil
}

func (c *categoryStorage) Create(ctx context.Context, category *model.Category) error {
	return nil
}

func (c *categoryStorage) Delete(ctx context.Context, category *model.Category) error {
	return nil
}

func (c *categoryStorage) Update(ctx context.Context, category *model.Category) error {
	return nil
}
