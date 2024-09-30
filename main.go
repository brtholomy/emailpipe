package main

import (
	"flag"
	"fmt"
)

var source = flag.String("source", "./sample.xml", "path to source .xml file")
var template_path = flag.String("template_path", "./template.html", "template path")
var slug = flag.String("slug", "why-fractals", "slug of the post to publish")
var status = flag.String("status", "draft", "status can be 'draft' or 'about_to_send'")
var email_id = flag.String("email_id", "", "id of draft email previously created")

func main() {
	flag.Parse()

	post, err := ExtractPost(*source, *slug)
	if err != nil {
		panic(err)
	}
	fmt.Println(post.Title)
	fmt.Println(post.Link)

	content, err := Assemble(post, template_path)
	if err != nil {
		panic(err)
	}

	body, err := SendEmail(content, &post.Title, status, email_id)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}
