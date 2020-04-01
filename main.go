package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

const k = 5

func main() {
	ch := make(chan struct{}, k)

	urlCh := make(chan string, 5)

	urls := `https://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org
				\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org`

	mUrls := strings.Split(urls, `\n`)

	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer wg.Done()
		for url := range urlCh {
			ch <- struct{}{}

			wg.Add(1)

			go func() {
				defer wg.Done()

				n, err := makeRequest(url)
				if err != nil {
					fmt.Println(err)
				}

				fmt.Println("count :", n)

				<-ch
			}()
		}
	}()

	for _, value := range mUrls {
			urlCh <- value
	}

	close(urlCh)

	wg.Wait()
}

func makeRequest(url string) (int, error) {
	client := http.Client{}

	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	var buf bytes.Buffer

	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return 0, err
	}

	return strings.Count(buf.String(), "Go"), nil
}
