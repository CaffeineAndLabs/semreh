package main

import (
	"fmt"
	"html"
	"strings"
	"time"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/mmcdole/gofeed"
)

type FeedItem struct {
	Source          string
	Title           string
	Description     string
	Link            string
	PublishedParsed *time.Time
}

func sanitizeHTML(s string) string {
	sanitized := strip.StripTags(s)
	sanitized = html.UnescapeString(sanitized)

	return sanitized
}

func truncDescription(s string) string {
	// If there is '\n' it's probably a big description
	// So we need to truncate it
	if strings.ContainsAny(s, "\n") {
		trunc := strings.SplitN(s, "\n", 2)[0]
		return fmt.Sprintf("%s ...\n", trunc)
	}

	return s
}

func getLastNews(n int) []*FeedItem {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("http://www.dofus.com/fr/rss/news.xml")
	var lastNews []*FeedItem

	for i := 0; i <= n-1; i++ {
		new := &FeedItem{
			Source:          feed.Title,
			Title:           feed.Items[i].Title,
			Description:     truncDescription(sanitizeHTML(feed.Items[i].Description)),
			Link:            feed.Items[i].Link,
			PublishedParsed: feed.Items[i].PublishedParsed,
		}
		lastNews = append(lastNews, new)
	}

	return lastNews
}
