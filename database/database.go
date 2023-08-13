package database

import (
	"fmt"
	"strings"
)

func PIDArrayToString(array []uint32) string {
	return strings.Trim(strings.Replace(fmt.Sprint(array), " ", ",", -1), "[]")
}
