package models

import (
	"github.com/kwanok/friday/config"
	"log"
)

type Post struct {
	Id        uint64
	AuthorId  uint64
	Title     string
	Content   string
	CreatedAt string
	UpdatedAt string
}

func GetAllPosts() ([]Post, error) {
	posts := make([]Post, 0)

	rows, err := config.DBCon.Query("SELECT id, author_id, title, content, created_at, updated_at FROM posts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.Id, &post.AuthorId, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}
