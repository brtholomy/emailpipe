package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const SECRET_SOURCE string = "./SECRETS.json"
const BASEURL string = "https://api.buttondown.email/v1/emails"
const ENDPOINT_SEND_DRAFT string = "send-draft"
const STATUS_DRAFT string = "draft"
const STATUS_FINAL string = "about_to_send"
const HTTP_POST string = "POST"
const HTTP_PATCH string = "PATCH"

type Secrets struct {
	Test_buttondown_api_key string `json:test_buttondown_api_key`
	Prod_buttondown_api_key string `json:prod_buttondown_api_key`
	Key                     string `json:key`
	Test_address            string `json:test_address`
	Test_subscriber         string `json:test_subscriber`
}

type EmailPayload struct {
	// pointers to allow them to be "optional": nil value will json.Marshal
	Subject     *string  `json:"subject"`
	Body        *string  `json:"body"`
	Status      *string  `json:"status"`
	Recipients  []string `json:"recipients"`
	Subscribers []string `json:"subscribers"`
}

type ResponsePayload struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

// Fill out Secrets struct from .gitignore'd SECRET_SOURCE.
func GetSecrets(prod bool, test_address string) (*Secrets, error) {
	dat, err := os.ReadFile(SECRET_SOURCE)
	if err != nil {
		return nil, err
	}
	var secrets Secrets
	if err := json.Unmarshal(dat, &secrets); err != nil {
		return nil, err
	}
	// NOTE: .Key is left initially empty by design.
	secrets.Key = secrets.Test_buttondown_api_key
	if prod {
		secrets.Key = secrets.Prod_buttondown_api_key
	}
	if test_address != "" {
		secrets.Test_address = test_address
	}
	return &secrets, nil
}

// Send the Post using the Options.
//
// Proceeds like this:
//
// 1. create and send draft, save opts.Email_id.
// 2. pause and ask for confirmation to send to all.
// 3. use the response email_id to send.
//
// Also possible to feed in the opts.Email_id directly, if opts.Status=="about_to_send"
func SendEmail(post *Post, opts *Options) ([]byte, error) {
	payload := EmailPayload{&post.Title, &post.Content, &opts.Status, nil, nil}

	// skip directly to prod if draft aleady exists:
	if opts.Status == STATUS_FINAL {
		if opts.Email_id == "" {
			return nil, errors.New("sending to prod requires an email_id of a draft")
		}
		var err error
		opts.Endpoint, err = url.JoinPath(BASEURL, opts.Email_id)
		if err != nil {
			return nil, err
		}
		return SendPayload(payload, opts)
	}

	// create draft
	opts.Endpoint = BASEURL
	res, err := SendPayload(payload, opts)
	if err != nil {
		return nil, err
	}
	var resp ResponsePayload
	if err := json.Unmarshal(res, &resp); err != nil {
		return nil, err
	}
	if resp.Status != STATUS_DRAFT {
		return nil, errors.New("return status is not draft, cancelling send")
	}

	// send draft using return id
	opts.Email_id = resp.Id
	fmt.Println("email_id:", opts.Email_id)
	opts.Endpoint, err = url.JoinPath(BASEURL, opts.Email_id, ENDPOINT_SEND_DRAFT)
	if err != nil {
		return nil, err
	}
	payload.Recipients = []string{opts.Secrets.Test_address}
	// payload.Subscribers = []string{opts.Secrets.Test_subscriber}
	res, err = SendPayload(payload, opts)
	if err != nil {
		return nil, err
	}

	// continue to prod
	fmt.Println("sent draft. send to all subscribers? Yes/n:")
	var answer string
	if _, err := fmt.Scanln(&answer); err != nil {
		return nil, err
	}
	if answer == "Yes" {
		// NOTE: difference is no /send-draft at the end, and
		// status="about_to_send", and "PATCH" method
		opts.Endpoint, err = url.JoinPath(BASEURL, opts.Email_id)
		if err != nil {
			return nil, err
		}
		opts.Method = HTTP_PATCH
		opts.Status = STATUS_FINAL
		payload.Status = &opts.Status
		return SendPayload(payload, opts)
	}
	fmt.Println("quitting")
	return res, nil
}

func SendPayload(payload EmailPayload, opts *Options) ([]byte, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(b)

	req, err := http.NewRequest(opts.Method, opts.Endpoint, r)
	if err != nil {
		return nil, err
	}
	// NOTE: buttondown expects "Token key", standard seems to be "Bearer key"
	// NOTE: mailgun expects BasicAuth, with "api" as username.
	req.Header.Add("Authorization", "Token "+opts.Secrets.Key)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
