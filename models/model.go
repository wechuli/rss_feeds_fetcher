package models

import (
	"database/sql"
)

// Feed struct
type Feed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	PubDate     string `json:"pubdate"`
	Link        string `json:"link"`
}

// query to create the table
const tableCreationQuery = `CREATE TABLE IF NOT EXISTS feeds
(
    title text NOT NULL ,
    description text,
    pubdate varchar,
    link varchar PRIMARY KEY
)`

// SearchRssFeeds fetch all rss feeds matching a search string
func SearchRssFeeds(db *sql.DB, searchString string) ([]Feed, error) {

	rows, err := db.Query("SELECT title,description,pubdate,link FROM feeds where to_tsvector(title ||' '|| description) @@ phraseto_tsquery($1)", searchString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	feeds := []Feed{}

	for rows.Next() {
		var f Feed
		if err := rows.Scan(&f.Title, &f.Description, &f.PubDate, &f.Link); err != nil {
			return nil, err
		}
		feeds = append(feeds, f)
	}

	return feeds, nil
}

//StoreRssFeeds stores feeds in the db
func StoreRssFeeds(db *sql.DB, feeds []Feed) error {
	if _, err := db.Exec(tableCreationQuery); err != nil {
		return err
	}
	for _, item := range feeds {
		// skip if a row already contains the link - duplicated rss feed
		rows, err := db.Query("INSERT INTO feeds(title,description,pubdate,link) VALUES($1,$2,$3,$4) ON CONFLICT (link) DO NOTHING", item.Title, item.Description, item.PubDate, item.Link)
		if err != nil {
			return err
		}
		defer rows.Close()

	}

	return nil
}
