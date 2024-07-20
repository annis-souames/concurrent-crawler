package main

import (
	"annis/webcrawler/crawler"
	"fmt"
	"sync"
)

const maxLinks = 200

func crawl(link string, urlList chan<- string, resultList chan<- string, wg *sync.WaitGroup, visitedUrls *sync.Map, counter *int, counterLock *sync.Mutex) {
	defer wg.Done()

	basePage, err := crawler.FetchUrl(link)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	resultList <- basePage.Title

	visitedUrls.Store(link, true)
	for _, url := range basePage.UrlList {
		counterLock.Lock()
		if *counter < maxLinks {
			*counter++
			urlList <- url
		}
		counterLock.Unlock()
	}
}

func main() {
	startUrl := "http://golang.org"
	basePage, err := crawler.FetchUrl(startUrl)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	var wg sync.WaitGroup
	urlChan := make(chan string)
	resultChan := make(chan string)
	visitedUrls := &sync.Map{}
	counter := 0
	var counterLock sync.Mutex

	go func() {
		for url := range urlChan {
			if _, loaded := visitedUrls.LoadOrStore(url, true); !loaded {
				wg.Add(1)
				go crawl(url, urlChan, resultChan, &wg, visitedUrls, &counter, &counterLock)
			}
		}
	}()

	for _, u := range basePage.UrlList {
		fullUrl, err := crawler.ResolveURL(u, startUrl)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		if _, loaded := visitedUrls.LoadOrStore(fullUrl, true); !loaded {
			counterLock.Lock()
			if counter < maxLinks {
				counter++
				wg.Add(1)
				go crawl(fullUrl, urlChan, resultChan, &wg, visitedUrls, &counter, &counterLock)
			}
			counterLock.Unlock()
		}
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for r := range resultChan {
		fmt.Println("Title: ", r)
	}

	close(urlChan)
}