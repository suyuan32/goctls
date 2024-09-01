package generator

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/suyuan32/goctls/rpc/execx"
	"github.com/suyuan32/goctls/util/format"
)

const initEntCodeTpl string = `    if err := l.svcCtx.DB.Schema.Create(l.ctx, schema.WithForeignKeys(false)); err != nil {
        logx.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
        return nil, errorx.NewInternalError(err.Error())
    }`

func (g *Generator) GenEntInitCode(zctx *ZRpcContext, abs string) error {
	// gen init code
	name, err := format.FileNamingFormat(g.cfg.NamingFormat, zctx.RpcName)
	if err != nil {
		return err
	}

	_, err = execx.Run(fmt.Sprintf("goctls rpc protoc ./%s.proto --go_out=./types --go-grpc_out=./types --zrpc_out=. --style=%s", name, g.cfg.NamingFormat), abs)
	if err != nil {
		return err
	}

	if zctx.Ent {
		initFileName, err := format.FileNamingFormat(g.cfg.NamingFormat, "initDatabaseLogic")
		if err != nil {
			return err
		}

		initFilePath := filepath.Join(zctx.Output, "internal/logic/base/"+initFileName+".go")

		if !fileutil.IsExist(initFilePath) {
			return fmt.Errorf("init database file not found: %s", initFilePath)
		}

		baseLogicStr, err := fileutil.ReadFileToString(initFilePath)
		if err != nil {
			return errors.Join(err, errors.New("failed to load init database file"))
		}

		baseLogicStr = strings.ReplaceAll(baseLogicStr, "\t// todo: add your logic here and delete this line", initEntCodeTpl)

		if zctx.I18n {
			baseLogicStr = strings.ReplaceAll(baseLogicStr, "BaseResp{},", "BaseResp{Msg: i18n.Success},")
			baseLogicStr = strings.ReplaceAll(baseLogicStr, "\"context\"", "\"context\"\n\t\"entgo.io/ent/dialect/sql/schema\"\n\t\"github.com/suyuan32/simple-admin-common/i18n\"\n\t\"github.com/suyuan32/simple-admin-common/msg/logmsg\"\n\t\"github.com/zeromicro/go-zero/core/errorx\"")
		} else {
			baseLogicStr = strings.ReplaceAll(baseLogicStr, "BaseResp{},", "BaseResp{Msg: errormsg.Success},")
			baseLogicStr = strings.ReplaceAll(baseLogicStr, "\"context\"", "\"context\"\n\t\"entgo.io/ent/dialect/sql/schema\"\n\t\"github.com/suyuan32/simple-admin-common/msg/errormsg\"\n\t\"github.com/suyuan32/simple-admin-common/msg/logmsg\"\n\t\"github.com/zeromicro/go-zero/core/errorx\"")
		}

		err = fileutil.WriteStringToFile(initFilePath, baseLogicStr, false)
		if err != nil {
			return err
		}
	}

	return err
}
