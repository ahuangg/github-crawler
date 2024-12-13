package crawler

import (
	"fmt"
	"sort"

	"net/url"
	"os"
	"sync"
	"time"

	"github.com/ahuangg/gh-crawler/internal/models"
	"github.com/ahuangg/gh-crawler/internal/utils"
	"github.com/gocolly/colly/v2"
)

type Crawler struct {
	collector *colly.Collector
	baseUrl string
}

func NewCrawler() *Crawler{
	outputDir := "locations"
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        utils.PrintError("failed to create locations directory: %v", err)
    }
	
	c := colly.NewCollector(
		colly.AllowedDomains("github.com"),
		colly.UserAgent(GetRandomUserAgent()),
	)

	c.Limit(&colly.LimitRule{
        DomainGlob:  "*github.com*",
        Delay:       1 * time.Second,
    })

	return &Crawler{collector: c, baseUrl: "https://github.com"}
}

func (c *Crawler) CrawlUsersByLocation(location string, userChannel chan<- *models.User) {
    page := 1
    var crawlWg sync.WaitGroup
    var callbackWg sync.WaitGroup
    hasUsers := false

    c.collector = colly.NewCollector(
        colly.AllowedDomains("github.com"),
        colly.UserAgent(GetRandomUserAgent()),
    )

    c.collector.Limit(&colly.LimitRule{
        DomainGlob:  "*github.com*",
        Delay:       1 * time.Second,
    })

    c.collector.OnHTML("div.Box-sc-g0xbh4-0.flszRz", func(e *colly.HTMLElement) {
        username := e.ChildText("h3 span.Text__StyledText-sc-17v1xeu-0.hBjWst")
  
        if username != "" {
            hasUsers = true
            callbackWg.Add(1)
            user := &models.User{
                Username: username,
                Location: "",
                LanguageStats: make(map[string]float64),
            }
            userChannel <- user
            utils.PrintUserRetrieved("%s", username)
            callbackWg.Done()
        }
    })

    c.collector.OnScraped(func(r *colly.Response) {
        time.Sleep(1 * time.Second)
        crawlWg.Done()
    })

    for {
        hasUsers = false
        crawlWg.Add(1)
        searchURL := fmt.Sprintf("https://github.com/search?q=location:%s&type=users&p=%d", 
            url.QueryEscape(location), page)
        utils.PrintInfo("Searching %s", searchURL)
        
        err := c.collector.Visit(searchURL)
        if err != nil {
            utils.PrintError("Error visiting page: %v", err)
            crawlWg.Done()
            break
        }
        
        crawlWg.Wait()

        if !hasUsers {
            break
        }

        page++
    }

    callbackWg.Wait()
    close(userChannel)
}

func (c *Crawler) ProcessUserLanguageStats(user *models.User, processedChannel chan<- *models.User) {
    collector := colly.NewCollector(
        colly.AllowedDomains("github.com"),
        colly.UserAgent(GetRandomUserAgent()),
    )

    collector.Limit(&colly.LimitRule{
        DomainGlob:  "*github.com*",
        Delay:       1 * time.Second,
    })

    profileURL := fmt.Sprintf("https://github.com/%s", user.Username)
    collector.OnHTML("li.vcard-detail[itemprop='homeLocation']", func(e *colly.HTMLElement) {
        location := e.ChildText("span.p-label")
        if location != "" {
            user.Location = location 
        }
    })

    err := collector.Visit(profileURL)
    if err != nil {
        utils.PrintError("Error visiting profile page for %s: %v", user.Username, err)
        return
    }

    if user.Location == "" {
        utils.PrintInfo("Skipping %s - Missing location", user.Username)
        return 
    }

    languageCounts := make(map[string]int)
    totalRepos := 0
    page := 1
    hasNextPage := true

    collector.OnHTML("div[id^='user-repositories-list'] div.col-10", func(e *colly.HTMLElement) {
        isFork := e.ChildText("span.Label") == "Fork"
        if !isFork {
            language := e.ChildText("span[itemprop='programmingLanguage']")
            if language != "" {
                languageCounts[language]++
                totalRepos++
            }
        }
    })

    collector.OnHTML("div.paginate-container", func(e *colly.HTMLElement) {
        nextLink := e.ChildAttr("a[rel='next']", "href")
        if nextLink == "" {
            hasNextPage = false
        }
    })

    for hasNextPage && page <= 10 {
        reposURL := fmt.Sprintf("https://github.com/%s?page=%d&tab=repositories", user.Username, page)
        err := collector.Visit(reposURL)
        if err != nil {
            utils.PrintError("Error visiting %s page %d: %v", user.Username, page, err)
            break
        }
        page++
        time.Sleep(1 * time.Second)
    }

    if totalRepos > 0 {
        type langStat struct {
            Language   string
            Percentage float64
        }
        stats := make([]langStat, 0)

        for lang, count := range languageCounts {
            percentage := (float64(count) / float64(totalRepos)) * 100
            stats = append(stats, langStat{lang, percentage})
        }

        sort.Slice(stats, func(i, j int) bool {
            return stats[i].Percentage > stats[j].Percentage
        })

        user.LanguageStats = make(map[string]float64)
        for _, stat := range stats {
            user.LanguageStats[stat.Language] = stat.Percentage
        }
    }

    utils.PrintUserProcessed("%s", user.Username)
    processedChannel <- user
}