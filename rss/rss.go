

package rss

import (
	"errors"
	"io/ioutil"
	"net/http"
	"github.com/wechuli/rss_feeds_fetcher/models"
	"strings"

	"github.com/mmcdole/gofeed/rss"
)


// FetchRssFeedRaw downloads a raw rss website and returns a string representation
func FetchRssFeedRaw(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {

		return "", errors.New("unable to fetch feed")
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("unable to  feed")
	}
	return string(bodyBytes), nil

}

// ParseRawRssString takes a string repr of an rss website and parses it to a slice of feeds and returns the feeds
func ParseRawRssString(rawRssString string) ([]models.Feed, error) {

	fp := rss.Parser{}
	rssFeed, err := fp.Parse(strings.NewReader(rawRssString))

	if err != nil {
		return nil, err
	}
	
	var feeds []models.Feed

	for _, item := range rssFeed.Items {

		feed := models.Feed{Title: item.Title, Description: item.Description, PubDate: item.PubDate, Link: item.Link}

		feeds = append(feeds, feed)
	}

	return feeds, nil

}
