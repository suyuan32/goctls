package generator

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"

	conf "github.com/suyuan32/goctls/config"
	"github.com/suyuan32/goctls/rpc/parser"
	"github.com/suyuan32/goctls/util"
	"github.com/suyuan32/goctls/util/format"
	"github.com/suyuan32/goctls/util/pathx"
)

//go:embed svc.tpl
var svcTemplate string

// GenSvc generates the servicecontext.go file, which is the resource dependency of a service,
// such as rpc dependency, model dependency, etc.
func (g *Generator) GenSvc(ctx DirContext, _ parser.Proto, cfg *conf.Config, c *ZRpcContext) error {
	dir := ctx.GetSvc()
	svcFilename, err := format.FileNamingFormat(cfg.NamingFormat, "service_context")
	if err != nil {
		return err
	}

	fileName := filepath.Join(dir.Filename, svcFilename+".go")
	text, err := pathx.LoadTemplate(category, svcTemplateFile, svcTemplate)
	if err != nil {
		return err
	}

	imports := strings.Builder{}
	imports.WriteString(fmt.Sprintf("\t\"%v\"\n", ctx.GetConfig().Package))
	if c.Ent {
		imports.WriteString(fmt.Sprintf("\t\"%s/ent\"\n\n", ctx.GetMain().Package))
		imports.WriteString("\t\"github.com/zeromicro/go-zero/core/logx\"\n\t\"github.com/zeromicro/go-zero/core/stores/redis\"\n")
	}

	return util.With("svc").GoFmt(true).Parse(text).SaveTo(map[string]any{
		"imports": imports.String(),
		"isEnt":   c.Ent,
	}, fileName, false)
}
