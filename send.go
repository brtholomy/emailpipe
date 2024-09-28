package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// foo@hautogdoad.com:
var test_subscriber string = "b7c993fc-1baf-4ba7-81f5-e5cb096dcfa0"
var tbartholomy string = "t@bartholomy.ooo"

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

func SendEmail(content *string, subject *string, status *string) ([]byte, error) {
	// TODO: fill in recepients every time
	payload := EmailPayload{subject, content, status, nil, nil}
	endpoint := "https://api.buttondown.email/v1/emails"
	if *status == "about_to_send" {
		return SendPayload(payload, endpoint)
	}

	// create draft
	res, err := SendPayload(payload, endpoint)
	if err != nil {
		return nil, err
	}
	var resp ResponsePayload
	if err := json.Unmarshal(res, &resp); err != nil {
		return nil, err
	}
	if resp.Status != "draft" {
		return res, errors.New("return status is not draft, cancelling send")
	}

	// send draft
	fmt.Println("email id: ", resp.Id)
	endpoint = fmt.Sprintf("https://api.buttondown.email/v1/emails/%s/send-draft", resp.Id)
	payload.Recipients = []string{tbartholomy}
	payload.Subscribers = []string{test_subscriber}
	return SendPayload(payload, endpoint)
}

func SendPayload(payload EmailPayload, endpoint string) ([]byte, error) {
	key := os.Getenv("BUTTONDOWN_TEST_API_KEY")
	if key == "" {
		return nil, errors.New("no api key found")
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(b)

	req, err := http.NewRequest("POST", endpoint, r)
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
