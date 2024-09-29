package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendRequestHTTP(
	url string,
	method string,
	headers map[string]string,
	body map[string]string,
) {

	bodyString, err := json.Marshal(body)

	if err != nil {
		fmt.Println(err)
	}

	// Send request HTTP
	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyString))
	if err != nil {
		fmt.Println(err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

}
