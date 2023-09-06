package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/chamzzzzzz/supersimplesoup"
)

type Record struct {
	Title string
	Name  string
	Value string
}

func main() {
	urls, err := getURLs()
	if err != nil {
		fmt.Printf("get urls error: %v\n", err)
		return
	}

	var records []*Record
	for _, url := range urls {
		r, err := getRecords(url)
		if err != nil {
			fmt.Printf("get records error: %v\n", err)
			return
		}
		records = append(records, r...)
	}

	b, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		fmt.Printf("marshal error: %v\n", err)
		return
	}

	err = os.WriteFile("records.json", b, 0644)
	if err != nil {
		fmt.Printf("write record file error: %v\n", err)
		return
	}
	fmt.Printf("write record file success\n")
}

func getURLs() ([]string, error) {
	client := &http.Client{}
	resp, err := client.Get("https://data.rmtc.org.cn/gis/listtype0M.html")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	dom, err := supersimplesoup.Parse(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	var urls []string
	for _, a := range dom.QueryAll("li", "class", "datali").Query("a") {
		urls = append(urls, "https://data.rmtc.org.cn/gis/"+a.Href())
	}
	return urls, nil
}

func getRecords(url string) ([]*Record, error) {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	dom, err := supersimplesoup.Parse(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	div, err := dom.Find("div", "class", "pagetitle")
	if err != nil {
		return nil, err
	}
	title := div.Text()

	var records []*Record
	for _, li := range dom.QueryAll("li", "class", "datali") {
		div, err := li.Find("div", "class", "divname")
		if err != nil {
			return nil, err
		}
		span, err := li.Find("span", "class", "label")
		if err != nil {
			return nil, err
		}
		name := strings.TrimSpace(div.Text())
		value := strings.TrimSpace(span.Text())
		record := &Record{
			Title: title,
			Name:  name,
			Value: value,
		}
		records = append(records, record)
	}
	return records, nil
}
