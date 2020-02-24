package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/wechuli/rss_feeds_fetcher/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "rssuser"
	password = "password"
	dbname   = "rssfeeds"
)

var app App

var sampleFeeds []models.Feed

func TestMain(m *testing.M) {
	app = App{}
	app.Initialize(host, port, user, password, dbname)

	code := m.Run()
	clearTable()
	os.Exit(code)

}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d \n", expected, actual)
	}
}

func clearTable() {
	app.DB.Exec("DELETE FROM feeds")
}
func addSampleFeeds() {
	feeds := []models.Feed{
		{Title: "Coronavirus: Evacuation flight for Britons on Diamond Princess lands in UK", Description: "British nationals evacuated from a cruise ship are on their way to a hospital, where they will be quarantined.", PubDate: "Sat, 22 Feb 2020 10:08:39 GMT", Link: "https://www.cnn.com/2020/02/22/us/iyw-nothing-but-love-notes-trnd/index.html"},
		{Title: "Kenya's longest serving president dies at 95", Description: "Former Kenyan President Daniel Arap Moi, who ruled the country for 24 years has died, President Uhuru Kenyatta announced Tuesday", PubDate: "Tue, 04 Feb 2020 15:56:27 GMT", Link: "http://rss.cnn.com/~r/rss/edition_africa/~3/Z7oI1OZU1uM/index.html"}}
	sampleFeeds = append(sampleFeeds, feeds...)

	models.StoreRssFeeds(app.DB, sampleFeeds)
}

func TestDBIsPopulatedByFeeds(t *testing.T) {
	app.FetchNewFeedsAndPopulateDB()
	var count int
	app.DB.QueryRow(`SELECT count(link) from feeds;`).Scan(&count)
	if count == 0 {
		t.Errorf("Expected rss feed count in the database not to be 0 \n")
	}
}

func TestEmptyFeedsTable(t *testing.T) {
	clearTable()

	payload := []byte(`{"term":"Coronavirus"}`)

	req, _ := http.NewRequest("POST", "/search", bytes.NewBuffer(payload))

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	strBody := string(response.Body.Bytes())

	if strBody != "[]" {
		t.Errorf("Expected body to be %s. Got %s \n", "[]", strBody)
	}

}

func TestSearchForExistingFeed(t *testing.T) {
	clearTable()
	addSampleFeeds()

	payload := []byte(`{"term":"Coronavirus"}`)

	req, _ := http.NewRequest("POST", "/search", bytes.NewBuffer(payload))

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	// strBody := string(response.Body.Bytes())

	var f []models.Feed

	err := json.NewDecoder(response.Body).Decode(&f)
	if err != nil {
		fmt.Println(err)
	}

	// check that the first rss feed was received back from the api matches the one in the sampleFeeds

	if f[0].Link != sampleFeeds[0].Link {
		t.Errorf("Expected link to be %s. Got %s \n", sampleFeeds[0].Link, f[0].Link)
	}

}
