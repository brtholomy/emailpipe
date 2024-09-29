package main

import (
	"flag"
	"fmt"
)

var source = flag.String("source", "./sample.xml", "path to source .xml file")
var slug = flag.String("slug", "why-fractals", "slug of the post to publish")
var status = flag.String("status", "draft", "status can be 'draft' or 'about_to_send'")
var email_id = flag.String("email_id", "", "id of draft email previously created")

func main() {
	flag.Parse()

	i, err := ExtractPost(*source, *slug)
	if err != nil {
		panic(err)
	}

	fmt.Println(i.Title)
	fmt.Println(i.Description)

	body, err := SendEmail(&i.Content, &i.Title, status, email_id)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}
