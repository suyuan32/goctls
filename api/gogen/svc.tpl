package svc

import (
	{{.configImport}}
	{{if .useI18n}}
	"github.com/suyuan32/simple-admin-common/i18n"{{end}}{{if .useCoreRpc}}
	"github.com/suyuan32/simple-admin-core/rpc/coreclient"{{end}}{{if .useEnt}}
	"{{.projectPackage}}/ent"
	_ "{{.projectPackage}}/ent/runtime"
	"github.com/zeromicro/go-zero/core/logx"{{end}}
    {{if .useCasbin}}
    "github.com/zeromicro/go-zero/rest"
    "github.com/casbin/casbin/v2"{{end}}{{if .useCoreRpc}}
	"github.com/zeromicro/go-zero/zrpc"{{end}}
)

type ServiceContext struct {
	Config {{.config}}{{if .hasMiddleware}}
	{{.middleware}}{{end}}{{if .useCasbin}}
	Casbin    *casbin.Enforcer
	Authority rest.Middleware{{end}}{{if .useEnt}}
	DB         *ent.Client{{end}}{{if .useI18n}}
	Trans     *i18n.Translator{{end}}{{if .useCoreRpc}}
	CoreRpc   coreclient.Core{{end}}
}

func NewServiceContext(c {{.config}}) *ServiceContext {
{{if .useCasbin}}
    rds := c.RedisConf.MustNewUniversalRedis()

    cbn := c.CasbinConf.MustNewCasbinWithOriginalRedisWatcher(c.CasbinDatabaseConf.Type, c.CasbinDatabaseConf.GetDSN(), c.RedisConf)
{{end}}
{{if .useI18n}}
    trans := i18n.NewTranslator(c.I18nConf, i18n2.LocaleFS)
{{end}}
{{if .useEnt}}
	entOpts := []ent.Option{
		ent.Log(logx.Info),
		ent.Driver(c.DatabaseConf.NewNoCacheDriver()),
	}

	if c.DatabaseConf.Debug {
		entOpts = append(entOpts, ent.Debug())
	}

	db := ent.NewClient(entOpts...)
{{end}}
	return &ServiceContext{
		Config: c,{{if .useCasbin}}
		Authority: middleware.NewAuthorityMiddleware(cbn, rds).Handle,{{end}}{{if .useI18n}}
		Trans:     trans,{{end}}{{if .useEnt}}
		DB:     db,{{end}}{{if .useCoreRpc}}
		CoreRpc:   coreclient.NewCore(zrpc.NewClientIfEnable(c.CoreRpc)),{{end}}
	}
}
