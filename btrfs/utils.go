package btrfs

import (
	"io/ioutil"
	"strconv"
	"strings"
)

func readString(path string) (string, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(raw)), nil
}

func readInt(path string) (int, error) {
	str, err := readString(path)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(str)
}

func readDir(path string) ([]string, error) {
	xs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, x := range xs {
		result = append(result, x.Name())
	}

	return result, nil
}
