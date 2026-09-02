package model

import "fmt"

type key string

func NewKey(ip, uri string) key {
	return key(fmt.Sprintf("%s:%s", ip, uri))
}
