package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Enclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	Url     string   `xml:"url,attr"`
}

type Item struct {
	XMLName     xml.Name  `xml:"item"`
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Enclosure   Enclosure `xml:"enclosure"`
	Content     string    `xml:"encoded"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Items   []*Item  `xml:"channel>item"`
}

func main() {
	dat, _ := os.ReadFile("./sample.xml")
	var rss RSS

	if err := xml.Unmarshal(dat, &rss); err != nil {
		panic(err)
	}
	i := rss.Items[1]
	fmt.Println(i.Title)
	fmt.Println(i.Description)
	fmt.Println(i.Link)
	fmt.Println(i.Enclosure.Url)
	fmt.Println(i.Content)
}
