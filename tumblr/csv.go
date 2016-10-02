package tumblr

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func ReadPostsFromCSV() (posts []Post) {
	file, err := os.Open("data.csv")
	defer file.Close()
	if err != nil {
		return
	}
	reader := csv.NewReader(file)
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}
		post := CSVToPost(row)
		posts = append(posts, *post)
	}
	return posts
}

func WritePostsToCSV(posts []Post) {
	var row []string
	file, err := os.Create("data.csv")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	writer := csv.NewWriter(file)
	for _, post := range posts {
		row = getRow(post)
		writer.Write(row)
		if err := writer.Error(); err != nil {
			fmt.Println(err)
		}
	}
	writer.Flush()
}

func getRow(post Post) (row []string) {
	row = []string{
		strconv.FormatInt(post.Id, 10),
		post.Title,
		post.Url,
	}
	return
}
