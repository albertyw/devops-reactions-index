package server

import (
	"fmt"
	"github.com/albertyw/reaction-pics/tumblr"
	// Used for getting tumblr env vars
	_ "github.com/joho/godotenv/autoload"
	"net/http"
	"os"
	"strings"
)

const dataURLPath = "/data.json"

var templateDir = os.Getenv("SERVER_TEMPLATES")
var indexPath = fmt.Sprintf("%s/index.htm", templateDir)
var jsPath = fmt.Sprintf("%s/app.js", templateDir)
var cssPath = fmt.Sprintf("%s/global.css", templateDir)
var uRLFilePaths = map[string]func() (string, error){}
var posts []tumblr.Post

// logURL is a closure that logs (to stdout) the url and query of requests
func logURL(
	targetFunc func(http.ResponseWriter, *http.Request),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()
		fmt.Println(url)
		targetFunc(w, r)
	}
}

// exactURL is a closure that checks that the http match is an exact url path
// match instead of allowing for net/http's loose match
func exactURL(
	targetFunc func(http.ResponseWriter, *http.Request),
	requestedPath string,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != requestedPath {
			http.NotFound(w, r)
			return
		}
		targetFunc(w, r)
		return
	}
}

// readFile returns a function that reads the file at a given path and makes a
// response from it
func readFile(p string) func(http.ResponseWriter, *http.Request) {
	path := p
	return func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		fileInfo, err := os.Stat(path)
		if err != nil {
			return
		}
		http.ServeContent(w, r, path, fileInfo.ModTime(), file)
	}
}

// dataURLHandler is an http handler for the dataURLPath response
func dataURLHandler(w http.ResponseWriter, r *http.Request) {
	html := tumblr.PostsToJSON(posts)
	fmt.Fprintf(w, html)
}

// searchHandler is an http handler to search data for keywords
func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	query = strings.ToLower(query)
	selectedPosts := []tumblr.Post{}
	for _, post := range posts {
		postData := strings.ToLower(post.Title)
		if strings.Contains(postData, query) {
			selectedPosts = append(selectedPosts, post)
		}
	}
	html := tumblr.PostsToJSON(selectedPosts)
	fmt.Fprintf(w, html)
}

// Run starts up the HTTP server
func Run(p []tumblr.Post) {
	posts = p
	address := ":" + os.Getenv("PORT")
	fmt.Println("server listening on", address)
	http.HandleFunc("/", logURL(exactURL(readFile(indexPath), "/")))
	http.HandleFunc("/app.js", logURL(readFile(jsPath)))
	http.HandleFunc("/global.css", logURL(readFile(cssPath)))
	http.HandleFunc(dataURLPath, logURL(dataURLHandler))
	http.HandleFunc("/search", logURL(searchHandler))
	http.ListenAndServe(address, nil)
}
