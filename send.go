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

func SendEmail(content string, subject string) error {
	key := os.Getenv("MAILGUN_TEST_API_KEY")
	if key == "" {
		panic("no api key found")
	}
	endpoint := "https://api.mailgun.net/v3/sandbox1e7a4321500241bc88fbd6fb1ad7d544.mailgun.org/messages"

	fields := make(map[string]string)
	fields["from"] = "bth  <mailgun@sandbox1e7a4321500241bc88fbd6fb1ad7d544.mailgun.org>"
	fields["to"] = "test@hautogdoad.com"
	fields["subject"] = subject
	fields["text"] = content

	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)

	for field, val := range fields {
		formfield, _ := writer.CreateFormField(field)
		_, err := io.Copy(formfield, strings.NewReader(val))
		if err != nil {
			return err
		}
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
