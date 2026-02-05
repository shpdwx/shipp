package sock

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func sockPath(s string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return home + "/" + s
}

// Curl sock 请求
func Curl[T any](url, method string) (*T, error) {

	path := sockPath(".local/share/containers/podman/machine/podman.sock")

	// sock
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial("unix", path)
		}}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	// handle api url
	api := "http://d/v5.0.0/" + strings.TrimLeft(url, "/")

	// handle request
	req, err := http.NewRequest(strings.ToUpper(method), api, nil)
	if err != nil {
		return nil, err
	}

	// handle response
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// body to []byte
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// to struct
	result := new(T) // create not nil pointer
	if err = json.Unmarshal(body, result); err != nil {
		return nil, err
	}
	return result, nil
}
