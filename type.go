package bakusai

import "time"

type Article struct {
	ID         int
	ThreadID   int
	PostAt     time.Time
	BodyText   string
	AuthorName string
}

type Thread struct {
	URI       string
	CreatedAt time.Time
	Title     string
	PrevURI   string
	NextURI   string
	Articles  []*Article
}
