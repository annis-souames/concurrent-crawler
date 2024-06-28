package main

import (
	"annis/webcrawler/crawler"
	"fmt"
)

func main() {
	page := crawler.FetchUrl("https://www.google.com")
	fmt.Println(string(page))
}
