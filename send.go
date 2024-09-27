package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func SendEmail(content string) error {
	key := os.Getenv("MAILGUN_TEST_API_KEY")
	endpoint := "https://api.mailgun.net/v3/sandbox1e7a4321500241bc88fbd6fb1ad7d544.mailgun.org/messages"

	From := "bth  <mailgun@sandbox1e7a4321500241bc88fbd6fb1ad7d544.mailgun.org>"
	To := "test@hautogdoad.com"
	Subject := "Foo go test"
	Text := "go test yeah"

	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)
	htmlFw, _ := writer.CreateFormField("html")
	_, err := io.Copy(htmlFw, strings.NewReader(Text))
	if err != nil {
		return err
	}
	fromFw, _ := writer.CreateFormField("from")
	_, err = io.Copy(fromFw, strings.NewReader(From))
	if err != nil {
		return err
	}
	toFw, _ := writer.CreateFormField("to")
	_, err = io.Copy(toFw, strings.NewReader(To))
	if err != nil {
		return err
	}
	subjectFw, _ := writer.CreateFormField("subject")
	_, err = io.Copy(subjectFw, strings.NewReader(Subject))
	if err != nil {
		return err
	}
	writer.Close()

	payload := bytes.NewReader(data.Bytes())
	req, err := http.NewRequest("POST", endpoint, payload)
	if err != nil {
		return err
	}

	req.SetBasicAuth("api", key)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(res)
	fmt.Println(string(body))
	return nil
}
