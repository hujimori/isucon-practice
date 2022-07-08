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

	// キャッシュなし

	fmt.Println("1回目")
	for _, p := range results {
		p.User = getUser(p.UserID)
		// fmt.Printf("[%d] %v\n", p.ID, p.User)
	}

	fmt.Println("2回目")
	// キャッシュあり

	ids := make([]int, len(results))
	for _, p := range results {
		p.User = getUser(p.UserID)
		ids = append(ids, p.UserID)
		// fmt.Printf("[%d] %v\n", p.UserID, p.User)
	}

	// fmt.Printf("%v\n", ids)

	fmt.Println("3回目")
	var u []User
	defer func() {
		u = getUserList(ids)
	}()
	fmt.Println(u)

	// for _, u := range users {
	// fmt.Printf("[%d] %v\n", u.ID, u)
	// }
}

func getUser(id int) User {
	// 参考
	// https://selfnote.work/20210719/programming/golang/golang-redis/
	var user User

	var ctx = context.Background()
	key := strconv.Itoa(id)
	it, err := rdb.Get(ctx, key).Result()

	if err == nil {
		err = json.Unmarshal([]byte(it), &user)
		if err == nil {
			fmt.Printf("hit! %v\n", id)

			return user
		}
	}

	err = db.Get(&user, "SELECT * FROM `users` WHERE `id` = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	j, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	_, err = rdb.Set(ctx, key, j, redis.KeepTTL).Result()
	if err != nil {
		log.Fatal(err)
	}

	return user

}

func getUserList(ids []int) []User {
	var ctx = context.Background()

	users := make([]User, len(ids))
	// its := make([]bytes, len(ids))
	pipe := rdb.Pipeline()
	m := map[string]*redis.StringCmd{}
	for _, id := range ids {

		m[strconv.Itoa(id)] = pipe.Get(ctx, strconv.Itoa(id))
		// u := User{}
		// if err == nil {
		// 	err = json.Unmarshal([]byte(it), &u)
		// 	if err == nil {
		// 		users = append(users, u)
		// 	}
		// }

	}

	_, err := pipe.Exec(ctx)
	fmt.Println(err)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println(len(m))

	for _, v := range m {
		res, _ := v.Result()
		fmt.Println(res)

		u := User{}
		if err == nil {
			err = json.Unmarshal([]byte(res), &u)
			if err == nil {
				users = append(users, u)
			}
		}

	}

	// pipe.Exec(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return users

}
