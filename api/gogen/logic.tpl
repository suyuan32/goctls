package {{.pkgName}}

import (
	{{.imports}}
)

type {{.logic}} struct {
	logx.Logger{{if .useSSE}}
	request *http.Request
	response http.ResponseWriter{{end}}
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func New{{.logic}}({{if not .useSSE}}ctx context.Context, {{else}}r *http.Request, w http.ResponseWriter, {{end}}svcCtx *svc.ServiceContext) *{{.logic}} {
	return &{{.logic}}{
		Logger: logx.WithContext({{if not .useSSE}}ctx{{else}}r.Context(){{end}}),
		ctx:    {{if not .useSSE}}ctx{{else}}r.Context(){{end}},{{if .useSSE}}
		request : r,
		response : w,{{end}}
		svcCtx: svcCtx,
	}
}

func (l *{{.logic}}) {{.function}}({{.request}}) {{.responseType}} {
	// todo: add your logic here and delete this line

	{{.returnString}}
}
