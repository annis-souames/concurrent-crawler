package main

import (
	"fmt"
)

func main() {
	page := crawler.FetchUrl("https://www.google.com")
	fmt.Println(string(page))
}
