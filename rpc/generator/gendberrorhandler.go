package generator

import (
	_ "embed"
	"path/filepath"

	"github.com/iancoleman/strcase"

	conf "github.com/suyuan32/goctls/config"
	"github.com/suyuan32/goctls/rpc/parser"
	"github.com/suyuan32/goctls/util"
	"github.com/suyuan32/goctls/util/format"
	"github.com/suyuan32/goctls/util/pathx"
)

//go:embed dberrorhandler.tpl
var dbErrorHandlerTemplateText string

func (g *Generator) GenErrorHandler(ctx DirContext, _ parser.Proto, cfg *conf.Config, c *ZRpcContext) error {
	filename, err := format.FileNamingFormat(cfg.NamingFormat, "error_handler.go")
	if err != nil {
		return err
	}

	handlerFilename := filepath.Join(ctx.GetInternal().Filename, "utils/dberrorhandler", filename)
	if err := pathx.MkdirIfNotExist(filepath.Join(ctx.GetInternal().Filename, "utils")); err != nil {
		return err
	}

	if err := pathx.MkdirIfNotExist(filepath.Join(ctx.GetInternal().Filename, "utils", "dberrorhandler")); err != nil {
		return err
	}

	err = util.With("dbErrorHandler").Parse(dbErrorHandlerTemplateText).SaveTo(map[string]any{
		"package":     ctx.GetMain().Package,
		"serviceName": strcase.ToCamel(c.RpcName),
		"useI18n":     c.I18n,
	}, handlerFilename, false)
	return err
}
