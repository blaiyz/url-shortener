package store

import (
	"github.com/jxskiss/base62"
)

type memoryStore struct {
	current uint64
	urls    map[string]string
	special string
}

func NewMemoryStore(special string) *memoryStore {
	return &memoryStore{
		current: 100,
		urls:    make(map[string]string),
		special: special,
	}
}

func (store *memoryStore) Get(id string) (string, bool) {
	url, ok := store.urls[id]
	return url, ok
}

func (store *memoryStore) SetNext(url string) string {
	newId := string(base62.FormatUint(store.current))
	store.current++

	// Skip the special path for collecting the shortened url
	if newId == store.special {
		newId = string(base62.FormatUint(store.current))
		store.current++
	}
	store.urls[newId] = url
	return newId
}
