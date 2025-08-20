package {{.PkgName}}

import (
	"net/http"

    "github.com/zeromicro/go-zero/rest/httpx"

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
		}{{else}}
		client := make(chan {{.ResponseType}}, 16)
        defer func() {
            close(client)
        }()

        l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)

        threading.GoSafeCtx(r.Context(), func() {
            err := l.{{.Call}}({{if .HasRequest}}&req, {{end}}client)
            if err != nil {
                logc.Errorw(r.Context(), "{{.HandlerName}}", logc.Field("error", err))
                return
            }
        })

        for {
            select {
            case data := <-client:
                output, err := json.Marshal(data)
                if err != nil {
                    logc.Errorw(r.Context(), "{{.HandlerName}}", logc.Field("error", err))
                    continue
                }

                if _, err := fmt.Fprintf(w, "data: %s\n\n", string(output)); err != nil {
                    logc.Errorw(r.Context(), "{{.HandlerName}}", logc.Field("error", err))
                    return
                }
               if flusher, ok := w.(http.Flusher); ok {
                   flusher.Flush()
               }
            case <-r.Context().Done():
                return
            }
        }{{end}}
	}
}
