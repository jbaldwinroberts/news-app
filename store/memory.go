// Package store implements functionality for creating, querying, manipulating and deleting stored RSS feeds
package store

import (
	"log"
	"sort"
	"sync"
	"time"

	"github.com/josephroberts/esqimo/reader"
	Swagger "github.com/josephroberts/esqimo/swagger"
	"github.com/mmcdole/gofeed"
)

// Store defines the interface that must be implemented by storage types.
// Should make it easier to replace MemoryStore in future
type Store interface {
	GetItems(titles []string, categories []string, limit int32) []Swagger.Item
	GetTitles() []string
	GetCategories() []string
}

// MemoryStore contains a RWMutex to prevent concurrent read/writes, a reader interface for reading RSS feeds, and data storage
type MemoryStore struct {
	sync.RWMutex
	reader     reader.RSS
	feeds      Swagger.Feeds
	titles     map[string]bool // Used a map so that the keys will be unique
	categories map[string]bool // Used a map so that the keys will be unique
}

// New initialises MemoryStore and starts a go routine to refresh the stored data every updateInterval.
func New(reader reader.RSS, updateInterval time.Duration) *MemoryStore {
	m := &MemoryStore{
		sync.RWMutex{},
		reader,
		Swagger.Feeds{},
		make(map[string]bool, 0),
		make(map[string]bool, 0),
	}

	// Populate the data in MemoryStore
	if err := m.refresh(); err != nil {
		log.Fatalf("unable to refesh memory store, %v", err)
	}

	// Refresh the data in MemoryStore every updateInterval
	ticker := time.NewTicker(updateInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := m.refresh(); err != nil {
					log.Fatalf("unable to refesh memory store, %v", err)
				}
			}
		}
	}()

	return m
}

// GetItems returns a list of items.
func (m *MemoryStore) GetItems(titles []string, categories []string, limit int32) []Swagger.Item {
	items := make([]Swagger.Item, 0)

	// Initialise the titles filter
	titlesFilter := m.titles

	// Overwrite filter if titles parameter is in the request
	if len(titles) != 0 {
		titlesFilter = convertSliceToMap(titles)
	}

	// Initialise the categories filter
	categoriesFilter := m.categories

	// Add an "" empty string to filter
	// This prevents items that don't have a category from being filtered out
	categoriesFilter[""] = true

	// Overwrite filter if categories parameter is in the request
	if len(categories) != 0 {
		categoriesFilter = convertSliceToMap(categories)
	}

	for key, value := range m.feeds.AdditionalProperties {
		// Filter by title
		if _, exists := titlesFilter[key]; !exists {
			continue
		}

		for _, value := range value.Items.AdditionalProperties {
			// Filter by category
			if _, exists := categoriesFilter[value.Category]; !exists {
				continue
			}

			items = append(items, value)
		}
	}

	// Sort by PublishParsed, most recent items first and oldest items last
	sort.Slice(items, func(i, j int) bool {
		return items[i].PublishedParsed.After(*items[j].PublishedParsed)
	})

	// Limit number of items returned
	if limit != 0 {
		return items[:limit]
	}

	return items
}

// GetTitles returns a list of titles.
func (m *MemoryStore) GetTitles() []string {
	// Aquire a read lock, defer release until the function returns
	m.RLock()
	defer m.RUnlock()

	return convertMapToSlice(m.titles)
}

// GetCategories returns a list of categories.
func (m *MemoryStore) GetCategories() []string {
	// Aquire a read lock, defer release until the function returns
	m.RLock()
	defer m.RUnlock()

	return convertMapToSlice(m.categories)
}

// Refresh stored data in MemoryStore
func (m *MemoryStore) refresh() error {
	// Read feeds
	feeds, err := m.reader.Read()
	if err != nil {
		return err
	}

	// Aquire a read/write lock, defer release until the function returns
	m.Lock()
	defer m.Unlock()

	// Refresh stored data
	m.feeds, m.titles, m.categories = convertFeeds(feeds)

	return nil
}

// Convert feeds from generic feed structs into a swagger feeds struct
func convertFeeds(feeds []*gofeed.Feed) (Swagger.Feeds, map[string]bool, map[string]bool) {
	swagFeeds := Swagger.Feeds{
		AdditionalProperties: make(map[string]Swagger.Feed),
	}
	titles := make(map[string]bool, 0)
	categories := make(map[string]bool, 0)

	for _, feed := range feeds {

		swagItems := Swagger.Items{
			AdditionalProperties: make(map[string]Swagger.Item),
		}

		for _, item := range feed.Items {
			swagItem := Swagger.Item{
				Title:           item.Title,
				Link:            item.Link,
				Description:     item.Description,
				Category:        extractCategory(item.Categories),
				Guid:            item.GUID,
				Published:       item.Published,
				PublishedParsed: item.PublishedParsed,
				// Copy the feed image to the item so that an image can be displayed along with each item in the frontend
				Image: Swagger.Image{
					Title: &feed.Image.Title,
					Url:   &feed.Image.URL,
				},
			}

			swagItems.AdditionalProperties[swagItem.Guid] = swagItem

			// Store categories so they can be used for the GetCategories endpoint
			if swagItem.Category != "" {
				categories[swagItem.Category] = true
			}
		}

		swagFeed := Swagger.Feed{
			Title:       feed.Title,
			Link:        feed.Link,
			Description: feed.Description,
			Categories:  &feed.Categories,
			Image: Swagger.Image{
				Title: &feed.Image.Title,
				Url:   &feed.Image.URL,
			},
			Items: swagItems,
		}

		// Uses title as the map key, assumes title is unique
		swagFeeds.AdditionalProperties[swagFeed.Title] = swagFeed

		// Store titles so they can be used for the GetTitles endpoint
		titles[swagFeed.Title] = true
	}

	return swagFeeds, titles, categories
}

// Extract the first category if it exists
func extractCategory(categories []string) string {
	if len(categories) == 0 {
		return ""
	}
	return categories[0]
}

// Convert a map to a slice
func convertMapToSlice(m map[string]bool) []string {
	s := make([]string, 0)
	for k := range m {
		s = append(s, k)
	}

	return s
}

// Convert a slice to a map
func convertSliceToMap(s []string) map[string]bool {
	m := make(map[string]bool, 0)
	for _, v := range s {
		m[v] = true
	}

	return m
}
