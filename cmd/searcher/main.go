package main

import (
	"encoding/json"
	elasticsearch8 "github.com/elastic/go-elasticsearch/v8"
	"io"
	"log/slog"
	"strings"
)

func main() {
	slog.Info("Starting the XKCD comic query")

	// create a new elasticsearch client configuration
	cfg := elasticsearch8.Config{
		Addresses: []string{
			"https://elastic.jeremyforan.com:443",
		},
		Username: "demo",
		Password: "demo123",
	}

	slog.Info("Create an Elasticsearch client")

	// create a new elasticsearch client
	es, err := elasticsearch8.NewClient(cfg)
	if err != nil {
		slog.Error("could not create client", "error", err)
		return
	}

	// create a query to search for all the XKCD comics
	query := `{ "size": 5, "query": { "match_all": {} } }`
	reader := strings.NewReader(query)

	slog.Info("Initiating search")

	res, err := es.Search(
		es.Search.WithIndex("xkcd"),
		es.Search.WithBody(reader),
	)
	if err != nil {
		slog.Error("unable to perform search", "error", err)
		return
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			slog.Error("unable to close body", "error", err)
		}
	}(res.Body)

	if res.IsError() {
		slog.Error("response reported an error", "error", res.String())
		return
	}

	var data ElasticSearchResponse

	// Deserialize the response into a map.
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		slog.Error("unable to decode response", "error", err)
		return
	}

	slog.Info("Results found", "total", data.Hits.Total.Value)

	// print out the results
	for _, hit := range data.Hits.Hits {
		slog.Info("////////////////////////////////////////////////////////////////////////////////////////////////")
		slog.Info("\tID", "id", hit.ID)
		slog.Info("\tScore", "score", hit.Score)
		slog.Info("\tAlt", "alt", hit.Source.Alt)
		slog.Info("\tDay", "day", hit.Source.Day)
		slog.Info("\tNews", "news", hit.Source.News)
		slog.Info("\tNum", "num", hit.Source.Num)
		slog.Info("\tTitle", "title", hit.Source.Title)
	}
}

// XKCD struct represents an XKCD comic.
type XKCD struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
}

// ElasticSearchResponse struct represents the response from an ElasticSearch query. It has multiple nested levels
// representing the shards, hits, and the actual hits.
// More information below.
type ElasticSearchResponse struct {
	Shards   Shards   `json:"_shards"`
	Hits     HitsHigh `json:"hits"`
	TimedOut bool     `json:"timed_out"`
	Took     float64  `json:"took"`
}

type HitsHigh struct {
	Hits     []Hit   `json:"hits"`
	MaxScore float64 `json:"max_score"`
	Total    Total   `json:"total"`
}

type Hit struct {
	ID     string  `json:"_id"`
	Index  string  `json:"_index"`
	Score  float64 `json:"_score"`
	Source XKCD    `json:"_source"`
}

type Total struct {
	Relation string  `json:"relation"`
	Value    float64 `json:"value"`
}

type Shards struct {
	Failed     int64 `json:"failed"`
	Skipped    int64 `json:"skipped"`
	Successful int64 `json:"successful"`
	Total      int64 `json:"total"`
}
