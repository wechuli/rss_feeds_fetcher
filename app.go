package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rss_feeds/models"
	"rss_feeds/rss"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// App struct, encapsulates the router and DB references
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Search struct
type Search struct {
	Term string `json:"term"`
}

// Initialize sets the DB and Router reference for the App
func (a *App) Initialize(host string, port int, user string, password string, dbname string) {

	var err error
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	go a.FetchNewFeedsAndPopulateDB() // go routine to fetch new rss feeds and store in DB every time the program is run

	a.Router = mux.NewRouter()
	a.initializeRoutes()

}

// Run initiates listening on the specified ports and registers the mux router
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {

	a.Router.HandleFunc("/search", a.searchFeeds).Methods("POST")

}

// respondWithJSON - utilility function to return JSON payload
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondWithError - utility function to respond with an Error JSON Payload
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// searchFeeds -
func (a *App) searchFeeds(w http.ResponseWriter, r *http.Request) {

	var s Search

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&s); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	feeds, err := models.SearchRssFeeds(a.DB, s.Term)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, feeds)
}

// FetchNewFeedsAndPopulateDB loops through a list of specified urls, gets the feeds and populates the DB with *new* rss feeds
func (a *App) FetchNewFeedsAndPopulateDB() {

	// CNN and BBC rss feeds links for various news categories
	rssFeedsURLs := []string{"http://rss.cnn.com/rss/edition.rss", "http://rss.cnn.com/rss/edition_world.rss", "http://rss.cnn.com/rss/edition_africa.rss", "http://rss.cnn.com/rss/edition_technology.rss", "http://feeds.bbci.co.uk/news/rss.xml", "http://feeds.bbci.co.uk/news/world/rss.xml", "http://feeds.bbci.co.uk/news/technology/rss.xml", "http://feeds.bbci.co.uk/news/world/africa/rss.xml"}

	for _, url := range rssFeedsURLs {
		rawString, err := rss.FetchRssFeedRaw(url)
		if err != nil {
			fmt.Println(err)
		}
		feeds, err := rss.ParseRawRssString(rawString)
		if err != nil {
			fmt.Println(err)
		}

		// populate db
		err = models.StoreRssFeeds(a.DB, feeds)
		if err != nil {
			fmt.Println(err)
		}
	}


	fmt.Println("Data fetched and updated in database")
}
