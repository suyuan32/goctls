package gogen

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/suyuan32/goctls/config"
	"github.com/suyuan32/goctls/rpc/execx"
	"github.com/suyuan32/goctls/util/format"
)

const initEntCodeTpl string = `    if err := l.svcCtx.DB.Schema.Create(l.ctx, schema.WithForeignKeys(false)); err != nil {
        logx.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
        return nil, errorx.NewInternalError(err.Error())
    }`

func GenEntInitCode(dir string, cfg *config.Config, g *GenContext) error {
	// gen init code
	_, err := execx.Run(fmt.Sprintf("goctls api go --api ./desc/all.api --dir ./ --trans_err=true --style=%s", cfg.NamingFormat), dir)
	if err != nil {
		return err
	}

	if g.UseEnt {
		initFileName, err := format.FileNamingFormat(cfg.NamingFormat, "initDatabaseLogic")
		if err != nil {
			return err
		}

		initFilePath := filepath.Join(dir, "internal/logic/base/"+initFileName+".go")

		if !fileutil.IsExist(initFilePath) {
			return fmt.Errorf("init database file not found: %s", initFilePath)
		}

		baseLogicStr, err := fileutil.ReadFileToString(initFilePath)
		if err != nil {
			return errors.Join(err, errors.New("failed to load init database file"))
		}

		baseLogicStr = strings.ReplaceAll(baseLogicStr, "\t// todo: add your logic here and delete this line", initEntCodeTpl)

		if g.UseI18n {
			baseLogicStr = strings.ReplaceAll(baseLogicStr, "return\n", "return &types.BaseMsgResp{Msg:  l.svcCtx.Trans.Trans(l.ctx, i18n.Success)},nil\n")
			baseLogicStr = strings.ReplaceAll(baseLogicStr, "\"context\"", "\"context\"\n\t\"entgo.io/ent/dialect/sql/schema\"\n\t\"github.com/suyuan32/simple-admin-common/i18n\"\n\t\"github.com/suyuan32/simple-admin-common/msg/logmsg\"\n\t\"github.com/zeromicro/go-zero/core/errorx\"")
		} else {
			baseLogicStr = strings.ReplaceAll(baseLogicStr, "return\n", "return &types.BaseMsgResp{Msg: errormsg.Success},nil\n")
			baseLogicStr = strings.ReplaceAll(baseLogicStr, "\"context\"", "\"context\"\n\t\"entgo.io/ent/dialect/sql/schema\"\n\t\"github.com/suyuan32/simple-admin-common/msg/errormsg\"\n\t\"github.com/suyuan32/simple-admin-common/msg/logmsg\"\n\t\"github.com/zeromicro/go-zero/core/errorx\"")
		}

		err = fileutil.WriteStringToFile(initFilePath, baseLogicStr, false)
		if err != nil {
			return err
		}
	}

	return err
}
