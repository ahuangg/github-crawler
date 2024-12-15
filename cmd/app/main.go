package main

import (
	"bufio"
	"os"
	"strings"
	"strconv"
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

	var currentWriter *writer.CSVWriter
	scanner := bufio.NewScanner(file)
	
	maxConcurrentProcessing := 5

	for scanner.Scan() {
		line := scanner.Text()
		
		if strings.HasPrefix(line, "#") {
			if currentWriter != nil {
				if err := currentWriter.Close(); err != nil {
					utils.PrintError("%v", err)
				}
			}
			
			regionName := strings.TrimSpace(strings.TrimPrefix(line, "#"))
			newWriter, err := writer.NewCSVWriter("locations", regionName)
			if err != nil {
				utils.PrintError("%v", err)
				continue
			}
			currentWriter = newWriter
			utils.PrintInfo("Created new file for region: %s", regionName)
			continue
		}

		if currentWriter == nil {
			continue
		}

		parts := strings.Split(strings.TrimSpace(line), " ")
		if len(parts) < 2 {
			continue
		}

		userCount, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			utils.PrintError("Invalid user count for line: %s", line)
			continue
		}

		location := strings.Join(parts[:len(parts)-1], " ")
		
		userChannel := make(chan *models.User)
		processedChannel := make(chan *models.User)
		var wg sync.WaitGroup

		utils.PrintInfo("Starting to Process %s (Target: %d users)", location, userCount)

		go func(loc string, count int) {
			c.CrawlUsersByLocation(loc, count, userChannel)
		}(location, userCount)

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
				if err := currentWriter.WriteUser(user); err != nil {
					utils.PrintError("%v", err)
				} else {
					utils.PrintUserWritten("%s - %s", user.Username, user.Location)
				}
			}
		}()

		wg.Wait()
	}

	if currentWriter != nil {
		if err := currentWriter.Close(); err != nil {
			utils.PrintError("%v", err)
		}
	}

	utils.PrintSuccess("Completed")
}