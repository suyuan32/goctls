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

type Get{{.modelName}}ByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGet{{.modelName}}ByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Get{{.modelName}}ByIdLogic {
	return &Get{{.modelName}}ByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Get{{.modelName}}ByIdLogic) Get{{.modelName}}ById(req *types.{{if .useUUID}}UU{{end}}IDReq) (resp *types.{{.modelName}}InfoResp, err error) {
{{if .optionalService}}	if l.svcCtx.{{.rpcName}}Rpc == nil {
		return nil, errorx.NewCodeUnavailableError({{if .useI18n}}i18n.ServiceUnavailable{{else}}errormsg.ServiceUnavailable{{end}})
	}
{{end}}	data, err := l.svcCtx.{{.rpcName}}Rpc.Get{{.modelName}}ById(l.ctx, &{{.rpcPbPackageName}}.{{if .useUUID}}UU{{end}}IDReq{Id: req.Id})
	if err != nil {
		return nil, err
	}

	return &types.{{.modelName}}InfoResp{
		BaseDataInfo: types.BaseDataInfo{
			Code: 0,
			Msg:  {{if .useI18n}}l.svcCtx.Trans.Trans(l.ctx, i18n.Success){{else}}"successful"{{end}},
		},
		Data: types.{{.modelName}}Info{
            Base{{if .useUUID}}UU{{end}}IDInfo: types.Base{{if .useUUID}}UU{{end}}IDInfo{
                Id:        data.Id,
                CreatedAt: data.CreatedAt,
                UpdatedAt: data.UpdatedAt,
            },{{.setLogic}}
		},
	}, nil
}

