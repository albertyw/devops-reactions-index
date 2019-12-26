package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/albertyw/reaction-pics/tumblr"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var b = tumblr.NewBoard([]tumblr.Post{})
var d = handlerDeps{
	logger: zap.NewNop().Sugar(),
	board:  &b,
}

func TestIndexFile(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	indexHandler(response, request, d)
	assert.Equal(t, response.Code, 200)

	cacheString := appCacheString()
	assert.Contains(t, response.Body.String(), cacheString)
}

func TestOnlyIndexFile(t *testing.T) {
	request, err := http.NewRequest("GET", "/asdf", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	indexHandler(response, request, d)
	assert.Equal(t, response.Code, 404)
}

func TestReadFile(t *testing.T) {
	request, err := http.NewRequest("GET", "/static/favicon/manifest.json", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	staticHandler(response, request, d)
	assert.Equal(t, response.Code, 200)
	assert.True(t, len(response.Body.String()) > 100)
}

func TestNoExactURL(t *testing.T) {
	request, err := http.NewRequest("GET", "/static/asdf.js", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	staticHandler(response, request, d)
	assert.Equal(t, response.Code, 404)

	response = httptest.NewRecorder()
	indexHandler(response, request, d)
	assert.Equal(t, response.Code, 404)
}

func TestSearchHandler(t *testing.T) {
	request, err := http.NewRequest("GET", "/search", nil)
	assert.NoError(t, err)

	q := request.URL.Query()
	q.Add("query", "searchTerm")
	response := httptest.NewRecorder()
	searchHandler(response, request, d)
	assert.Equal(t, response.Code, 200)
	assert.Equal(t, response.Body.String(), "{\"data\":[],\"offset\":0,\"totalResults\":0}")
}

func TestSearchHandlerOffset(t *testing.T) {
	request, err := http.NewRequest("GET", "/search?offset=1", nil)
	assert.NoError(t, err)

	q := request.URL.Query()
	q.Add("query", "searchTerm")
	response := httptest.NewRecorder()
	searchHandler(response, request, d)
	assert.Equal(t, response.Code, 200)
	assert.Equal(t, response.Body.String(), "{\"data\":[],\"offset\":1,\"totalResults\":0}")
}

func TestSearchHandlerMalformedOffset(t *testing.T) {
	request, err := http.NewRequest("GET", "/search?offset=asdf", nil)
	assert.NoError(t, err)

	q := request.URL.Query()
	q.Add("query", "searchTerm")
	response := httptest.NewRecorder()
	searchHandler(response, request, d)
	assert.Equal(t, response.Code, 200)
	assert.Equal(t, response.Body.String(), "{\"data\":[],\"offset\":0,\"totalResults\":0}")
}

func TestPostHandlerMalformed(t *testing.T) {
	request, err := http.NewRequest("GET", "/post/asdf", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	postHandler(response, request, d)
	assert.Equal(t, response.Code, 404)
}

func TestPostHandlerNotFound(t *testing.T) {
	request, err := http.NewRequest("GET", "/post/1234", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	postHandler(response, request, d)
	assert.Equal(t, response.Code, 404)
}

func TestPostHandler(t *testing.T) {
	post := tumblr.Post{ID: 1234}
	d.board.AddPost(post)
	defer func() { d.board.Reset() }()
	request, err := http.NewRequest("GET", "/post/1234", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	postHandler(response, request, d)
	assert.Equal(t, response.Code, 200)
	assert.NotEqual(t, len(response.Body.String()), 0)
}

func TestPostDataHandler(t *testing.T) {
	post := tumblr.Post{ID: 1234}
	d.board.AddPost(post)
	defer func() { d.board.Reset() }()
	request, err := http.NewRequest("GET", "/postdata/1234", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	postDataHandler(response, request, d)
	assert.Equal(t, response.Code, 200)
	assert.NotEqual(t, len(response.Body.String()), 0)
}

func TestPostDataPercentHandler(t *testing.T) {
	post := tumblr.Post{ID: 1234, Title: `asdf% qwer`}
	d.board.AddPost(post)
	defer func() { d.board.Reset() }()
	request, err := http.NewRequest("GET", "/postdata/1234", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	postDataHandler(response, request, d)
	var data map[string][]map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &data)
	title := data["data"][0]["title"].(string)
	assert.Equal(t, `asdf% qwer`, title)
}

func TestPostDataHandlerMalformed(t *testing.T) {
	request, err := http.NewRequest("GET", "/postdata/asdf", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	postDataHandler(response, request, d)
	assert.Equal(t, response.Code, 404)
}

func TestPostDataHandlerUnknown(t *testing.T) {
	request, err := http.NewRequest("GET", "/postdata/1234", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	postDataHandler(response, request, d)
	assert.Equal(t, response.Code, 404)
}

func TestStatsHandler(t *testing.T) {
	request, err := http.NewRequest("GET", "/stats.json", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	statsHandler(response, request, d)
	assert.Equal(t, response.Code, 200)
	assert.Equal(t, response.Body.String(), "{\"keywords\":[],\"postCount\":\"0\"}")
}
