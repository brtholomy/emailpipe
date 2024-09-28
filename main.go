package main

import (
	"flag"
	"fmt"
)

var source = flag.String("source", "./sample.xml", "path to source .xml file")
var slug = flag.String("slug", "why-fractals", "slug of the post to publish")
var status = flag.String("status", "draft", "status can be 'draft' or 'about_to_send'")

func main() {
	flag.Parse()

	i, err := ExtractPost(*source, *slug)
	if err != nil {
		panic(err)
	}

	fmt.Println(i.Title)
	fmt.Println(i.Description)

	body, err := SendEmail(&i.Content, &i.Title, status)
	if err != nil {
		panic(err)
	}
	// TODO: loop every time like this:
	// 1. send draft, save id
	// 2. pause and ask for confirmation to send to all
	// 3. use the return id to send
	fmt.Println(string(body))
}
