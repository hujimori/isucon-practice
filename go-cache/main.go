package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID          int       `db:"id" redis:"id"`
	AccountName string    `db:"account_name" redis:"account_name"`
	Passhash    string    `db:"passhash" redis:"passhash"`
	Authority   int       `db:"authority" redis:"authority"`
	DelFlg      int       `db:"del_flg" redis:"del_flg"`
	CreatedAt   time.Time `db:"created_at" redis:"created_at"`
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
var rdb *redis.Client

func main() {

	var err error
	db, err = sqlx.Open("mysql", "root:root@tcp(127.0.0.1:3306)/isuconp?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	results := []Post{}
	err = db.Select(&results, "SELECT `id`, `user_id`, `body`, `mime`, `created_at` FROM `posts` ORDER BY `created_at` DESC LIMIT 30")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	rdb.FlushAll(ctx)

	// キャッシュなし
	fmt.Println("1回目")
	for _, p := range results {
		p.User = getUser(p.UserID)
	}

	fmt.Println("キャッシュから読み込み")
	ids := make([]int, len(results))
	for _, p := range results {
		p.User = getUser(p.UserID)
		ids = append(ids, p.UserID)
	}

	fmt.Println("一気に読み込む")
	u := getUserList(ids)
	fmt.Println(u)

}

func getUser(id int) User {
	// 参考
	// https://selfnote.work/20210719/programming/golang/golang-redis/
	var user User

	var ctx = context.Background()

	user, err := getUserFromCache(ctx, id)

	if err == nil {
		fmt.Printf("[Cache hit] %v\n", user)
		return user

	}

	err = db.Get(&user, "SELECT * FROM `users` WHERE `id` = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[DB hit] %v\n", user)
	setUserToCache(ctx, user)

	return user

}

func getUserList(ids []int) []User {
	var ctx = context.Background()

	users := make([]User, 0, len(ids))

	pipe := rdb.Pipeline()
	m := map[string]*redis.StringCmd{}
	for _, id := range ids {

		m[strconv.Itoa(id)] = pipe.Get(ctx, strconv.Itoa(id))

	}

	pipe.Exec(ctx)

	for _, c := range m {
		v, err := c.Result()

		u := User{}
		if err == nil {
			err = json.Unmarshal([]byte(v), &u)
			if err == nil {
				users = append(users, u)
			}
		}
	}

	return users

}

func getUserFromCache(ctx context.Context, id int) (User, error) {

	var u User

	_, err := rdb.Pipelined(ctx, func(p redis.Pipeliner) error {

		v, err := rdb.Get(ctx, strconv.Itoa(id)).Result()

		if err != nil {
			return err
		}

		err = json.Unmarshal([]byte(v), &u)
		if err != nil {
			return err
		}

		return nil

	})

	return u, err
}

func setUserToCache(ctx context.Context, u User) {

	j, err := json.Marshal(u)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := rdb.Pipelined(ctx, func(p redis.Pipeliner) error {

		_, err = rdb.Set(ctx, strconv.Itoa(u.ID), j, redis.KeepTTL).Result()
		if err != nil {
			log.Fatal(err)
		}

		return nil

	}); err != nil {
		log.Fatal(err)
	}
}
