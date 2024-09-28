package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type EmailPayload struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Status  string `json:"status"`
}

func SendEmail(content string, css string, subject string, status string) error {
	endpoint := "https://api.buttondown.email/v1/emails"
	key := os.Getenv("BUTTONDOWN_TEST_API_KEY")
	if key == "" {
		panic("no api key found")
	}
	// FIXME: buttondown doesn't like css here. may be impossible to inject via
	// the API?
	content = css + content + "</div></body></html>"

	payload := EmailPayload{subject, content, status}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	r := bytes.NewReader(b)

	req, err := http.NewRequest("POST", endpoint, r)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+key)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}
