package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Post struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

var db *sql.DB
var err error

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	/*fmt.Printf("godotenv : %s = %s \n", "Site Title", siteTitle)
	fmt.Printf("godotenv : %s = %s \n", "DB Host", dbHost)
	fmt.Printf("godotenv : %s = %s \n", "DB Port", dbPort)
	fmt.Printf("godotenv : %s = %s \n", "username", dbUsername)
	fmt.Printf("godotenv : %s = %s \n", "password", dbPassword)
	fmt.Printf("godotenv : %s = %s \n", "db Name", dbname)*/
	//db, err = sql.Open("mysql", "root:akash@tcp(127.0.0.1:3306)/data")
	dsn := dbUsername + ":" + dbPassword + "@" + dbHost + "/" + dbname + "?charset=utf8"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// getting env variables SITE_TITLE and DB_HOST

	router := mux.NewRouter()
	router.HandleFunc("/", add).Methods("GET")
	router.HandleFunc("/table", createPostTable).Methods("GET")
	router.HandleFunc("/user", getPosts).Methods("GET")
	router.HandleFunc("/useradd", createPost).Methods("POST")
	router.HandleFunc("/user/{id}", getPost).Methods("GET")
	//router.HandleFunc("/health", HealthCheckHandler)
	router.HandleFunc("/health", healthApi)
	http.ListenAndServe(":8000", router)
}

func add(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome")
}

func createPostTable(w http.ResponseWriter, r *http.Request) {
	query := `CREATE TABLE IF NOT EXISTS posts(id INT AUTO_INCREMENT primary key NOT NULL , title varchar(255))`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating product table", err)
		return
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return
	}
	log.Printf("Rows affected when creating table: %d", rows)
	return
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var posts []Post
	result, err := db.Query("SELECT id, title from posts")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var post Post
		err := result.Scan(&post.ID, &post.Title)
		if err != nil {
			panic(err.Error())
		}
		posts = append(posts, post)
	}
	json.NewEncoder(w).Encode(posts)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stmt, err := db.Prepare("INSERT INTO posts(title) VALUES(?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	title := keyVal["title"]
	_, err = stmt.Exec(title)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New post was created")
}

func getPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, title FROM posts WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var post Post
	for result.Next() {
		err := result.Scan(&post.ID, &post.Title)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(post)
}

/*func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"responce": 200}`)

}*/

func healthApi(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:8000/user")
	if err != nil {
		log.Fatalf("HTTP GET request failed, %v\n", err)
	}
	fmt.Fprintf(w, "<h1>Health check is done  %v</h1>", resp.Status)
	defer resp.Body.Close()
	fmt.Println("Response status:", resp.Status)
	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Body read failed: %v\n", err)
	}
}
