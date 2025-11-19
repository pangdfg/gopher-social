package db

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/pangdfg/gopher-social/internal/store"
	"gorm.io/gorm"
)

func Seed(ctx context.Context, db *gorm.DB) {
	gofakeit.Seed(time.Now().UnixNano())

	db.Exec("DELETE FROM comments")
	db.Exec("DELETE FROM posts")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")

	role := store.Role{Name: "user", Description: "Regular user", Level: 1}
	if err := db.Create(&role).Error; err != nil {
		log.Println("Error creating role:", err)
		return
	}

	users := generateUsers(50, role)
	if err := db.Create(&users).Error; err != nil {
		log.Println("Error creating users:", err)
		return
	}

	posts := generatePosts(10, users)
	if err := db.Create(&posts).Error; err != nil {
		log.Println("Error creating posts:", err)
		return
	}

	comments := generateComments(20, users, posts)
	if err := db.Create(&comments).Error; err != nil {
		log.Println("Error creating comments:", err)
		return
	}

	log.Println("Seeding complete")
}

func generateUsers(num int, role store.Role) []*store.User {
	users := make([]*store.User, num)
	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Role:     role,
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	convertToTags := func(words []string) []store.Tag {
		tags := make([]store.Tag, len(words))
		for i, word := range words {
			tags[i] = store.Tag{Name: word}
		}
		return tags
	}

	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		randomWords := []string{gofakeit.Word(), gofakeit.Word()}

		posts[i] = &store.Post{
			UserID:  uint(user.ID),
			Title:   gofakeit.Sentence(5),
			Content: gofakeit.Paragraph(1, 3, 10, " "),
			Tags:    convertToTags(randomWords),
		}
	}

	return posts
}


func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID:  uint(posts[rand.Intn(len(posts))].ID),
			UserID:  uint(posts[rand.Intn(len(posts))].ID),
			Content: gofakeit.Sentence(8),
		}
	}
	return cms
}
