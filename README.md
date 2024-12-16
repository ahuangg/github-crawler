[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)

# Github User Crawler

This program efficiently crawls and collects users programming language data from **multiple** Github users based on geographic locations. It processes **one location at a time** (such as San Francisco, Tokyo, or Singapore) and writes the results to **CSV files** organized by region. The collected data includes comprehensive **user information** such as usernames, locations, and detailed programming language statistics from their repositories. Using the Colly web scraping framework, the crawler implements smart rate limiting, concurrent processing, and automatic retries to respect Github's website guidelines. It's easily configurable through a simple locations text file where you can specify target locations and the number of users to crawl per location.

## Key Features

-   Geographic-based GitHub user crawling
-   Concurrent processing with rate limiting
-   Language statistics calculation for each user
-   CSV output organized by regions
-   Configurable user counts per location
-   Random user agent rotation for reliable crawling
-   Automatic retry mechanism
-   Progress tracking and logging

## Prerequisites

-   Go 1.20 or higher
-   Git
-   GitHub access (no authentication required)

## Output Format

The crawler generates CSV files in the `locations` directory with the following structure:

### User Information Fields

-   username: GitHub username (e.g., "ahuangg")
-   location: User's listed location (e.g., "New York City, NY")
-   language_stats: Language statistics JSON object (e.g., {"Python": "49", "TypeScript": "11", "Go": "40"})

Example CSV output:

## Example

If you want to know the specific results of this program, you can view the locations directory of this project. which has examples of data crawled and stored in the format described above.

## Operating Environment

-   Development language: Go
-   System: Windows/Linux/macOS

## Configuration

### Location File Format

The `locations.txt` file supports two types of entries:

-   Region headers: Lines starting with # (e.g., `# San Francisco Bay Area`)
-   Location entries: `<city_name> <target_count>` (e.g., `San Francisco 255`)
-   Target user count determines how many GitHub users to crawl from that location
-   A new csv file is created for each region header

### Rate Limiting

The crawler implements:

-   1 second delay between requests
-   Concurrent processing of up to 5 users (can be updated in main by changing maxConcurrentProcessing)

### Source Code Installation

```bash
$ git clone git@github.com:ahuangg/github-crawler.git
$ cd github-crawler
$ go install
```

#### Building and running the crawler

```bash
$  go build -o bin/github-crawler ./cmd/app
$  ./bin/github-crawler
```

## Notes

Please open an issue if you encounter any bugs or would like to request new features for this project
