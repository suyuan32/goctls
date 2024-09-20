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

package vben

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"

	"github.com/suyuan32/goctls/api/parser"
	"github.com/suyuan32/goctls/api/spec"
	"github.com/suyuan32/goctls/util/pathx"
)

var (
	// VarStringOutput describes the output.
	VarStringOutput string
	// VarStringApiFile describes the api file path
	VarStringApiFile string
	// VarStringFolderName describes the folder name to output dir
	VarStringFolderName string
	// VarStringApiPrefix describes the request URL's prefix
	VarStringApiPrefix string
	// VarStringModelName describes the model name
	VarStringModelName string
	// VarStringSubFolder describes the sub folder name
	VarStringSubFolder string
	// VarBoolOverwrite describes whether to overwrite the files, it will overwrite all generated files.
	VarBoolOverwrite bool
	// VarStringFormType describes the form type
	VarStringFormType string
	// VarStringModelChineseName describes the Chinese name of model
	VarStringModelChineseName string
	// VarStringModelEnglishName describes the English name of model
	VarStringModelEnglishName string
)

type GenContext struct {
	ApiDir           string
	ModelDir         string
	ViewDir          string
	Prefix           string
	ModelName        string
	LocaleDir        string
	FolderName       string
	SubFolderName    string
	ApiSpec          *spec.ApiSpec
	UseUUID          bool
	HasStatus        bool
	HasState         bool
	FormType         string
	Overwrite        bool
	ModelChineseName string
	ModelEnglishName string
}

func (g GenContext) Validate() error {
	if g.ApiDir == "" {
		return errors.New("please set the api file path via --api_file")
	} else if !strings.HasSuffix(g.ApiDir, "api") {
		return errors.New("please input correct api file path")
	}
	return nil
}

// GenCRUDLogic is used to generate CRUD file for simple admin backend UI
func GenCRUDLogic(_ *cobra.Command, _ []string) error {
	outputDir, err := filepath.Abs(VarStringOutput)
	if err != nil {
		return err
	}

	apiFile, err := parser.Parse(VarStringApiFile)
	if err != nil {
		return err
	}

	apiOutputDir := filepath.Join(outputDir, "src/api", VarStringFolderName)
	if err := pathx.MkdirIfNotExist(apiOutputDir); err != nil {
		return err
	}
	modelOutputDir := filepath.Join(outputDir, "src/api", VarStringFolderName, "model")
	if err := pathx.MkdirIfNotExist(modelOutputDir); err != nil {
		return err
	}
	viewOutputDir := filepath.Join(outputDir, "src/views", VarStringFolderName)
	if err := pathx.MkdirIfNotExist(viewOutputDir); err != nil {
		return err
	}
	if VarStringSubFolder != "" {
		viewOutputDir = filepath.Join(viewOutputDir, VarStringSubFolder)
		if err := pathx.MkdirIfNotExist(viewOutputDir); err != nil {
			return err
		}
	}
	localeDir := filepath.Join(outputDir, "src/locales/lang")

	var modelName string
	if VarStringModelName != "" {
		modelName = VarStringModelName
	} else {
		modelName = strcase.ToCamel(strings.TrimSuffix(filepath.Base(VarStringApiFile), ".api"))
	}

	genCtx := &GenContext{
		ApiDir:           apiOutputDir,
		ModelDir:         modelOutputDir,
		ViewDir:          viewOutputDir,
		Prefix:           VarStringApiPrefix,
		ModelName:        modelName,
		ApiSpec:          apiFile,
		LocaleDir:        localeDir,
		FolderName:       VarStringFolderName,
		Overwrite:        VarBoolOverwrite,
		FormType:         VarStringFormType,
		ModelChineseName: VarStringModelChineseName,
		ModelEnglishName: VarStringModelEnglishName,
	}

	err = genCtx.Validate()

	if err := genModel(genCtx); err != nil {
		return err
	}

	if err := genApi(genCtx); err != nil {
		return err
	}

	if err := genData(genCtx); err != nil {
		return err
	}

	if err := genLocale(genCtx); err != nil {
		return err
	}

	if err := genIndex(genCtx); err != nil {
		return err
	}

	if err := genDrawer(genCtx); err != nil {
		return err
	}

	if err := genModalIndex(genCtx); err != nil {
		return err
	}

	color.Green.Println("Generate vben files successfully")
	return nil
}
