package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type App struct {
    urls map[string]string
}

var app  = App {
    urls: map[string]string{},
}
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func badReqRes(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusBadRequest)
    http.Error(w, "bad request", http.StatusBadRequest)
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

func main() {
    mux := http.NewServeMux()
    rand.Seed(time.Now().UnixNano())
    mux.HandleFunc("/new", newUrl)
    mux.HandleFunc("/", getShort)
    err := http.ListenAndServe(":3000", mux)
    log.Fatal(err)
}

func getShort(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path
    code := strings.Split(path, "/")[1]
    url := app.urls[code]

    if len(url) == 0  {
        notFoundReqRes(w)
        return
    }

    http.Redirect(w, r, url, http.StatusPermanentRedirect)
    return
}

func newUrl(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost:
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
        code := randCode(7)
        app.urls[code] = rawUrl
        log.Println(app.urls)
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
