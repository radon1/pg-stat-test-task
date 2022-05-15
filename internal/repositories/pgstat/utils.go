package pgstat

import (
	"errors"
	"fmt"
	"strings"
)

var malformedArgsErr = errors.New("malformed args")

func buildLikeArgsArray(items []string) (string, error) {
	likeFilterArgs := strings.Builder{}
	likeFilterArgs.WriteString("{")
	for i, item := range items {
		if len(item) == 0 {
			return "", malformedArgsErr
		}

		likeFilterArgs.WriteString("%")
		likeFilterArgs.WriteString(fmt.Sprintf("%s", item))
		likeFilterArgs.WriteString("%")
		if i < len(items)-1 {
			likeFilterArgs.WriteString(",")
		}
	}
	likeFilterArgs.WriteString("}")
	return likeFilterArgs.String(), nil
}
