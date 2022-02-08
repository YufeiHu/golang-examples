package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type resultStruct struct {
	url        string
	readLength int
	err        error
}

func getPage(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	return len(body), nil
}

func main() {
	urls := []string{
		"http://www.google.com/",
		"http://www.yahoo.com",
		"http://www.bing.com",
		"http://bbc.co.uk",
	}

	resultChan := make(chan resultStruct)
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			readLength, err := getPage(url)
			resultChan <- resultStruct{url: url, readLength: readLength, err: err}
		}(url)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		if result.err != nil {
			fmt.Printf("failed to fetch %q: %v\n", result.url, result.err)
			continue
		}
        fmt.Printf("%s has a response of length: %d\n", result.url, result.readLength)
	}
}
