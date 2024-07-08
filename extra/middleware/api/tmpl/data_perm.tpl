package middleware

import (
	"errors"
	"github.com/redis/go-redis/v9"{{if .useTrans}}
    "github.com/suyuan32/simple-admin-common/i18n"{{end}}
	"github.com/suyuan32/simple-admin-common/orm/ent/entctx/datapermctx"
	"github.com/suyuan32/simple-admin-common/orm/ent/entctx/deptctx"
	"github.com/suyuan32/simple-admin-common/orm/ent/entctx/rolectx"
	"github.com/suyuan32/simple-admin-common/orm/ent/entenum"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

type DataPermMiddleware struct {
	Rds   redis.UniversalClient{{if .useTrans}}
    Trans *i18n.Translator{{end}}
}

func NewDataPermMiddleware(rds redis.UniversalClient{{if .useTrans}}, trans *i18n.Translator{{end}}) *DataPermMiddleware {
	return &DataPermMiddleware{
		Rds:   rds,{{if .useTrans}}
        Trans: trans,{{end}}
	}
}

func (m *DataPermMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var subDept, dataScope, customDept string

		deptId, err := deptctx.GetDepartmentIDFromCtx(ctx)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		roleCodes, err := rolectx.GetRoleIDFromCtx(ctx)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		dataScope, err = m.Rds.Get(ctx, datapermctx.GetRoleScopeDataPermRedisKey(roleCodes)).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				httpx.Error(w, errorx.NewApiUnauthorizedError("please log in again"))
				return
			} else {
				logx.Error("redis error", logx.Field("detail", err))
				httpx.Error(w, errorx.NewInternalError({{if .useTrans}}m.Trans.Trans(ctx, i18n.RedisError){{else}}"Redis Error"{{end}}))
				return
			}
		}

		ctx = datapermctx.WithScopeContext(ctx, dataScope)

		if dataScope == entenum.DataPermOwnDeptAndSubStr {
			subDept, err = m.Rds.Get(ctx, datapermctx.GetSubDeptDataPermRedisKey(deptId)).Result()
			if err != nil {
				if errors.Is(err, redis.Nil) {
					httpx.Error(w, errorx.NewApiUnauthorizedError("please log in again"))
					return
				} else {
					logx.Error("redis error", logx.Field("detail", err))
                    httpx.Error(w, errorx.NewInternalError({{if .useTrans}}m.Trans.Trans(ctx, i18n.RedisError){{else}}"Redis error"{{end}}))
					return
				}
			}

			ctx = datapermctx.WithSubDeptContext(ctx, subDept)
		}

		if dataScope == entenum.DataPermCustomDeptStr {
			customDept, err = m.Rds.Get(ctx, datapermctx.GetRoleCustomDeptDataPermRedisKey(roleCodes)).Result()
			if err != nil {
				if errors.Is(err, redis.Nil) {
					httpx.Error(w, errorx.NewApiUnauthorizedError("please log in again"))
					return
				} else {
					logx.Error("redis error", logx.Field("detail", err))
                    httpx.Error(w, errorx.NewInternalError({{if .useTrans}}m.Trans.Trans(ctx, i18n.RedisError){{else}}"Redis error"{{end}}))
					return
				}
			}

			ctx = datapermctx.WithCustomDeptContext(ctx, customDept)
		}

		next(w, r.WithContext(ctx))
	}
}
