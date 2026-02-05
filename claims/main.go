package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/shpdwx/claims/auth"
)

var (
	ctx = context.Background()
)

func main() {

	aut := auth.NewJwtToken(ctx, "user-center")

	refresh(aut)
}

func create(aut auth.JwtToken) {

	aut.User(10023, "dalls")
	aut.Device("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36")

	result, _ := aut.Gen()
	b, _ := json.Marshal(result)
	fmt.Println(string(b))
}

func verify(aut auth.JwtToken) {
	str := `
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEwMDIzLCJ1biI6ImRhbGxzIiwianRpIjoiWldOak5UWXdOek10TmpVeVpTMDBNVEF4TFRsa09Ea3RPVGxsT1RnNVpEYzFNVEEyIiwiaXNzIjoidXNlci1jZW50ZXIiLCJzdWIiOiIxMDAyMyIsImV4cCI6MTc2ODk4MjYzMSwiaWF0IjoxNzY4OTgxNzMxfQ.6rDNISXZs7F69d4PxyoFAWmmWanl-lSXG78iaiRkFmM
`
	aut.Validate(str)
}

func refresh(aut auth.JwtToken) {

	str := `ZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SjFhV1FpT2pFd01ESXpMQ0pxZEdraU9pSmFWMDVxVGxSWmQwNTZUWFJPYWxWNVdsTXdNRTFVUVhoTVZHeHJUMFJyZEU5VWJHeFBWR2MxV2tSak1VMVVRVElpTENKa2FXUWlPaUpVVnprMllWZDRjMWxUT0RGTWFrRm5TMFV4YUZreWJIVmtSemw2WVVSeloxTlhOVEJhVjNkblZGZEdha2xGT1ZSSlJtZG5UVlJDWmsxVVZtWk9lV3RuVVZoQ2QySkhWbGhhVjBwTVlWaFJkazVVVFROTWFrMHlTVU5vVEZOR1VrNVVRM2RuWWtkc2NscFRRa2hhVjA1eVlubHJaMUV5YUhsaU1qRnNUSHBGTUUxNU5IZE1ha0YxVFVOQ1ZGbFhXbWhqYld0MlRsUk5NMHhxVFRJaUxDSnBjM01pT2lKMWMyVnlMV05sYm5SbGNpSXNJbk4xWWlJNkltUmhiR3h6SWl3aVpYaHdJam94Tnpjd01Ua3hNek14TENKcFlYUWlPakUzTmpnNU9ERTNNekY5LjhjalhTbDI1NFRwa01NU2phc0E0aWhjZHZIU0QzYzNuUGdBZUJHbmJycWs=`

	fmt.Println(aut.Refresh(str))

}
