package main

import (
	"flag"
	"fmt"
)

var source = flag.String("source", "./templates/sample.xml", "path to source .xml file")
var template_path = flag.String("template_path", "./templates/template.html", "template path")
var slug = flag.String("slug", "why-fractals", "slug of the post to publish")
var status = flag.String("status", "draft", "status can be 'draft' or 'about_to_send'")
var email_id = flag.String("email_id", "", "id of draft email previously created")
var prod = flag.Bool("prod", false, "whether to send to real prod account")

type Options struct {
	Source   string
	Template string
	Status   string
	Slug     string
	Email_id string
	Endpoint string
	Method   string
	Secrets  *Secrets
}

func main() {
	flag.Parse()

	secrets, err := GetSecrets(*prod)
	if err != nil {
		panic(err)
	}
	opts := Options{*source, *template_path, *status, *slug, *email_id, "", "POST", secrets}

	post, err := ExtractPost(&opts)
	if err != nil {
		panic(err)
	}
	fmt.Println(post.Title)
	fmt.Println(post.Link)

	content, err := Assemble(post, &opts)
	if err != nil {
		panic(err)
	}

	body, err := SendEmail(content, &post.Title, &opts)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}
