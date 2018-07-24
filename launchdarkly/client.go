package launchdarkly

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Client struct {
	AccessToken string
}

func (c *Client) GetStatus(url string) (int, error) {
	status, _, err := c.execute("GET", url, nil, []int{})
	return status, err
}

func (c *Client) Get(url string, expectedStatus []int) (interface{}, error) {
	_, response, err := c.execute("GET", url, nil, expectedStatus)

	var parsedResponse interface{}
	json.Unmarshal(response, &parsedResponse)

	return parsedResponse, err
}

func (c *Client) GetInto(url string, expectedStatus []int, target interface{}) error {
	_, response, err := c.execute("GET", url, nil, expectedStatus)

	json.Unmarshal(response, target)

	return err
}

func (c *Client) Post(url string, body interface{}, expectedStatus []int) ([]byte, error) {
	_, response, err := c.execute("POST", url, body, expectedStatus)
	return response, err
}

func (c *Client) Patch(url string, body interface{}, expectedStatus []int) ([]byte, error) {
	_, response, err := c.execute("PATCH", url, body, expectedStatus)
	return response, err
}

func (c *Client) Delete(url string, expectedStatus []int) error {
	_, _, err := c.execute("DELETE", url, nil, expectedStatus)
	return err
}

func (c *Client) execute(method string, url string, body interface{}, expectedStatus []int) (int, []byte, error) {
	requestBody, err := json.Marshal(body)
	if err != nil {
		return 0, nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("Authorization", c.AccessToken)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}

	println(method + " " + url + " returned HTTP status " + strconv.Itoa(resp.StatusCode))

	if len(expectedStatus) > 0 {
		found := false
		for _, status := range expectedStatus {
			if status == resp.StatusCode {
				found = true
				break
			}
		}

		if !found {
			return resp.StatusCode, nil, errors.New(method + " " + url + " did not return one of the expected HTTP status codes. Got HTTP " + strconv.Itoa(resp.StatusCode) + "\n" + string(responseBody))
		}
	}

	return resp.StatusCode, responseBody, nil
}
