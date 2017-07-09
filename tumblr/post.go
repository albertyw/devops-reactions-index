package tumblr

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"github.com/MariaTerzieva/gotumblr"
	"github.com/gosimple/slug.git"
	"golang.org/x/net/html"
)

// Post is a representation of a single tumblr post
type Post struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Image string `json:"image"`
	Likes int64  `json:"likes"`
}

// PostJSON is a representation of Post for creating JSON values
type PostJSON struct {
	Post
	InternalURL string `json:"internalURL"`
}

// InternalURL returns the path to the post
func (p Post) InternalURL() string {
	slug := slug.Make(p.Title)
	if len(slug) > 30 {
		slug = slug[0:30]
	}
	return "/post/" + strconv.FormatInt(p.ID, 10) + "/" + slug
}

// ToJSONStruct builds a PostJSON based on the Post
func (p Post) ToJSONStruct() PostJSON {
	return PostJSON{
		Post:        p,
		InternalURL: p.InternalURL(),
	}
}

// GoTumblrToPost converts a gotumblr.TextPost into a Post
func GoTumblrToPost(tumblrPost *gotumblr.TextPost) *Post {
	image := getImageFromPostBody(tumblrPost.Body)
	likes := tumblrPost.Note_count
	post := Post{
		ID:    tumblrPost.Id,
		Title: strings.TrimSpace(tumblrPost.Title),
		URL:   tumblrPost.Post_url,
		Image: image,
		Likes: likes,
	}
	return &post
}

// CSVToPost converts a CSV row into a Post
func CSVToPost(row []string) *Post {
	id, err := strconv.ParseInt(row[0], 10, 64)
	if err != nil {
		id = 0
	}
	likes, err := strconv.ParseInt(row[4], 10, 64)
	if err != nil {
		likes = 0
	}
	post := Post{
		ID:    id,
		Title: row[1],
		URL:   row[2],
		Image: row[3],
		Likes: likes,
	}
	return &post
}

// PostsToJSON converts a Post into a JSON string
func PostsToJSON(posts []Post) string {
	postsJSON := make([]PostJSON, len(posts))
	for i := 0; i < len(posts); i++ {
		postsJSON[i] = posts[i].ToJSONStruct()
	}
	marshalledPosts, _ := json.Marshal(postsJSON)
	return string(marshalledPosts)
}

// SortPostsByID sorts Posts in reverse ID order
func SortPostsByID(posts *[]Post) *[]Post {
	sort.Sort(sort.Reverse(SortByID(*posts)))
	return posts
}

// SortByID is an interface for Sorting
type SortByID []Post

func (a SortByID) Len() int           { return len(a) }
func (a SortByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByID) Less(i, j int) bool { return a[i].ID < a[j].ID }

// SortPostsByLikes sorts Posts in reverse number of likes order
func SortPostsByLikes(posts *[]Post) *[]Post {
	sort.Sort(sort.Reverse(SortByLikes(*posts)))
	return posts
}

// SortByLikes is an interface for Sorting
type SortByLikes []Post

func (a SortByLikes) Len() int           { return len(a) }
func (a SortByLikes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByLikes) Less(i, j int) bool { return a[i].Likes < a[j].Likes }

// getImageFromPostBody parses the body of a post and returns the url of the image
func getImageFromPostBody(body string) string {
	bodyReader := strings.NewReader(body)
	tokenizer := html.NewTokenizer(bodyReader)

	for tokenizer.Next() != html.ErrorToken {
		token := tokenizer.Token()
		if token.Data == "img" {
			for _, attr := range token.Attr {
				if attr.Key == "src" {
					return attr.Val
				}
			}
		}
	}
	return ""
}
