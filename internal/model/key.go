package model

import "fmt"

type Key string

func NewKey(ip, uri string) Key {
	return Key(fmt.Sprintf("%s:%s", ip, uri))
}
