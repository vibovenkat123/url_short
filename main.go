package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type App struct {
	db *sql.DB
}

var app = App{}
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func badReqRes(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	http.Error(w, "bad request", http.StatusBadRequest)
}

func serverErrReqRes(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func notFoundReqRes(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	http.Error(w, "not found", http.StatusNotFound)
}

func createdReqRes(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
}

func randCode(length int) string {
	bytes := make([]rune, length)
	for i := range bytes {
		bytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(bytes)
}

var schema = `
CREATE TABLE IF NOT EXISTS url (raw_url TEXT, code VARCHAR(7));
`

func main() {
	dbFile := os.Getenv("url_db_file")
	if len(dbFile) == 0 {
		dbFile = "main.db"
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	app.db = db
	execShema(schema, app.db)
	defer db.Close()

	if err != nil {
		log.Panicf("Cannot connect to database `%s`", dbFile)
	}

	mux := http.NewServeMux()
	rand.Seed(time.Now().UnixNano())
	mux.HandleFunc("/new", newUrl)
	mux.HandleFunc("/", getShort)
	err = http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}

func execShema(schema string, db *sql.DB) {
	_, err := db.Exec(schema)
	if err != nil {
		log.Printf("%q: %s\n", err, schema)
	}
}

func getShort(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	params := strings.Split(path, "/")
	code := params[1]
	if len(code) == 0 || len(params) != 2 {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Go to `/new?url=urlhere` to create a new shortened url")
		return
	}
	var url string
	stmt, err := app.db.Prepare("select raw_url from url where code = ?")

	if err != nil {
		serverErrReqRes(err, w)
		return
	}

	err = stmt.QueryRow(code).Scan(&url)

	if err != nil {
		serverErrReqRes(err, w)
		return
	}

	if len(url) == 0 {
		notFoundReqRes(w)
		return
	}

	http.Redirect(w, r, url, http.StatusPermanentRedirect)
	return
}

func newUrl(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		code := randCode(7)
		for ok := true; ok; ok = len(code) == 0 {
			code = randCode(7)
		}
		rawUrl := r.FormValue("url")
		if len(rawUrl) == 0 {
			badReqRes(w)
			return
		}
		url, err := url.ParseRequestURI(rawUrl)
		if err != nil || url == nil {
			badReqRes(w)
			return
		}
		stmt, err := app.db.Prepare("INSERT INTO url(raw_url, code) values (?, ?)")
		if err != nil {
			serverErrReqRes(err, w)
			return
		}
		defer stmt.Close()
		_, err = stmt.Exec(rawUrl, code)
		if err != nil {
			serverErrReqRes(err, w)
			return
		}
		createdReqRes(w)
		fmt.Fprintf(w, code)
		return

	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)

	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
