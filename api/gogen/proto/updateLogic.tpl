package {{.modelNameLowerCase}}

import (
	"context"

	"{{.projectPackage}}{{.importPrefix}}/internal/svc"
	"{{.projectPackage}}{{.importPrefix}}/internal/types"
	"{{.rpcPackage}}"
{{if .useI18n}}
	"github.com/suyuan32/simple-admin-common/i18n"{{end}}{{if .optionalService}}{{if not .useI18n}}
	"github.com/suyuan32/simple-admin-common/msg/errormsg"{{end}}
	"github.com/zeromicro/go-zero/core/errorx"{{end}}
	"github.com/zeromicro/go-zero/core/logx"
)

type Update{{.modelName}}Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdate{{.modelName}}Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Update{{.modelName}}Logic {
	return &Update{{.modelName}}Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Update{{.modelName}}Logic) Update{{.modelName}}(req *types.{{.modelName}}Info) (resp *types.BaseMsgResp, err error) {
{{if .optionalService}}	if !l.svcCtx.Config.{{.rpcName}}Rpc.Enabled {
		return nil, errorx.NewCodeUnavailableError({{if .useI18n}}i18n.ServiceUnavailable{{else}}errormsg.ServiceUnavailable{{end}})
	}
{{end}}	data, err := l.svcCtx.{{.rpcName}}Rpc.Update{{.modelName}}(l.ctx,
		&{{.rpcPbPackageName}}.{{.modelName}}Info{
			Id:          req.Id,{{.setLogic}}
		})
	if err != nil {
		return nil, err
	}
	return &types.BaseMsgResp{Msg: {{if .useI18n}}l.svcCtx.Trans.Trans(l.ctx, data.Msg){{else}}data.Msg{{end}}}, nil
}
