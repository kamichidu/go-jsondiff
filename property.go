package jsondiff

import (
	"strconv"
	"strings"
)

func createPathTester(pathItems []string) interface{} {
	if len(pathItems) == 0 {
		return struct{}{}
	}
	head, tail := pathItems[0], pathItems[1:]
	val := createPathTester(tail)
	if head == "$" {
		return val
	}
	if strings.HasSuffix(head, "]") {
		head = strings.TrimSuffix(head, "]")
		idx := strings.LastIndex(head, "[")
		if idx < 0 {
			panic("invalid path string: " + strings.Join(pathItems, "."))
		}
		i, err := strconv.Atoi(head[idx+1:])
		if err != nil {
			panic("invalid path string: " + strings.Join(pathItems, "."))
		}
		ary := make([]interface{}, i+1)
		ary[i] = val
		val = ary
		head = head[:idx]
	}
	return map[string]interface{}{
		head: val,
	}
}
