package httphelper

import (
	"bytes"
	"context"
	"github.com/ProninIgorr/alcomarket-nearst/pkg/helpers/utils"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	Timeout time.Duration
}

func (c *Client) MakeRequest(ctx context.Context, method string, url string, body []byte, headers map[string]string) ([]byte, error) {
	var jsonStr = body
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("X-Request-ID", utils.GetRequestId(ctx))
	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)

	return resBody, err
}
