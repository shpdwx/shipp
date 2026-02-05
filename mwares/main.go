package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/shpdwx/mwares/conf"
	"github.com/shpdwx/mwares/internal"
	"github.com/shpdwx/mwares/sc"
)

var ctx = context.Background()

func main() {

	c := conf.LoadConfig()

	serv := sc.NewServiceContext(c)
	serv.Ctx = ctx

	// sc.Ctx = context.WithValue(sc.Ctx, "request_id", uuid.NewString())

	c.Minio.RootPath = "/" + c.Server.Name
	internal.UploadMinio(serv.Ctx, &c.Minio, uuid.NewString())
}
