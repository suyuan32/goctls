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

type Create{{.modelName}}Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreate{{.modelName}}Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Create{{.modelName}}Logic {
	return &Create{{.modelName}}Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Create{{.modelName}}Logic) Create{{.modelName}}(req *types.{{.modelName}}Info) (resp *types.BaseMsgResp, err error) {
{{if .optionalService}}	if l.svcCtx.{{.rpcName}}Rpc == nil {
		return nil, errorx.NewCodeUnavailableError({{if .useI18n}}i18n.ServiceUnavailable{{else}}errormsg.ServiceUnavailable{{end}})
	}
{{end}}	data, err := l.svcCtx.{{.rpcName}}Rpc.Create{{.modelName}}(l.ctx,
		&{{.rpcPbPackageName}}.{{.modelName}}Info{ {{.setLogic}}
		})
	if err != nil {
		return nil, err
	}
	return &types.BaseMsgResp{Msg: {{if .useI18n}}l.svcCtx.Trans.Trans(l.ctx, data.Msg){{else}}data.Msg{{end}}}, nil
}
