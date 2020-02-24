## Introduction

The project fetches RSS Feeds from BBC and CNN and enables the feeds to be accessed and searched through a RESTful web service. Feeds are stored in a Postgres database and full text search can be performed on them. The search functionality is exposed as an API endpoint that accepts a search keyword and returns a JSON array of RSS feeds that matched the search criteria.

## Running the Project

1. The project makes use of a Postgres database to store the rss feeds and uses Postgres's full text search capability to search through the feeds . A running Postgres database will need to be set up before the application is run. A sample script to create a Postgres database is provided below:

```SQL


CREATE DATABASE rssfeeds;
CREATE USER rssuser WITH PASSWORD 'password';

ALTER ROLE rssuser SET client_encoding TO 'utf8';
ALTER ROLE rssuser SET default_transaction_isolation TO 'read committed';
ALTER ROLE rssuser SET timezone TO 'UTC';


GRANT ALL PRIVILEGES ON DATABASE rssfeeds TO rssuser;

CREATE TABLE feeds(
    title text NOT NULL ,
    description text,
    pubdate varchar,
    link varchar PRIMARY KEY
);


```

The table does not need to be created before-hand, once the application is started, the necessary table will be created if it does not exist.

2.  The project can be run using:

        go run main.go app.go

This initializes a http server on port 8080 while concurrently fetching new RSS Feeds from CNN and BBC and storing them on the database.

The database connection details are queried from the environment variables, so they will need to be setup before hand. Refer to code snippet below to name the environment variables correctly:

```GO
host, user, password, dbname := os.Getenv("RSS_DB_HOST"), os.Getenv("RSS_DB_USERNAME"), os.Getenv("RSS_DB_PASSWORD"), os.Getenv("RSS_DB_NAME")
```

3. There is only one route - for searching for the feeds `/search` which allows a `POST` method with a JSON body describing the search term

**sample request**

`POST /search`

```JSON
{
	"term":"africa"
}
```

**sample response**

```JSON

[
    {
        "title": "CNN Travel's 20 best places to visit in 2020",
        "description": "Whether you want to relax on a remote island off the coast of Africa, ride Germany's coolest trains or spot howling monkeys in South America, there is much to explore heading into a new decade in 2020.<img src=\"http://feeds.feedburner.com/~r/rss/edition_world/~4/AmgOSHUmQzE\" height=\"1\" width=\"1\" alt=\"\"/>",
        "pubdate": "Mon, 06 Jan 2020 23:10:01 GMT",
        "link": "http://rss.cnn.com/~r/rss/edition_world/~3/AmgOSHUmQzE/index.html"
    },
    {
        "title": "Britain seeks closer economic ties with Africa following Brexit ",
        "description": "Britain is leaving the European Union on Friday, starting the clock on an 11-month transition period during which the country will try to sign as many new trade deals as possible. African countries are a prime target.<img src=\"http://feeds.feedburner.com/~r/rss/edition_africa/~4/VrZIYMYDIxU\" height=\"1\" width=\"1\" alt=\"\"/>",
        "pubdate": "Fri, 31 Jan 2020 16:23:35 GMT",
        "link": "http://rss.cnn.com/~r/rss/edition_africa/~3/VrZIYMYDIxU/index.html"
    }
]

```

## Application Structure

To separate concerns, the project was implemented using two packages(additional to `main`) :- `models` and `rss`

### Package rss

Package rss implements functions concerned with fetching and parsing rss feeds from CNN and BBC (it can fetch rss feeds from any rss enabled website) and returning it in a format ready to be stored in the database (where the `models` package takes over)

- **`func FetchRssFeedRaw`**

  ```GO

    func FetchRssFeedRaw(url string) (string, error)
  ```

  `FetchRssFeedRaw` fetches the raw rss feeds (usually in xml) and returns a string representation of the same.

- **`func ParseRawRssString`**

      ```GO
      func ParseRawRssString(rawRssString string) ([]models.Feed, error)

`ParseRawRssString` receives a raw string rss website and parses it to extract the rss feeds represented in a convenient `Feed` struct, defined in the package `models`. It returns a slice of feeds extracted from the website.

### Package models

Package models implements types and functions related to reading and writing the rss feeds to the database. It expects feeds to have already been formatted correctly (from the rss package).

- **`type Feed`**

  ```GO
  type Feed struct {
  `
  `
  `
  Link        string `json:"link"`
  }

  ```

  `Feed` struct is a minimal representation of the rss feeds, with just enough information to enable searching. JSON tags are included to allow marshalling and unmarshalling when working with the web servive

- **`func StoreRssFeeds`**

  ```GO
  func StoreRssFeeds(db *sql.DB, feeds []Feed) error{...}

  ```

  `StoreRssFeeds` takes a slice of `Feeds` and runs a query on the database(pointed at by the `db` parameter) to store the Feeds.

- **`func SearchRssFeeds`**

  ```GO

  func SearchRssFeeds(db *sql.DB, searchString string) ([]Feed, error){...}

  ```

  `SearchRssFeeds` performs a full text search on the database (pointed at by the `db` parameter) for the phrase specified in the `searchString` string parameter and returns a slice of Feeds that match the search criteria.

## Testing

To test the project run `go test -v`.
