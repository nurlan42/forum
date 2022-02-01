package internal

import (
	"net/http"
	"strconv"
	"strings"
)

func GetPostID(r *http.Request) (int, error) {
	str := strings.TrimPrefix(r.URL.Path, "/post/")
	postID, err := strconv.Atoi(str)

	if err != nil {
		return 0, err
	}
	return postID, nil
}
