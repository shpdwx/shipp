package internal

import (
	"fmt"
	"os"
	"time"
)

var (
	prefix = "mwares-Ce4air"
)

type HandleFile func(string) interface{}

func TempFile(content string, handle HandleFile) (interface{}, error) {

	if len(content) == 0 {
		return nil, nil
	}

	t := time.Now().UnixMicro()
	tmpf, err := os.CreateTemp("", fmt.Sprintf("%s-%d", prefix, t))
	if err != nil {
		return nil, err
	}

	defer func() {
		tmpf.Close()
		os.Remove(tmpf.Name())
	}()

	if _, err := tmpf.WriteString(content); err != nil {
		return nil, err
	}

	return handle(tmpf.Name()), nil
}
