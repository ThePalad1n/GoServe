package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Post struct {
	ID       int    `json:"id"`
	Created  string `json:"created"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Author   string `json:"author"`
	Imageurl string `json:"imageurl"`
}

func CreateRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", CreatePost).Methods("POST")
	router.HandleFunc("/post/{id}", FindPost).Methods("GET")

	return router
}

func GetPost() Post {
	var post Post
	return post
}

func GetPosts() []Post {
	var posts []Post
	return posts
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/mysocial")
	if err != nil {
		panic(err.Error())
	}
	enableCors(&w)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS mysocial.post(id int auto_increment primary key, created timestamp default current_timestamp not null, title varchar(128), content text, author varchar(64), imageurl varchar(128))")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("\nEvent fields:\n")

	fmt.Println("\nJSON representation:\n")

	var Posts Post
	w.Header().Set("Content-Type", "application/json")
	str, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(str, &Posts)
	fmt.Println(Posts)

	insertEvent(db, Posts)

}

func FindPost(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/mysocial")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("\nEvent fields:\n")

	fmt.Println("\nJSON representation:\n")

	var Posts Post
	w.Header().Set("Content-Type", "application/json")
	str, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(str, &Posts)
	fmt.Println(str)
	fmt.Println(Posts)
	selectEventById(db, int64(Posts.ID), &Posts)
}

var (
	insertEventQuery     = (`INSERT INTO mysocial.post(created,title,content,author,imageurl) VALUES (?, ?, ?, ?, ?)`)
	selectEventByIdQuery = `SELECT * FROM post WHERE id = ?`
)

func insertEvent(db *sql.DB, post Post) (int64, error) {
	fmt.Println(post)
	res, err := db.Exec(insertEventQuery, post.Created, post.Title, post.Content, post.Author, post.Imageurl)
	if err != nil {
		return 0, err
	}
	lid, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lid, nil
}

func selectEventById(db *sql.DB, id int64, post *Post) error {
	row := db.QueryRow(selectEventByIdQuery, id)
	err := row.Scan(&post.ID, &post.Created, &post.Title, &post.Content, &post.Author, &post.Imageurl)
	if err != nil {
		return err
	}
	return nil
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	fmt.Println("Listening on port 6969...")
	var a = CreateRouter()
	log.Fatal(http.ListenAndServe(":6969", a))

}
