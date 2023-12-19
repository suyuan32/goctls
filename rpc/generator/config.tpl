package config

import (
{{if .isEnt}}   "github.com/suyuan32/simple-admin-common/config"
{{end}}
    "github.com/zeromicro/go-zero/zrpc"

)

type Config struct {
	zrpc.RpcServerConf
{{if .isEnt}}   DatabaseConf config.DatabaseConf
    RedisConf    config.RedisConf
{{end}}
}

