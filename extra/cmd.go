package extra

import (
	"github.com/suyuan32/goctls/extra/ent/importschema"
	"github.com/suyuan32/goctls/extra/ent/localmixin"
	"github.com/suyuan32/goctls/extra/ent/schema"
	"github.com/suyuan32/goctls/extra/ent/template"
	"github.com/suyuan32/goctls/extra/i18n"
	"github.com/suyuan32/goctls/extra/initlogic"
	"github.com/suyuan32/goctls/extra/logviewer"
	"github.com/suyuan32/goctls/extra/makefile"
	"github.com/suyuan32/goctls/extra/proto2api"
	"github.com/suyuan32/goctls/internal/cobrax"
)

var (
	ExtraCmd = cobrax.NewCommand("extra")

	i18nCmd = cobrax.NewCommand("i18n", cobrax.WithRunE(i18n.Gen))

	initCmd = cobrax.NewCommand("init_code", cobrax.WithRunE(initlogic.Gen))

	entCmd = cobrax.NewCommand("ent")

	templateCmd = cobrax.NewCommand("template", cobrax.WithRunE(template.GenTemplate))

	mixinCmd = cobrax.NewCommand("mixin", cobrax.WithRunE(localmixin.GenLocalMixin))

	entImportCmd = cobrax.NewCommand("import", cobrax.WithRunE(importschema.Gen))

	entSchemaCmd = cobrax.NewCommand("schema", cobrax.WithRunE(schema.GenSchema))

	proto2apiCmd = cobrax.NewCommand("proto2api", cobrax.WithRunE(proto2api.Gen))

	makefileCmd = cobrax.NewCommand("makefile", cobrax.WithRunE(makefile.Gen))

	logViewerCmd = cobrax.NewCommand("view_log", cobrax.WithRunE(logviewer.Gen))
)

func init() {
	var (
		i18nCmdFlags      = i18nCmd.Flags()
		initCmdFlags      = initCmd.Flags()
		templateCmdFlags  = templateCmd.Flags()
		mixinCmdFlags     = mixinCmd.Flags()
		makefileCmdFlags  = makefileCmd.Flags()
		entImportCmdFlags = entImportCmd.Flags()
		proto2apiCmdFlags = proto2apiCmd.Flags()
		logViewerCmdFlags = logViewerCmd.Flags()
		entSchemaCmdFlags = entSchemaCmd.Flags()
	)

	i18nCmdFlags.StringVarP(&i18n.VarStringTarget, "target", "t")
	i18nCmdFlags.StringVarP(&i18n.VarStringModelName, "model_name", "m")
	i18nCmdFlags.StringVarP(&i18n.VarStringModelNameZh, "model_name_zh", "z")
	i18nCmdFlags.StringVarP(&i18n.VarStringOutputDir, "output", "o")

	initCmdFlags.StringVarP(&initlogic.VarStringTarget, "target", "t")
	initCmdFlags.StringVarP(&initlogic.VarStringModelName, "model_name", "m")
	initCmdFlags.StringVarP(&initlogic.VarStringOutputPath, "output", "o")

	templateCmdFlags.StringVarP(&template.VarStringDir, "dir", "d")
	templateCmdFlags.StringVarP(&template.VarStringAdd, "add", "a")
	templateCmdFlags.BoolVarP(&template.VarBoolList, "list", "l")
	templateCmdFlags.BoolVarP(&template.VarBoolUpdate, "update", "u")

	makefileCmdFlags.StringVarP(&makefile.VarStringServiceName, "service_name", "n")
	makefileCmdFlags.StringVarP(&makefile.VarStringStyle, "style", "s")
	makefileCmdFlags.StringVarP(&makefile.VarStringDir, "dir", "d")
	makefileCmdFlags.StringVarP(&makefile.VarStringServiceType, "service_type", "t")
	makefileCmdFlags.BoolVarP(&makefile.VarBoolI18n, "i18n", "i")
	makefileCmdFlags.BoolVarP(&makefile.VarBoolEnt, "ent", "e")

	proto2apiCmdFlags.StringVarP(&proto2api.VarStringApiPath, "api_path", "a")
	proto2apiCmdFlags.StringVarP(&proto2api.VarStringProtoPath, "proto_path", "p")
	proto2apiCmdFlags.StringVarP(&proto2api.VarStringModelName, "model_name", "m")
	proto2apiCmdFlags.StringVarP(&proto2api.VarStringGroupName, "group_name", "g")
	proto2apiCmdFlags.BoolVarWithDefaultValue(&proto2api.VarBoolMultiple, "multiple", false)
	proto2apiCmdFlags.StringVarPWithDefaultValue(&proto2api.VarStringJsonStyle, "json_style", "j", "goZero")

	mixinCmdFlags.StringVarP(&localmixin.VarStringDir, "dir", "d")
	mixinCmdFlags.StringVarP(&localmixin.VarStringAdd, "add", "a")
	mixinCmdFlags.BoolVarP(&localmixin.VarBoolList, "list", "l")
	mixinCmdFlags.BoolVarP(&localmixin.VarBoolUpdate, "update", "u")

	entImportCmdFlags.StringVarP(&importschema.VarStringDSN, "dsn", "d")
	entImportCmdFlags.StringVarP(&importschema.VarStringTables, "tables", "t")
	entImportCmdFlags.StringVarP(&importschema.VarStringOutputDir, "output", "o")
	//entImportCmdFlags.BoolVarP(&importschema.VarBoolAutoMixin, "auto_mixin", "a")
	entImportCmdFlags.StringVarP(&importschema.VarStringExcludeTables, "exclude_tables", "e")

	entSchemaCmdFlags.StringVarP(&schema.VarStringModelName, "model_name", "m")

	logViewerCmdFlags.StringVarP(&logviewer.VarStringFilePath, "path", "p")
	logViewerCmdFlags.StringVarP(&logviewer.VarStringWorkspaceSetting, "workspace_setting", "k")
	logViewerCmdFlags.StringVarP(&logviewer.VarStringWorkspace, "workspace", "w")
	logViewerCmdFlags.StringVarP(&logviewer.VarStringLogType, "log_type", "t")
	logViewerCmdFlags.BoolVarP(&logviewer.VarBoolResetWorkspace, "reset_workspace", "r")
	logViewerCmdFlags.BoolVarP(&logviewer.VarBoolList, "list", "l")
	logViewerCmdFlags.IntVarPWithDefaultValue(&logviewer.VarIntMessageCapacity, "size", "s", 10)
	logViewerCmdFlags.StringVarP(&logviewer.VarStringRemoveConfig, "delete_config", "d")

	ExtraCmd.AddCommand(i18nCmd)
	ExtraCmd.AddCommand(initCmd)
	entCmd.AddCommand(templateCmd)
	entCmd.AddCommand(mixinCmd)
	entCmd.AddCommand(entImportCmd)
	entCmd.AddCommand(entSchemaCmd)
	ExtraCmd.AddCommand(entCmd)
	ExtraCmd.AddCommand(makefileCmd)
	ExtraCmd.AddCommand(proto2apiCmd)
	ExtraCmd.AddCommand(logViewerCmd)
}
