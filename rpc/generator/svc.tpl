package svc

import (
{{.imports}}
)

type ServiceContext struct {
	Config config.Config
    {{if .isEnt}}DB     *ent.Client
    Redis  redis.UniversalClient
{{end}}

}

func NewServiceContext(c config.Config) *ServiceContext {
{{if .isEnt}}   entOpts := []ent.Option{
		ent.Log(logx.Info),
		ent.Driver(c.DatabaseConf.NewNoCacheDriver()),
	}

	if c.DatabaseConf.Debug {
		entOpts = append(entOpts, ent.Debug())
	}

	db := ent.NewClient(entOpts...)
    
    {{end}}
	return &ServiceContext{
		Config: c,
		{{if .isEnt}}DB:     db,
		Redis:  c.RedisConf.MustNewUniversalRedis(),{{end}}
	}
}
