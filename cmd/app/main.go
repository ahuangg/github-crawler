package main

import (
	"bufio"
	"os"
	"sync"

	"github.com/ahuangg/gh-crawler/internal/crawler"
	"github.com/ahuangg/gh-crawler/internal/models"
	"github.com/ahuangg/gh-crawler/internal/utils"
	"github.com/ahuangg/gh-crawler/internal/writer"
)

func main() {
    c := crawler.NewCrawler()
    
    file, err := os.Open("locations.txt")
    if err != nil {
        utils.PrintError(err.Error())
        return
    }
    defer file.Close()

    var locations []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        locations = append(locations, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        utils.PrintError(err.Error())
        return
    }

    maxConcurrentProcessing := 25

    csvWriter, err := writer.NewCSVWriter("locations", "data")
    if err != nil {
        utils.PrintError(err.Error())
    }

    for _, location := range locations {
        userChannel := make(chan *models.User)
        processedChannel := make(chan *models.User)
        var wg sync.WaitGroup

        utils.PrintInfo("Starting to Process %s", location)

        go func(loc string) {
            c.CrawlUsersByLocation(loc, userChannel)
        }(location)

        processingWg := sync.WaitGroup{}
        for i := 0; i < maxConcurrentProcessing; i++ {
            processingWg.Add(1)
            go func() {
                defer processingWg.Done()
                for user := range userChannel {
                    c.ProcessUserLanguageStats(user, processedChannel)
                }
            }()
        }

        go func() {
            processingWg.Wait()
            close(processedChannel)
        }()

        wg.Add(1)
        go func() {
            defer wg.Done()
            for user := range processedChannel {
                if err := csvWriter.WriteUser(user); err != nil {
                    utils.PrintError("%v", err)
                } else {
                    utils.PrintUserWritten("%s - %s", user.Username, user.Location)
                }
            }
        }()

        wg.Wait()

        if err := csvWriter.Close(); err != nil {
            utils.PrintError("%s - %v", location, err)
        }
    }

    utils.PrintSuccess("Completed")
}