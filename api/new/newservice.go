// Copyright (C) 2023  Ryan SU (https://github.com/suyuan32)

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package new

import (
	_ "embed"
	"errors"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/duke-git/lancet/v2/fileutil"

	"github.com/spf13/cobra"

	"github.com/iancoleman/strcase"

	"github.com/suyuan32/goctls/api/gogen"
	conf "github.com/suyuan32/goctls/config"
	"github.com/suyuan32/goctls/util"
	"github.com/suyuan32/goctls/util/pathx"
)

var (
	// VarStringHome describes the goctl home.
	VarStringHome string
	// VarStringRemote describes the remote git repository.
	VarStringRemote string
	// VarStringBranch describes the git branch.
	VarStringBranch string
	// VarStringStyle describes the style of output files.
	VarStringStyle string
	// VarBoolErrorTranslate describes whether to translate error
	VarBoolErrorTranslate bool
	// VarBoolUseCasbin describe whether to use Casbin
	VarBoolUseCasbin bool
	// VarBoolUseI18n describe whether to use i18n
	VarBoolUseI18n bool
	// VarModuleName describe the module name
	VarModuleName string
	// VarIntServicePort describe the service port exposed
	VarIntServicePort int
	// VarBoolEnt describes whether to use ent in api
	VarBoolEnt bool
	// VarBoolCore describes whether to generate core rpc code
	VarBoolCore bool
)

// CreateServiceCommand fast create service
func CreateServiceCommand(_ *cobra.Command, args []string) error {
	dirName := args[0]
	if len(VarStringStyle) == 0 {
		VarStringStyle = conf.DefaultFormat
	}
	if strings.Contains(dirName, "-") {
		return errors.New("api new command service name not support strikethrough, because this will used by function name")
	}

	abs, err := filepath.Abs(dirName)
	if err != nil {
		return err
	}

	err = pathx.MkdirIfNotExist(abs)
	if err != nil {
		return err
	}

	err = pathx.MkdirIfNotExist(filepath.Join(abs, "desc"))
	if err != nil {
		return err
	}

	if len(VarStringRemote) > 0 {
		repo, _ := util.CloneIntoGitHome(VarStringRemote, VarStringBranch)
		if len(repo) > 0 {
			VarStringHome = repo
		}
	}

	if len(VarStringHome) > 0 {
		pathx.RegisterGoctlHome(VarStringHome)
	}

	apiFilePath := filepath.Join(abs, "desc", "all.api")

	text, err := pathx.LoadTemplate(category, apiTemplateFile, baseApiTmpl)
	if err != nil {
		return err
	}

	baseApiFile, err := os.Create(filepath.Join(abs, "desc", "base.api"))
	if err != nil {
		return err
	}
	defer baseApiFile.Close()

	t := template.Must(template.New("baseApiTemplate").Parse(text))
	if err := t.Execute(baseApiFile, map[string]string{
		"name": strcase.ToCamel(dirName),
	}); err != nil {
		return err
	}

	allApiFile, err := os.Create(filepath.Join(abs, "desc", "all.api"))
	if err != nil {
		return err
	}
	defer allApiFile.Close()

	allTpl := template.Must(template.New("allApiTemplate").Parse(allApiTmpl))
	if err := allTpl.Execute(allApiFile, map[string]string{
		"name": strcase.ToCamel(dirName),
	}); err != nil {
		return err
	}

	var moduleName string

	if VarModuleName != "" {
		moduleName = VarModuleName
	} else {
		moduleName = dirName
	}

	genCtx := &gogen.GenContext{
		UseCasbin:     VarBoolUseCasbin,
		UseI18n:       VarBoolUseI18n,
		TransErr:      VarBoolErrorTranslate,
		ModuleName:    moduleName,
		Port:          VarIntServicePort,
		UseMakefile:   true,
		UseDockerfile: true,
		UseEnt:        VarBoolEnt,
		IsNewProject:  true,
		UseCoreRpc:    VarBoolCore,
	}

	err = gogen.DoGenProject(apiFilePath, abs, VarStringStyle, genCtx)

	err = fileutil.WriteStringToFile(filepath.Join(abs, ".gitignore"), GitIgnoreTmpl, false)
	if err != nil {
		return err
	}

	return err
}
