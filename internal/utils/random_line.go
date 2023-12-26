package utils

import (
	"io/fs"
	"math/rand"
	"strings"
)

func RandomLine(fs fs.ReadFileFS, file string) (string, error) {
	data, err := fs.ReadFile(file)
	if err != nil {
		return "", err
	}
	split := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	intn := rand.Intn(len(split) - 1)

	return split[intn], nil
}
