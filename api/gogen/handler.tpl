package {{.PkgName}}

import (
	"net/http"

{{if not .UseSSE}}	"github.com/zeromicro/go-zero/rest/httpx"{{end}}

	{{.ImportPackages}}
)

{{.HandlerDoc}}

func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{if .HasRequest}}var req types.{{.RequestType}}
		if err := httpx.Parse(r, &req, {{if .UseValidator}}true{{else}}false{{end}}); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		{{end}}{{if not .UseSSE}}l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
		if err != nil { {{if .TransErr}}
		    err = svcCtx.Trans.TransError(r.Context(), err){{end}}
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			{{if .HasResp}}httpx.OkJsonCtx(r.Context(), w, resp){{else}}httpx.Ok(w){{end}}
		}{{else}}l := {{.LogicName}}.New{{.LogicType}}(r, w, svcCtx)
		l.{{.Call}}({{if .HasRequest}}&req{{end}}){{end}}
	}
}
