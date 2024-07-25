package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	elasticsearch8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/schollz/progressbar/v3"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// maxID represents the maximum XKCD comic ID. This increments each day as a new
const maxID = 2962

//the url for the XKCD comic is created by appending the comic number to the base url
// https://xkcd.com/{comic_number}

// To get a json description of that comic, we append /info.0.json to the url
// https://xkcd.com/{comic_number}/info.0.json

func main() {
	slog.Info("Starting XKCD comic downloader")

	// create a new wait group to wait for all comic pulling go routines to complete.
	var wg sync.WaitGroup

	// create a new XKCD store to hold all the XKCD comics
	var data XkcdStore

	// create a new elasticsearch client configuration
	// details on the client can be found at https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/getting-started-go.html
	esCred := elasticsearch8.Config{
		Addresses: []string{
			"https://elastic.jeremyforan.com:443",
		},
		Username: "demo",
		Password: "demo123",
	}

	slog.Info("Creating an Elasticsearch client")

	// create a new elasticsearch client
	es, err := elasticsearch8.NewClient(esCred)
	if err != nil {
		slog.Error("could not create client", "error", err)
		return
	}

	// create a new reader and item to be used by the bulk indexer
	var reader io.ReadSeeker
	var item esutil.BulkIndexerItem

	slog.Info("Spawning download threads", "threads", maxID-1)

	// I know this is not necessary, but I couldn't help myself
	bar := progressbar.Default(maxID - 2)

	// loop through all the XKCD comics, starting from 1 to the maxID
	for i := 1; i <= maxID; i++ {

		// there is no xkcd with id 404, if you try to grab it you get a 404. ;)
		if i == 404 {
			continue
		}

		var comic XKCD
		go func(id int) {
			// increment the wait group counter
			wg.Add(1)

			// decrement the wait group counter when the go routine completes
			defer wg.Done()

			// fetch the XKCD comic from the XKCD API and store it in the XKCD store
			err = GetJSON(makeURL(id), &comic)

			// if there is an error, log it
			if err != nil {
				slog.Error("error", err)
				return
			} else {

				// lock the data store to ensure that only one go routine can write to it at a time,
				// then unlock after the comic has been added
				data.Lock()
				data.XKCDs = append(data.XKCDs, comic)
				data.Unlock()
				bar.Add(1)
			}
		}(i)
	}

	// wait until all threads have completed downloading the XKCD comics
	wg.Wait()

	// close the progress bar
	bar.Close()

	// create a new bulk indexer to index the XKCD comics into elasticsearch
	bulkConfig := esutil.BulkIndexerConfig{
		Client: es,
		Index:  "xkcd",

		// how many workers to use to index the data. runtime.NumCPU() is the default.
		NumWorkers: 4,

		// how many bytes the bulk indexer should be before it flushes the data to elasticsearch
		FlushBytes: 5e+6,

		// how often the bulk indexer should flush the data to elasticsearch
		FlushInterval: 250 * time.Millisecond,
	}

	// create a new bulk indexer
	indexer, err := esutil.NewBulkIndexer(bulkConfig)
	if err != nil {
		slog.Error("unable to create bulk indexer", "error", err)
		return
	}

	slog.Info("Indexing XKCD comics into Elasticsearch")

	// loop through all the XKCD comics and add each one to the bulk indexer
	for _, xkcd := range data.XKCDs {

		// convert the XKCD comic to a reader that can be passed to the bulk indexer
		reader, err = xkcd.Reader()
		if err != nil {
			slog.Error("unable to create reader for comic", "error", err)
			continue
		}

		// create a new bulk indexer item
		// more info here https://github.com/elastic/go-elasticsearch/tree/main/_examples/bulk#indexergo
		item = esutil.BulkIndexerItem{
			Action: "index",
			Index:  "xkcd",

			// the document id can be specified here, if not specified, elasticsearch will generate one.
			// the comic number is a natural key which can be used for the document id
			DocumentID: fmt.Sprintf("%d", xkcd.Num),

			// Pass the reader to the bulk indexer
			Body: reader,
		}

		// Add the item to the indexer with the background context. Context can be used to cancel operations if needed.
		err = indexer.Add(context.Background(), item)
		if err != nil {
			slog.Error("unable to add document to bulk indexer", "error", err)
		}
	}

	// Close waits until all added items are flushed and closes the indexer. Similar to a Sync Wait operation.
	err = indexer.Close(context.Background())
	if err != nil {
		slog.Error("unable to close indexer", "error", err)
	}

	slog.Info("Indexing complete")

	// get the stats of the bulk indexer
	stats := indexer.Stats()
	slog.Info("bulk indexer report", "stats", stats)
}

// makeURL creates a url for the XKCD comic with the given id
// given the id of 123, the url will be https://xkcd.com/123/info.0.json
func makeURL(id int) string {
	return fmt.Sprintf("https://xkcd.com/%d/info.0.json", id)
}

// GetJSON fetches a JSON object from the given url and stores it in the target interface
func GetJSON(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	// close the response body when the function completes
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			slog.Error("unable to close response reader", "error", err)
		}
	}(resp.Body)

	// decode the JSON object from the response body and store it in the target interface
	return json.NewDecoder(resp.Body).Decode(target)
}

// XKCD represents the JSON structure of an XKCD comic. The data can be pulled from https://xkcd.com/1/info.0.json
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

// XkcdStore represents a store of XKCD comics
type XkcdStore struct {
	// hold all the XKCD comics in this slice
	XKCDs []XKCD

	// add the ability to lock the store to prevent concurrent writes corrupting the
	sync.Mutex
}

// Reader converts the XKCD comic to a reader that can be passed to the bulk indexer
func (x XKCD) Reader() (io.ReadSeeker, error) {

	// Marshal the XKCD comic to a byte slice
	data, err := json.Marshal(x)
	if err != nil {
		slog.Error("unable to marshal comic to byte slice", "error", err)
		return nil, err
	}

	// Create a bytes.Reader from the byte slice
	reader := bytes.NewReader(data)

	return reader, nil
}
