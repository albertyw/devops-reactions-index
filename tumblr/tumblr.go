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
	blogName   = "devopsreactions.tumblr.com"
	blogTypes  = "text"
	postsLimit = 20
)

func GetPosts() []gotumblr.TextPost {
	var posts, newPosts []gotumblr.TextPost
	offset := 0
	client := getTumblrClient()
	for len(newPosts) == postsLimit || offset == 0 {
		fmt.Println("Downloading", offset)
		options := getTumblrOptions(offset)
		postsResponse := client.Posts(blogName, blogTypes, options)
		newPosts = parsePosts(postsResponse)
		posts = append(posts, newPosts...)
		offset += postsLimit
	}
	return posts
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
	return options
}

func parsePosts(postsResponse gotumblr.PostsResponse) []gotumblr.TextPost {
	var posts []gotumblr.TextPost
	var post gotumblr.TextPost
	for _, element := range postsResponse.Posts {
		err := json.Unmarshal(element, &post)
		if err != nil {
			fmt.Println(err)
		} else {
			posts = append(posts, post)
		}
	}
	return posts
}