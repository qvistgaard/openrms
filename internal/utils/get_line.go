package utils

import (
	"io/fs"
	"strings"
)

func GetLine(fs fs.ReadFileFS, file string, line int) (string, error) {
	data, err := fs.ReadFile(file)
	if err != nil {
		return "", err
	}
	split := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")

	intn := len(split)/line - 1

	return split[intn], nil
}
