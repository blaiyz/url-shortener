package store

// Datastore for shortened URLs and their id (the url parameter for the shortener api).
type Datastore interface {
	SetNext(url string) (id string)
	Get(id string) (url string, ok bool)
}
