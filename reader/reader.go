// Package reader implements functionality for reading user defined RSS feeds.
package reader

import (
	"github.com/mmcdole/gofeed"
)

// RSS contains the RSS feed URLs to read.
type RSS struct {
	URLS []string
}

// Read iterates over the list of RSS feed URls and return a list of generic feed structs.
func (r *RSS) Read() ([]*gofeed.Feed, error) {
	feeds := make([]*gofeed.Feed, 0)

	for _, url := range r.URLS {
		feed, err := gofeed.NewParser().ParseURL(url)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}

	return feeds, nil
}
