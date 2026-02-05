package pod

import (
	"fmt"
	"net/http"
	"os"

	"lib.go.io/shplib/ctr/sock"
)

// Container detail
type Container struct {
	Id           string
	Names        string
	Status       string
	RestartCount int64
}

// Pod detail
type PodDetail struct {
	Cgroup     string
	Containers []Container
	Created    string
	Id         string
	InfraId    string
	Name       string
	Namespace  string
	Networks   []string
	Status     string
	Labels     map[string]string
}

func All() {

	body, err := sock.Curl[[]PodDetail]("/libpod/pods/json", http.MethodGet)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, i := range *body {

		fmt.Println(i.Name, i.Id)
		fmt.Println(i.Labels)
		fmt.Println(i.Networks)
	}
}
