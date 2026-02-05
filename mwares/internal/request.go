package internal

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type reqNet struct {
	Method string
	Api    string
	Header map[string]string
}

func NewReqNet(api string) *reqNet {
	return &reqNet{
		Api:    api,
		Method: "POST",
	}
}

func (r *reqNet) setMethod(method string) {
	if method == "" {
		return
	}
	r.Method = strings.ToUpper(method)
}

type AddHeader func(map[string]string) map[string]string

func (r *reqNet) header(ah AddHeader, m map[string]string) {
	ah(m)
	r.Header = m
}

type FetchResp func(io.ReadCloser) interface{}

func (r reqNet) run(data string, resp FetchResp) (interface{}, error) {
	payload := strings.NewReader(data)
	req, err := http.NewRequest(r.Method, r.Api, payload)
	if err != nil {
		return nil, err
	}

	for k, v := range r.Header {
		req.Header.Add(k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return resp(res.Body), nil
}

func CommonHeader(m map[string]string) map[string]string {

	if len(m) == 0 {
		return map[string]string{"Content-Type": "application/json"}
	}
	newMap := make(map[string]string, len(m))

	for k, v := range m {
		switch strings.ToLower(k) {
		case "bearer":
			newMap["Authorization"] = "Bearer " + v
		default:
			newMap[k] = v
		}
	}

	return newMap
}

func do() {

	payload := strings.NewReader("")

	req, _ := http.NewRequest("POST", "", payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
}
