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

// foo@hautogdoad.com : specific to test account:
var test_subscriber string = "b7c993fc-1baf-4ba7-81f5-e5cb096dcfa0"
var tbartholomy string = "t@bartholomy.ooo"
var baseurl string = "https://api.buttondown.email/v1/emails"
var final_status string = "about_to_send"

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

// Loops like this:
// 1. send draft, save id
// 2. pause and ask for confirmation to send to all
// 3. use the return id to send
// also possible to feed in the email_id directly, if status="about_to_send"
func SendEmail(content *string, subject *string, status *string, email_id *string) ([]byte, error) {
	payload := EmailPayload{subject, content, status, nil, nil}
	var endpoint string

	// skip directly to prod if draft aleady exists:
	if *status == "about_to_send" {
		if *email_id == "" {
			return nil, errors.New("sending to prod requires an email_id of a draft")
		}
		endpoint, _ = url.JoinPath(baseurl, *email_id)
		return SendPayload(payload, endpoint)
	}

	// create draft
	endpoint = baseurl
	res, err := SendPayload(payload, endpoint)
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

	// send draft
	fmt.Println("email_id:", resp.Id)
	endpoint, _ = url.JoinPath(baseurl, resp.Id, "send-draft")
	payload.Recipients = []string{tbartholomy}
	// payload.Subscribers = []string{test_subscriber}
	res, err = SendPayload(payload, endpoint)
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
		endpoint, _ = url.JoinPath(baseurl, resp.Id)
		payload = EmailPayload{nil, nil, &final_status, nil, nil}
		return SendPayload(payload, endpoint, "PATCH")
	}
	fmt.Println("quitting")
	return res, nil
}

func SendPayload(payload EmailPayload, endpoint string, methods ...string) ([]byte, error) {
	key := os.Getenv("BUTTONDOWN_TEST_API_KEY")
	if key == "" {
		return nil, errors.New("no api key found")
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(b)

	// NOTE: effective default param
	method := "POST"
	if len(methods) > 0 {
		method = methods[0]
	}
	req, err := http.NewRequest(method, endpoint, r)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Token "+key)
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
