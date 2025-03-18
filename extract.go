package main

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
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

func GetSlug(rss *RSS, opts *Options) (*Post, error) {
	if opts.Slug == "" {
		return nil, errors.New("--slug flag can't be empty")
	}
	for _, i := range rss.Items {
		u, _ := url.Parse(i.Link)
		if filepath.Join("/posts/", opts.Slug) == filepath.Clean(u.Path) {
			return i, nil
		}
	}
	return nil, errors.New("didn't find the slug")
}

// Returns a *Post by matching the opts.Slug against the opts.Source.
func ExtractPost(opts *Options) (*Post, error) {
	resp, err := http.Get(opts.Source)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rss RSS
	if err := xml.Unmarshal(dat, &rss); err != nil {
		return nil, err
	}
	i, err := GetSlug(&rss, opts)
	if err != nil {
		return nil, err
	}
	return i, nil
}
