package data

import (
	"errors"
	"strings"
)

type Entry struct {
	Id       int64
	Title    string
	Username string
	Password string
	Address  string
	Notes    string
}

type State struct {
	LastId int64
	Count  int64
}

type Configuration struct {
	EntryThreshold int64
}

var config = Configuration{
	EntryThreshold: 1_000_000,
}

type Order int8

const (
	ByTitle Order = iota
	ByUsername
)

func OrderFromString(status string) (Order, error) {
	switch strings.ToLower(status) {
	case "title":
		return ByTitle, nil
	case "username":
		return ByUsername, nil
	default:
		return -1, errors.New("invalid Order")
	}
}
