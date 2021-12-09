package models

import (
	"friday/server"
	"friday/server/utils"
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

	rows, err := server.DBCon.Query("SELECT id, author_id, title, content, created_at, updated_at FROM posts")
	utils.FatalError{Error: err}.Handle()
	defer rows.Close()

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.Id, &post.AuthorId, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
		utils.FatalError{Error: err}.Handle()
		posts = append(posts, post)
	}

	return posts, nil
}
