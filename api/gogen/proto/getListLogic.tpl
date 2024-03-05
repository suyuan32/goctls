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

type Get{{.modelName}}ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGet{{.modelName}}ListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Get{{.modelName}}ListLogic {
	return &Get{{.modelName}}ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Get{{.modelName}}ListLogic) Get{{.modelName}}List(req *types.{{.modelName}}ListReq) (resp *types.{{.modelName}}ListResp, err error) {
{{if .optionalService}}	if !l.svcCtx.Config.{{.rpcName}}Rpc.Enabled {
		return nil, errorx.NewCodeUnavailableError({{if .useI18n}}i18n.ServiceUnavailable{{else}}errormsg.ServiceUnavailable{{end}})
	}
{{end}}	data, err := l.svcCtx.{{.rpcName}}Rpc.Get{{.modelName}}List(l.ctx,
		&{{.rpcPbPackageName}}.{{.modelName}}ListReq{
			Page:        req.Page,
			PageSize:    req.PageSize,{{.searchKeys}}
		})
	if err != nil {
		return nil, err
	}
	resp = &types.{{.modelName}}ListResp{}
	resp.Msg = {{if .useI18n}}l.svcCtx.Trans.Trans(l.ctx, i18n.Success){{else}}"successful"{{end}}
	resp.Data.Total = data.GetTotal()

	for _, v := range data.Data {
		resp.Data.Data = append(resp.Data.Data,
			types.{{.modelName}}Info{
{{if .HasCreated}}				Base{{if .useUUID}}UU{{end}}ID{{.IdType}}Info: types.Base{{if .useUUID}}UU{{end}}ID{{.IdType}}Info{
					Id:        v.Id,
					CreatedAt: v.CreatedAt,
					UpdatedAt: v.UpdatedAt,
				},{{else}}			Id:  v.Id,{{end}}{{.setLogic}}
			})
	}
	return resp, nil
}
