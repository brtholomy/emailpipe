package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type EmailForm struct {
	From    string
	To      string
	Subject string
	Text    string
}

func main() {
	key := os.Getenv("MAILGUN_TEST_API_KEY")
	endpoint := "https://api.mailgun.net/v3/sandbox1e7a4321500241bc88fbd6fb1ad7d544.mailgun.org/messages"

	var e EmailForm
	e.From = "bth  <mailgun@sandbox1e7a4321500241bc88fbd6fb1ad7d544.mailgun.org>"
	e.To = "test@hautogdoad.com"
	e.Subject = "Foo go test"
	e.Text = "go test yeah"

	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	j := bytes.NewReader(b)

	req, _ := http.NewRequest("POST", endpoint, j)
	req.Header.Add("Authorization", "Bearer "+key)
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}
