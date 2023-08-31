package database

import (
	"strconv"
	"strings"
)

func intArrayToString(arr []int) string {
	values := make([]string, len(arr))
	for i, v := range arr {
		values[i] = strconv.Itoa(v)
	}
	return "{" + strings.Join(values, ",") + "}"
}
