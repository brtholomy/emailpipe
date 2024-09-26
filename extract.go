package main

import (
	"encoding/xml"
	"errors"
	"net/url"
	"os"
	"strings"
)

type Enclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	Url     string   `xml:"url,attr"`
}

type Post struct {
	XMLName     xml.Name  `xml:"item"`
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Enclosure   Enclosure `xml:"enclosure"`
	Content     string    `xml:"encoded"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Items   []*Post  `xml:"channel>item"`
}

func GetSlug(rss *RSS, slug string) (*Post, error) {
	for _, i := range rss.Items {
		u, _ := url.Parse(i.Link)
		if strings.Contains(u.Path, slug) {
			return i, nil
		}
	}
	return nil, errors.New("didn't find the slug")
}

func ExtractPost(source string, slug string) (*Post, error) {
	dat, _ := os.ReadFile(source)
	var rss RSS
	if err := xml.Unmarshal(dat, &rss); err != nil {
		return nil, err
	}
	i, err := GetSlug(&rss, slug)
	if err != nil {
		return nil, err
	}
	return i, nil
}
