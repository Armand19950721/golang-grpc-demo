package httpHelper

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type IKeyValuePair struct {
	Key   string
	Value string
}

func SendPostRequest(_url string, header map[string][]string, formData map[string][]string) (string, error) {
	defaultStringReturn := ""
	client := &http.Client{}
	form := url.Values{}
	form = formData

	req, err := http.NewRequest("POST", _url, strings.NewReader(form.Encode()))
	if err != nil {
		return defaultStringReturn, err
	}

	req.Header = header
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)

	if err != nil {
		return defaultStringReturn, err
	}
	defer res.Body.Close()

	cnt, err := io.ReadAll(res.Body)
	if err != nil {
		return defaultStringReturn, err
	}

	fmt.Println("SendPostRequest res")
	fmt.Println(string(cnt))

	return string(cnt), nil
}
