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
const FINAL_STATUS string = "about_to_send"

type Secrets struct {
	Test_buttondown_api_key string `json:test_buttondown_api_key`
	Prod_buttondown_api_key string `json:prod_buttondown_api_key`
	Key                     string `json:key`
	Test_email              string `json:test_email`
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

func GetSecrets(prod bool) (*Secrets, error) {
	dat, err := os.ReadFile(SECRET_SOURCE)
	if err != nil {
		return nil, err
	}
	var secrets Secrets
	if err := json.Unmarshal(dat, &secrets); err != nil {
		return nil, err
	}
	secrets.Key = secrets.Test_buttondown_api_key
	if prod {
		secrets.Key = secrets.Prod_buttondown_api_key
	}
	return &secrets, nil
}

// Loops like this:
// 1. send draft, save id
// 2. pause and ask for confirmation to send to all
// 3. use the return id to send
// also possible to feed in the email_id directly, if status="about_to_send"
func SendEmail(content *string, subject *string, opts *Options) ([]byte, error) {
	payload := EmailPayload{subject, content, &opts.Status, nil, nil}

	// skip directly to prod if draft aleady exists:
	if opts.Status == "about_to_send" {
		if opts.Email_id == "" {
			return nil, errors.New("sending to prod requires an email_id of a draft")
		}
		opts.Endpoint, _ = url.JoinPath(BASEURL, opts.Email_id)
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
	if resp.Status != "draft" {
		return nil, errors.New("return status is not draft, cancelling send")
	}

	// send draft using return id
	opts.Email_id = resp.Id
	fmt.Println("email_id:", opts.Email_id)
	opts.Endpoint, _ = url.JoinPath(BASEURL, opts.Email_id, "send-draft")
	payload.Recipients = []string{opts.Secrets.Test_email}
	// payload.Subscribers = []string{opts.Secrets.Test_subscriber}
	res, err = SendPayload(payload, opts)
	if err != nil {
		return nil, err
	}

	// continue to prod
	fmt.Println("sent draft. continue to prod? Y/n:")
	var answer string
	_, err = fmt.Scanln(&answer)
	if err != nil {
		return nil, err
	}
	if answer == "Y" {
		// NOTE: difference is no /send-draft at the end, and
		// status="about_to_send", and "PATCH" method
		opts.Endpoint, _ = url.JoinPath(BASEURL, opts.Email_id)
		opts.Method = "PATCH"
		opts.Status = FINAL_STATUS
		payload = EmailPayload{nil, nil, &opts.Status, nil, nil}
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
