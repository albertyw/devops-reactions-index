package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/albertyw/reaction-pics/tumblr"
)

func TestReadFile(t *testing.T) {
	request, err := http.NewRequest("GET", "/app.js", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	indexHandler(response, request)
	assert.Equal(t, response.Code, 200)

	cacheString := appCacheString()
	assert.Contains(t, response.Body.String(), cacheString)
}

func TestIndexFile(t *testing.T) {
	handler := rewriteFS(http.FileServer(http.Dir(staticPath)).ServeHTTP)
	request, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	handler(response, request)
	assert.Equal(t, response.Code, 200)
	assert.True(t, len(response.Body.String()) > 100)
}

func TestExactURL(t *testing.T) {
	handler := rewriteFS(http.FileServer(http.Dir(staticPath)).ServeHTTP)
	request, err := http.NewRequest("GET", "/asdf", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	handler(response, request)
	assert.Equal(t, response.Code, 404)
}

func TestNoExactURL(t *testing.T) {
	handler := rewriteFS(http.FileServer(http.Dir(staticPath)).ServeHTTP)
	request, err := http.NewRequest("GET", "/asdf", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	handler(response, request)
	assert.Equal(t, response.Code, 404)
}
func TestSearchHandler(t *testing.T) {
	request, err := http.NewRequest("GET", "/search", nil)
	assert.NoError(t, err)

	q := request.URL.Query()
	q.Add("query", "searchTerm")
	response := httptest.NewRecorder()
	searchHandler(response, request)
	assert.Equal(t, response.Code, 200)
	assert.Equal(t, response.Body.String(), "{\"data\":[],\"offset\":0,\"totalResults\":0}")
}

func TestSearchHandlerOffset(t *testing.T) {
	request, err := http.NewRequest("GET", "/search?offset=1", nil)
	assert.NoError(t, err)

	q := request.URL.Query()
	q.Add("query", "searchTerm")
	response := httptest.NewRecorder()
	searchHandler(response, request)
	assert.Equal(t, response.Code, 200)
	assert.Equal(t, response.Body.String(), "{\"data\":[],\"offset\":1,\"totalResults\":0}")
}

func TestSearchHandlerMalformedOffset(t *testing.T) {
	request, err := http.NewRequest("GET", "/search?offset=asdf", nil)
	assert.NoError(t, err)

	q := request.URL.Query()
	q.Add("query", "searchTerm")
	response := httptest.NewRecorder()
	searchHandler(response, request)
	assert.Equal(t, response.Code, 200)
	assert.Equal(t, response.Body.String(), "{\"data\":[],\"offset\":0,\"totalResults\":0}")
}

func TestPostHandlerMalformed(t *testing.T) {
	request, err := http.NewRequest("GET", "/post/asdf", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	postHandler(response, request)
	assert.Equal(t, response.Code, 404)
}

func TestPostHandlerNotFound(t *testing.T) {
	request, err := http.NewRequest("GET", "/post/1234", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	postHandler(response, request)
	assert.Equal(t, response.Code, 404)
}

func TestPostHandler(t *testing.T) {
	post := tumblr.Post{ID: 1234}
	board.AddPost(post)
	defer func() { board.Reset() }()
	request, err := http.NewRequest("GET", "/post/1234", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	postHandler(response, request)
	assert.Equal(t, response.Code, 200)
	assert.NotEqual(t, len(response.Body.String()), 0)
}

func TestPostDataHandler(t *testing.T) {
	post := tumblr.Post{ID: 1234}
	board.AddPost(post)
	defer func() { board.Reset() }()
	request, err := http.NewRequest("GET", "/postdata/1234", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	postDataHandler(response, request)
	assert.Equal(t, response.Code, 200)
	assert.NotEqual(t, len(response.Body.String()), 0)
}

func TestPostDataHandlerMalformed(t *testing.T) {
	request, err := http.NewRequest("GET", "/postdata/asdf", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	postDataHandler(response, request)
	assert.Equal(t, response.Code, 404)
}

func TestPostDataHandlerUnknown(t *testing.T) {
	request, err := http.NewRequest("GET", "/postdata/1234", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	postDataHandler(response, request)
	assert.Equal(t, response.Code, 404)
}

func TestStatsHandler(t *testing.T) {
	request, err := http.NewRequest("GET", "/stats.json", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	statsHandler(response, request)
	assert.Equal(t, response.Code, 200)
	assert.Equal(t, response.Body.String(), "{\"keywords\":[],\"postCount\":\"0\"}")
}

func TestAddPost(t *testing.T) {
	oldLength := len(board.Posts)
	post1 := tumblr.Post{ID: 1, Title: "1"}
	post2 := tumblr.Post{ID: 2, Title: "2"}
	post3 := tumblr.Post{ID: 3, Title: "3"}
	postChan := make(chan tumblr.Post, 1)
	defer func() { close(postChan) }()
	go loadPosts(postChan)
	postChan <- post1
	postChan <- post2
	postChan <- post3
	assert.NotEqual(t, len(board.Posts), oldLength)
	assert.True(t, board.Posts[0].ID < board.Posts[1].ID)
}

func TestCacheString(t *testing.T) {
	tempStaticPath := staticPath
	defer func() { staticPath = tempStaticPath }()
	cacheString := appCacheString()
	assert.NotEqual(t, cacheString, "")

	staticPath = ""
	cacheString = appCacheString()
	assert.Equal(t, cacheString, "")
}
