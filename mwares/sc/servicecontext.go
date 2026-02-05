package sc

import (
	"context"

	"github.com/shpdwx/mwares/conf"
	"github.com/shpdwx/mwares/internal"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Ctx    context.Context
	Cfg    *conf.Config
	Logger *zap.SugaredLogger
}

func NewServiceContext(c conf.Config) *ServiceContext {
	return &ServiceContext{
		Cfg:    &c,
		Logger: internal.InitZap(),
	}
}
