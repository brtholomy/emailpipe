package main

import (
	"flag"
	"fmt"
)

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
	var source = flag.String("source", "https://www.bartholomy.ooo/posts/index.xml", "path to source .xml file")
	var template_path = flag.String("template_path", "./templates/template.html", "template path")
	var slug = flag.String("slug", "", "slug of the post to publish")
	var status = flag.String("status", STATUS_DRAFT, "status can be 'draft' or 'about_to_send'")
	var method = flag.String("method", HTTP_POST, "HTTP method for initial call. Sending a draft requires PATCH")
	var email_id = flag.String("email_id", "", "id of draft email previously created")
	var test_address = flag.String("test_address", "", "email address to receive the draft, defaults to Secrets.Test_address")
	var prod = flag.Bool("prod", false, "whether to send to real prod account")
	flag.Parse()

	secrets, err := GetSecrets(*prod, *test_address)
	if err != nil {
		panic(err)
	}
	opts := Options{*source, *template_path, *status, *slug, *email_id, "", *method, secrets}

	post, err := ExtractPost(&opts)
	if err != nil {
		panic(err)
	}
	fmt.Println(post.Title)
	fmt.Println(post.Link)

	post, err = Assemble(post, &opts)
	if err != nil {
		panic(err)
	}

	resp, err := SendEmail(post, &opts)
	if err != nil {
		if resp != nil {
			fmt.Println(resp.Status)
		}
		panic(err)
	}
	fmt.Println("final response status:", resp.Status)
}
