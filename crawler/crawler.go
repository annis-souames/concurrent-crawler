package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// FetchUrl fetches the content of a URL and returns it as a byte slice.
func FetchUrl(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	}
	return body
}
