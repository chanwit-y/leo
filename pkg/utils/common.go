package utils

import (
	"bufio"
	"os"
)

func ternary[T any](condition bool, v1 T, v2 T) T {
	if condition {
		return v1
	} else {
		return v2
	}

}

func space(len int) string {
	space := ""
	for i := 0; i <= len; i++ {
		space += " "
	}
	return space
}

func createFile(name string, schema []string) {
	file, _ := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	datawrite := bufio.NewWriter(file)

	for _, data := range schema {
		_, _ = datawrite.WriteString(data)
	}

	datawrite.Flush()
	file.Close()
}
