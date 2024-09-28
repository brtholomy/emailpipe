package main

import (
	"flag"
	"fmt"
	"os"
)

var source = flag.String("source", "./sample.xml", "path to source .xml file")
var css = flag.String("css", "./custom.css", "path to .css file")
var slug = flag.String("slug", "why-fractals", "slug of the post to publish")
var draft = flag.Bool("draft", true, "whether to send a draft")

func main() {
	flag.Parse()

	i, err := ExtractPost(*source, *slug)
	if err != nil {
		panic(err)
	}

	dat, err := os.ReadFile(*css)
	if err != nil {
		panic(err)
	}
	cssstr := string(dat)

	fmt.Println(i.Title)
	fmt.Println(i.Description)

	var status string
	if *draft {
		status = "draft"
	} else {
		status = "sent"
	}
	if err := SendEmail(i.Content, cssstr, i.Title, status); err != nil {
		panic(err)
	}
}
