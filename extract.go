package main

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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
	Items   []Post   `xml:"channel>item"`
}

// NOTE: take the RSS pointer since this is fairly large.
func GetSlug(rss *RSS, opts Options) (Post, error) {
	if opts.Slug == "" {
		return Post{}, errors.New("--slug flag can't be empty")
	}
	for _, p := range rss.Items {
		u, _ := url.Parse(p.Link)
		if filepath.Join("/posts/", opts.Slug) == filepath.Clean(u.Path) {
			return p, nil
		}
	}
	return Post{}, errors.New("didn't find the slug")
}

func isURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func isFilePath(str string) bool {
	if !filepath.IsAbs(str) {
		// make absolute
		absPath, err := filepath.Abs(str)
		if err != nil {
			return false
		}
		str = absPath
	}
	_, err := os.Stat(str)
	return err == nil
}

// gets the []byte content of either a URL or file.
func GetUrlOrFile(opts Options) ([]byte, error) {
	var dat []byte
	var err error
	if isURL(opts.Source) {
		resp, err := http.Get(opts.Source)
		defer resp.Body.Close()
		if err != nil {
			return nil, err
		}
		dat, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	} else if isFilePath(opts.Source) {
		dat, err = os.ReadFile(opts.Source)
		if err != nil {
			return nil, err
		}
	}
	return dat, nil
}

// Returns a Post by matching the opts.Slug against the opts.Source.
func ExtractPost(opts Options) (Post, error) {
	dat, err := GetUrlOrFile(opts)
	if err != nil {
		return Post{}, err
	}
	var rss RSS
	if err := xml.Unmarshal(dat, &rss); err != nil {
		return Post{}, err
	}
	p, err := GetSlug(&rss, opts)
	if err != nil {
		return Post{}, err
	}
	return p, nil
}
