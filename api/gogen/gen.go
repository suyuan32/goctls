package gogen

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gookit/color"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/suyuan32/goctls/api/gogen/ent"
	"github.com/suyuan32/goctls/extra/ent/template"
	"github.com/suyuan32/goctls/util/format"

	apiformat "github.com/suyuan32/goctls/api/format"
	"github.com/suyuan32/goctls/api/gogen/proto"
	"github.com/suyuan32/goctls/api/parser"
	apiutil "github.com/suyuan32/goctls/api/util"
	"github.com/suyuan32/goctls/config"
	"github.com/suyuan32/goctls/pkg/golang"
	"github.com/suyuan32/goctls/rpc/execx"
	"github.com/suyuan32/goctls/util"
	"github.com/suyuan32/goctls/util/pathx"
)

const tmpFile = "%s-%d"

var (
	tmpDir = path.Join(os.TempDir(), "goctls")
	// VarStringDir describes the directory.
	VarStringDir string
	// VarStringAPI describes the API.
	VarStringAPI string
	// VarStringHome describes the go home.
	VarStringHome string
	// VarStringRemote describes the remote git repository.
	VarStringRemote string
	// VarStringBranch describes the branch.
	VarStringBranch string
	// VarStringStyle describes the style of output files.
	VarStringStyle string
	// VarBoolErrorTranslate describes whether to translate error
	VarBoolErrorTranslate bool
	// VarBoolUseCasbin describe whether to use Casbin
	VarBoolUseCasbin bool
	// VarBoolUseI18n describe whether to use i18n
	VarBoolUseI18n bool
	// VarStringProto describe the proto file path
	VarStringProto string
	// VarStringAPIServiceName describe the API service name
	VarStringAPIServiceName string
	// VarStringRPCServiceName describe the RPC service name
	VarStringRPCServiceName string
	// VarStringModelName describe which model for generating
	VarStringModelName string
	// VarIntSearchKeyNum describe the number of search keys
	VarIntSearchKeyNum int
	// VarStringOutput describes the output.
	VarStringOutput string
	// VarStringRpcName describes the rpc name in service context
	VarStringRpcName string
	// VarStringGrpcPbPackage describes the grpc package
	VarStringGrpcPbPackage string
	// VarBoolMultiple describes whether the proto contains multiple services
	VarBoolMultiple bool
	// VarStringJSONStyle describes the JSON tag format.
	VarStringJSONStyle string
	// VarBoolOverwrite describes whether to overwrite the files, it will overwrite all generated files.
	VarBoolOverwrite bool
	// VarStringSchema describes the ent schema path
	VarStringSchema string
	// VarStringGroupName describes whether to use group
	VarStringGroupName string
	// VarStringImportPrefix describes the prefix in import
	VarStringImportPrefix string
	// VarStringExtraField describes the extra field for api logic code such as user id
	VarStringExtraField string
	// VarBoolDisabledValidator describes whether to disable validator
	VarBoolDisabledValidator bool
	// VarBoolOptionalService describes whether to generate for optional service
	VarBoolOptionalService bool
)

// GoCommand gen go project files from command line
func GoCommand(_ *cobra.Command, _ []string) error {
	apiFile := VarStringAPI
	dir := VarStringDir
	namingStyle := VarStringStyle
	home := VarStringHome
	remote := VarStringRemote
	branch := VarStringBranch
	genCtx := &GenContext{
		UseCasbin:    VarBoolUseCasbin,
		UseI18n:      VarBoolUseI18n,
		TransErr:     VarBoolErrorTranslate,
		ModuleName:   "",
		ExtraField:   VarStringExtraField,
		UseValidator: !VarBoolDisabledValidator,
	}
	if len(remote) > 0 {
		repo, _ := util.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}

	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}
	if len(apiFile) == 0 {
		return errors.New("missing -api")
	}
	if len(dir) == 0 {
		return errors.New("missing -dir")
	}

	return DoGenProject(apiFile, dir, namingStyle, genCtx)
}

// GenContext describes the data used for api file generation
type GenContext struct {
	UseCasbin     bool
	UseI18n       bool
	TransErr      bool
	ModuleName    string
	Port          int
	UseMakefile   bool
	UseDockerfile bool
	ImportPrefix  string
	UseEnt        bool
	ExtraField    string
	IsNewProject  bool
	UseValidator  bool
}

// DoGenProject gen go project files with api file
func DoGenProject(apiFile, dir, style string, g *GenContext) error {
	color.Green.Println("Generating...")

	api, err := parser.Parse(apiFile)
	if err != nil {
		return err
	}

	if err := api.Validate(); err != nil {
		return err
	}

	cfg, err := config.NewConfig(style)
	if err != nil {
		return err
	}

	logx.Must(pathx.MkdirIfNotExist(dir))
	if g.ModuleName != "" {
		_, err = execx.Run("go mod init "+g.ModuleName, dir)
		if err != nil {
			return err
		}
	}

	rootPkg, err := golang.GetParentPackage(dir)
	if err != nil {
		return err
	}

	logx.Must(genEtc(dir, cfg, api, g))
	logx.Must(genConfig(dir, cfg, api, g))
	logx.Must(genMain(dir, rootPkg, cfg, api, g))
	logx.Must(genServiceContext(dir, rootPkg, cfg, api, g))
	logx.Must(genTypes(dir, cfg, api))
	logx.Must(genRoutes(dir, rootPkg, cfg, api))
	logx.Must(genHandlers(dir, rootPkg, cfg, api, g))
	logx.Must(genLogic(dir, rootPkg, cfg, api, g))
	logx.Must(genMiddleware(dir, cfg, api))

	if g.UseDockerfile {
		service, err := format.FileNamingFormat(cfg.NamingFormat, api.Service.Name)
		if err != nil {
			return err
		}

		_, err = execx.Run(fmt.Sprintf("goctls docker -p %d -s %s -t api -l", g.Port, service), dir)
	}

	if g.UseMakefile {
		var serviceType string
		if g.UseEnt {
			serviceType = "single"
		} else {
			serviceType = "api"
		}

		makefileCmd := fmt.Sprintf("goctls extra makefile -t %s -s %s -n %s", serviceType, style, api.Service.Name)
		if g.UseI18n {
			makefileCmd += " -i"
		}

		if g.UseEnt {
			makefileCmd += " -e"
		}

		_, err = execx.Run(makefileCmd, dir)

		if err != nil {
			return err
		}
	}

	if g.UseCasbin {
		logx.Must(genCasbin(dir, cfg, g))
	}

	if g.UseI18n {
		logx.Must(genI18n(dir))
	}

	if g.UseEnt {
		_, err := execx.Run(fmt.Sprintf("go run -mod=mod entgo.io/ent/cmd/ent new %s",
			strcase.ToCamel(api.Service.Name)), dir)
		if err != nil {
			return err
		}

		_, err = execx.Run("go mod tidy", dir)
		if err != nil {
			return err
		}

		// _, err = execx.Run("go run -mod=mod entgo.io/ent/cmd/ent generate ./ent/schema", dir)
		// if err != nil {
		// 	return err
		// }

		err = pathx.MkdirIfNotExist(filepath.Join(dir, "ent", "template"))
		if err != nil {
			return err
		}

		paginationTplPath := filepath.Join(dir, "ent", "template", "pagination.tmpl")
		notNilTplPath := filepath.Join(dir, "ent", "template", "set_not_nil.tmpl")
		if !pathx.FileExists(paginationTplPath) {
			err = os.WriteFile(paginationTplPath, []byte(template.PaginationTmpl), os.ModePerm)
			if err != nil {
				return err
			}

			err = os.WriteFile(notNilTplPath, []byte(template.NotNilTmpl), os.ModePerm)
			if err != nil {
				return err
			}
		}

		// gen ent error handler
		err = GenErrorHandler(dir, style, rootPkg)
		if err != nil {
			return err
		}

		// gen ent transaction util
		err = GenEntTx(dir, style, rootPkg)
		if err != nil {
			return err
		}

		_, err = execx.Run("go mod tidy", dir)
		if err != nil {
			return err
		}
	}

	if g.IsNewProject {
		_, err = execx.Run("goctls project upgrade", dir)
		if err != nil {
			return err
		}
	}

	if err := backupAndSweep(apiFile); err != nil {
		return err
	}

	if err := apiformat.ApiFormatByPath(apiFile, false); err != nil {
		return err
	}

	color.Green.Println("Done.")
	return nil
}

func backupAndSweep(apiFile string) error {
	var err error
	var wg sync.WaitGroup

	wg.Add(2)
	_ = os.MkdirAll(tmpDir, os.ModePerm)

	go func() {
		_, fileName := filepath.Split(apiFile)
		_, e := apiutil.Copy(apiFile, fmt.Sprintf(path.Join(tmpDir, tmpFile), fileName, time.Now().Unix()))
		if e != nil {
			err = e
		}
		wg.Done()
	}()
	go func() {
		if e := sweep(); e != nil {
			err = e
		}
		wg.Done()
	}()
	wg.Wait()

	return err
}

func sweep() error {
	keepTime := time.Now().AddDate(0, 0, -7)
	return filepath.Walk(tmpDir, func(fpath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		pos := strings.LastIndexByte(info.Name(), '-')
		if pos > 0 {
			timestamp := info.Name()[pos+1:]
			seconds, err := strconv.ParseInt(timestamp, 10, 64)
			if err != nil {
				// print error and ignore
				fmt.Println(color.Red.Sprintf("sweep ignored file: %s", fpath))
				return nil
			}

			tm := time.Unix(seconds, 0)
			if tm.Before(keepTime) {
				if err := os.RemoveAll(fpath); err != nil {
					fmt.Println(color.Red.Sprintf("failed to remove file: %s", fpath))
					return err
				}
			}
		}

		return nil
	})
}

func GenCRUDLogicByProto(_ *cobra.Command, _ []string) error {
	params := &proto.GenLogicByProtoContext{
		ProtoDir:        VarStringProto,
		OutputDir:       VarStringOutput,
		RPCServiceName:  VarStringRPCServiceName,
		APIServiceName:  VarStringAPIServiceName,
		Style:           VarStringStyle,
		ModelName:       VarStringModelName,
		SearchKeyNum:    VarIntSearchKeyNum,
		RpcName:         VarStringRpcName,
		GrpcPackage:     VarStringGrpcPbPackage,
		Multiple:        VarBoolMultiple,
		JSONStyle:       VarStringJSONStyle,
		UseI18n:         VarBoolUseI18n,
		ImportPrefix:    VarStringImportPrefix,
		Overwrite:       VarBoolOverwrite,
		OptionalService: VarBoolOptionalService,
	}

	err := params.Validate()
	if err != nil {
		return err
	}

	err = proto.GenLogicByProto(params)
	if err != nil {
		return err
	}

	return err
}

func GenCRUDLogicByEnt(_ *cobra.Command, _ []string) error {
	params := &ent.GenEntLogicContext{
		Schema:       VarStringSchema,
		Output:       VarStringOutput,
		ServiceName:  VarStringAPIServiceName,
		Style:        VarStringStyle,
		ModelName:    VarStringModelName,
		SearchKeyNum: VarIntSearchKeyNum,
		GroupName:    VarStringGroupName,
		JSONStyle:    VarStringJSONStyle,
		UseI18n:      VarBoolUseI18n,
		ImportPrefix: VarStringImportPrefix,
		Overwrite:    VarBoolOverwrite,
	}

	err := params.Validate()
	if err != nil {
		return err
	}

	err = ent.GenEntLogic(params)
	if err != nil {
		return err
	}

	return err
}
