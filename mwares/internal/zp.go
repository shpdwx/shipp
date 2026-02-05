package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

var (
	api   = "https://open.bigmodel.cn/api/paas/v4/chat/completions"
	token = "0a12f4c0dfde49ab99247bdae5703708.mjvWWdFTHYzRhrep"
)

type ChatReq struct{}

func Chat(cr ChatReq) {

	header := map[string]string{
		"Bearer": token,
	}

	r := NewReqNet(api)
	r.header(CommonHeader, header)

	b, err := json.Marshal(cr)
	if err != nil {
		log.Fatal("request struct format error")
	}

	v, err := r.run(string(b), get)
	if err != nil {
		log.Fatalf("request failed: %v", err)
	}

	resp, ok := v.(chatResp)
	if !ok {
		log.Fatal("parse response failed")
	}

	fmt.Println(resp)
}

type chatResp struct{}

func get(body io.ReadCloser) interface{} {
	return chatResp{}
}
