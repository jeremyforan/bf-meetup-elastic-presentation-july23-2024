# Go Elastically - Meetup Toronto Live Coding Sessions

This repository contains the code samples demonstrated during the [Toronto Elastic Meetup from July 25th 2024](https://www.meetup.com/elasticsearch-toronto/events/301905900/). The primary focus is to showcase the usage of the Elasticsearch Go client for indexing and searching XKCD comics.


![KXCD Standards](https://imgs.xkcd.com/comics/standards.png)

[**XKCD**](https://xkcd.com/) - a webcomic of romance, sarcasm, math, and language

## Table of Contents
- [Overview](#overview)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - [Indexer](#indexer)
  - [Searcher](#Searcher)
- [License](#license)

## Overview

In this repository, you will find two main code samples in the *cmd* directory:

1. **Indexer**: A Go program that fetches XKCD comics metadata and indexes them into an elasticsearch cluster.
2. **Searcher**: A Go program that executes quries to pull stored XKCD comics from the Elasticsearch XKCD index.

These examples illustrate how to use the Elasticsearch Go client to efficiently index and search data.

## Getting Started

To get started with the code samples, follow the instructions below.

If you have the time I highly recommend you watch [The Go Language](https://www.youtube.com/watch?v=rKnDgT73v8s) video introducing Go to the world.

### API

Each XKCD can be fetched by appending the comic number to the base url
```
https://xkcd.com/{comic_number}
```

To get a json description of that comic, we append /info.0.json to the url
```
https://xkcd.com/{comic_number}/info.0.json
```

This URL will fetch the JSON for the 'Standards' comic
```
https://xkcd.com/2963/info.0.json
```

and return the following
```json
{
  "month": "7",
  "num": 2963,
  "link": "",
  "year": "2024",
  "news": "",
  "safe_title": "House Inputs and Outputs",
  "transcript": "",
  "alt": "People think power over ethernet is so great, and yet when I try to do water over ethernet everyone yells at me.",
  "img": "https://imgs.xkcd.com/comics/house_inputs_and_outputs.png",
  "title": "House Inputs and Outputs",
  "day": "24"
}
```


### Prerequisites

- Go (1.21+)
- Elasticsearch instance (You can use the official [Docker image](https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html), or [Elastic Cloud](https://www.elastic.co/cloud))
- An Elastic index to store the XKCD comics:
```json
{
  "mappings": {
    "properties": {
      "month": {
        "type": "keyword"
      },
      "num": {
        "type": "integer"
      },
      "link": {
        "type": "keyword"
      },
      "year": {
        "type": "keyword"
      },
      "news": {
        "type": "text"
      },
      "safe_title": {
        "type": "keyword"
      },
      "transcript": {
        "type": "text"
      },
      "alt": {
        "type": "text"
      },
      "img": {
        "type": "keyword"
      },
      "title": {
        "type": "text"
      },
      "day": {
        "type": "keyword"
      }
    }
  }
}
```

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/jeremyforan/bf-meetup-elastic-presentation-july23-2024.git
    cd bf-meetup-elastic-presentation-july23-2024
    ```

2. Install the required Go packages:

    ```sh
    go mod tidy
    ```

## Usage

### Indexer

The indexer fetches [XKCD](https://xkcd.com/) comics and indexes them into Elasticsearch. It configures an Elasticsearch client and uses goroutines to download comics asynchronously. The comics are stored in a thread-safe structure. After downloading, the comics are indexed into Elasticsearch using a bulk indexer, which flushes data to Elasticsearch periodically. The program logs progress and errors throughout.

![Indexer Demo](indexer.gif)

### Searcher

The Searcher program queries an Elasticsearch index for XKCD comics and loops through the results. It configures an Elasticsearch client, and executes a search query for XKCD comics. The response is decoded, and logs each comic, such as ID, score, alt text, day, news, number, and title.

![Searcher Demo](searcher.gif)

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Happy coding! If you have any questions or need further assistance, feel free to open an issue or contact the repository maintainers.

And joiun the Elasticsearch '#meetup-toronto' channel

---

**Maintainers:**
- Jeremy Foran ([@jeremyforan](https://github.com/jeremyforan))

---

### Additional Resources

- [Elasticsearch Go Client Documentation](https://github.com/elastic/go-elasticsearch)
- [XKCD API Documentation](https://xkcd.com/json.html)
