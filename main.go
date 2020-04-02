package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

const k = 5

func main() {
	ch := make(chan struct{}, k)

	urlCh := make(chan string, 5)

	reader := bufio.NewReader(os.Stdin)
	urls, _ := reader.ReadString('\n')

	re := regexp.MustCompile(`[[:space:]]`)

	urls = re.ReplaceAllString(urls, "")

	//urls := `https://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org
	//			\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org`

	mUrls := strings.Split(urls, `\n`)

	fmt.Println(mUrls)

	wg := sync.WaitGroup{}

	wg.Add(1)

	var count int

	go func() {
		defer wg.Done()
		for url := range urlCh {
			ch <- struct{}{}

			wg.Add(1)

			go func(url string) {
				defer wg.Done()

				n, err := makeRequest(url)
				if err != nil {
					fmt.Println(err)
				}

				count += n

				fmt.Printf("count %s: %d \n", url, n)

				<-ch
			}(url)
		}
	}()

	for _, value := range mUrls {
		urlCh <- value
	}

	close(urlCh)

	wg.Wait()

	close(ch)

	println("total: ", count)
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
