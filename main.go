package main

import (
	"flag"
	"fmt"
)

var source = flag.String("source", "./sample.xml", "path to source .xml file")
var slug = flag.String("slug", "why-fractals", "slug of the post to publish")

func main() {
	flag.Parse()

	i, err := ExtractPost(*source, *slug)
	if err != nil {
		panic(err)
	}

	fmt.Println(i.Title)
	fmt.Println(i.Description)

	if err := SendEmail(i.Content, i.Title); err != nil {
		panic(err)
	}
}
