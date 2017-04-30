package tumblr

import (
	"encoding/json"
	"fmt"
	"github.com/MariaTerzieva/gotumblr"
	// Used for getting tumblr env vars
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strconv"
)

const (
	tumblrURL  = "http://api.tumblr.com"
	blogTypes  = "text"
	postsLimit = 20
)

// Blogs is a list of all blogs to read from
var Blogs = []string{
	"devopsreactions.tumblr.com",
}

func duplicateChan(in <-chan Post, out1, out2 chan<- Post) {
	for p := range in {
		out1 <- p
		out2 <- p
	}
	close(out1)
	close(out2)
}

// GetPosts returns a list of all Posts
func GetPosts(getNewPosts bool, out1 chan<- Post) {
	for _, blogName := range Blogs {
		posts := make(chan Post)
		out2 := make(chan Post)
		go duplicateChan(posts, out1, out2)
		go getBlogPosts(blogName, getNewPosts, posts)
		go WritePostsToCSV(blogName, out2)
	}
}

func getBlogPosts(blogName string, getNewPosts bool, posts chan<- Post) {
	defer func() { close(posts) }()
	var newPosts []Post
	existingPosts := ReadPostsFromCSV(blogName)
	maxPostID := int64(0)
	for _, p := range existingPosts {
		if p.ID > maxPostID {
			maxPostID = p.ID
		}
		posts <- p
	}
	if !getNewPosts {
		return
	}
	offset := 0
	client := getTumblrClient()
	for len(newPosts) == postsLimit || offset == 0 {
		fmt.Println("Downloading", offset)
		options := getTumblrOptions(offset)
		postsResponse := client.Posts(blogName, blogTypes, options)
		newPosts = parsePosts(postsResponse)
		for _, p := range newPosts {
			if p.ID <= maxPostID {
				return
			}
			posts <- p
		}
		offset += postsLimit
	}
}

func getTumblrClient() *gotumblr.TumblrRestClient {
	client := gotumblr.NewTumblrRestClient(
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"),
		os.Getenv("TOKEN"),
		os.Getenv("TOKEN_SECRET"),
		"https://www.albertyw.com/",
		tumblrURL,
	)
	return client
}

func getTumblrOptions(offset int) map[string]string {
	options := map[string]string{}
	options["offset"] = strconv.Itoa(offset)
	options["limit"] = strconv.Itoa(postsLimit)
	options["notes_info"] = strconv.FormatBool(true)
	return options
}

func parsePosts(postsResponse gotumblr.PostsResponse) []Post {
	var posts []Post
	var tumblrPost gotumblr.TextPost
	for _, element := range postsResponse.Posts {
		err := json.Unmarshal(element, &tumblrPost)
		if err != nil {
			fmt.Println(err)
		} else {
			post := GoTumblrToPost(&tumblrPost)
			posts = append(posts, *post)
		}
	}
	return posts
}
