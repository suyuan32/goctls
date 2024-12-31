package {{.packageName}}

import (
	"context"{{if .hasTime}}
	"time"{{end}}

	"{{.projectPath}}{{.importPrefix}}/ent/{{.modelNameLowerCase}}"
	"{{.projectPath}}{{.importPrefix}}/ent/predicate"
	"{{.projectPath}}{{.importPrefix}}/internal/svc"
	"{{.projectPath}}{{.importPrefix}}/internal/utils/dberrorhandler"
	"{{.projectPath}}{{.importPrefix}}/types/{{.projectName}}"

{{if .hasUUID}}    "github.com/suyuan32/simple-admin-common/utils/uuidx"
{{end}}{{if .HasCreated}}	"github.com/suyuan32/simple-admin-common/utils/pointy"{{end}}
    "github.com/zeromicro/go-zero/core/logx"
)

type Get{{.modelName}}ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGet{{.modelName}}ListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Get{{.modelName}}ListLogic {
	return &Get{{.modelName}}ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *Get{{.modelName}}ListLogic) Get{{.modelName}}List(in *{{.projectName}}.{{.modelName}}ListReq) (*{{.projectName}}.{{.modelName}}ListResp, error) {
{{.predicateData}}

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	resp := &{{.projectName}}.{{.modelName}}ListResp{}
	resp.Total = result.PageDetails.Total

	for _, v := range result.List {
		resp.Data = append(resp.Data, &{{.projectName}}.{{.modelName}}Info{
			Id:          {{if .useUUID}}pointy.GetPointer(v.ID.String()){{else}}&v.ID{{end}},{{if .HasCreated}}
			CreatedAt:   pointy.GetPointer(v.CreatedAt.UnixMilli()),
			UpdatedAt:   pointy.GetPointer(v.UpdatedAt.UnixMilli()),{{end}}
{{.listData}}
		})
	}

	return resp, nil
}
