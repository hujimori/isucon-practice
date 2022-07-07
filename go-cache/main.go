package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID          int       `db:"id"`
	AccountName string    `db:"account_name"`
	Passhash    string    `db:"passhash"`
	Authority   int       `db:"authority"`
	DelFlg      int       `db:"del_flg"`
	CreatedAt   time.Time `db:"created_at"`
}

type Post struct {
	ID           int       `db:"id"`
	UserID       int       `db:"user_id"`
	Imgdata      []byte    `db:"imgdata"`
	Body         string    `db:"body"`
	Mime         string    `db:"mime"`
	CreatedAt    time.Time `db:"created_at"`
	CommentCount int       `db:"comment_count"`
	Comments     []Comment
	User         User
	CSRFToken    string
}

type Comment struct {
	ID        int       `db:"id"`
	PostID    int       `db:"post_id"`
	UserID    int       `db:"user_id"`
	Comment   string    `db:"comment"`
	CreatedAt time.Time `db:"created_at"`
	User      User
}

var db *sqlx.DB
var rd *redis.Client

func main() {

	var err error
	db, err = sqlx.Open("mysql", "root:root@tcp(127.0.0.1:3306)/isuconp?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	rd = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	results := []Post{}
	err = db.Select(&results, "SELECT `id`, `user_id`, `body`, `mime`, `created_at` FROM `posts` ORDER BY `created_at` DESC LIMIT 30")
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range results {
		getUser(p.UserID)
	}
}

func getUser(id int) {
	// user := User{}

	var ctx = context.Background()

	it, err := rd.Get(ctx, strconv.Itoa(id)).Result()
	if err == nil {
		fmt.Printf("%v", it)
	}
}
