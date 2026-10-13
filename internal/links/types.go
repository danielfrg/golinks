package links

import (
	"net/url"
	"time"
)

type Link struct {
	Short    string    `json:"short"`
	Long     string    `json:"long"`
	Created  time.Time `json:"created"`
	LastEdit time.Time `json:"lastEdit"`
	Clicks   int64     `json:"clicks"`
}

type Environment struct {
	Path  string
	Now   time.Time
	Query url.Values
}
