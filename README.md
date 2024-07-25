# Go Elastically - Meetup Toronto Live Coding Sessions

This repository contains the code samples demonstrated during the Elastic Meetup from July 25th 2024. The primary focus is to showcase the usage of the Elasticsearch Go client for indexing and searching XKCD comics.

![KXCD Standards](https://imgs.xkcd.com/comics/standards.png)

https://www.meetup.com/elasticsearch-toronto/events/301905900/

## Table of Contents
- [Overview](#overview)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - [Indexer](#indexer)
  - [Searcher](#Searcher)
- [Contributing](#contributing)
- [License](#license)

## Overview

In this repository, you will find two main code samples:

1. **Inderxer**: A Go program that fetches XKCD comics metadata and indeexes them into an elasticsearch cluster.
2. **Searcher**: A Go program that executes quries to pull stored XKCD comics from the Elasticsearch XKCD index.

These examples illustrate how to use the Elasticsearch Go client to efficiently index and search data.

## Getting Started

To get started with the code samples, follow the instructions below.

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

![Demo](out.gif)

### Searcher

The Indexer program queries an Elasticsearch index for XKCD comics and loops through the results. It configures an Elasticsearch client, and executes a search query for XKCD comics. The response is decoded, and logs each comic, such as ID, score, alt text, day, news, number, and title.

1. Ensure your Elasticsearch instance is running. You can use Docker to start an instance:

    ```sh
    docker run -d -p 9200:9200 -e "discovery.type=single-node" elasticsearch:7.10.1
    ```

2. Run the uploader:

    ```sh
    go run uploader/main.go
    ```

3. The data from `comics.json` will be indexed into Elasticsearch, and you can use Kibana or any other client to search for the comics.

## Contributing

We welcome contributions from the community. If you have suggestions or improvements, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Happy coding! If you have any questions or need further assistance, feel free to open an issue or contact the repository maintainers.

---

**Maintainers:**
- Your Name ([@yourusername](https://github.com/yourusername))

---

### Additional Resources

- [Elasticsearch Go Client Documentation](https://github.com/elastic/go-elasticsearch)
- [XKCD API Documentation](https://xkcd.com/json.html)
