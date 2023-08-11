// Copyright 2023 The Ryan SU Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package importschema

import (
	"context"
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/suyuan32/goctls/extra/ent/importschema/mux"
	"path/filepath"
	"strings"
)

var (
	VarStringDSN       string
	VarStringOutputDir string
	VarStringTables    string
	VarBoolAutoMixin   bool
)

type GenContext struct {
	Dsn       string
	OutputDir string
	Tables    []string
	AutoMixin bool
}

func Gen(_ *cobra.Command, _ []string) (err error) {
	if err = validateData(); err != nil {
		return err
	}

	var outputDir string

	if VarStringOutputDir == "" {
		outputDir, err = filepath.Abs(".")
		if err != nil {
			return err
		}

		if fileutil.IsExist(outputDir + "/ent/schema/") {
			outputDir = outputDir + "/ent/schema/"
		}
	} else if !strings.HasSuffix(VarStringOutputDir, "/") {
		VarStringTables = VarStringOutputDir + "/"

		outputDir, err = filepath.Abs(VarStringOutputDir)
		if err != nil {
			return err
		}
	}

	ctx := &GenContext{}
	ctx.OutputDir = outputDir
	ctx.Dsn = VarStringDSN
	ctx.AutoMixin = VarBoolAutoMixin
	ctx.Tables = strings.Split(VarStringTables, ",")

	drv, err := mux.Default.OpenImport(ctx.Dsn)
	if err != nil {
		return fmt.Errorf("failed to create import driver - %v", err)
	}
	i, err := NewImport(
		WithTables(ctx.Tables),
		WithDriver(drv),
	)
	if err != nil {
		return fmt.Errorf("create importer failed: %v", err)
	}

	mutations, err := i.SchemaMutations(context.Background())
	if err != nil {
		return fmt.Errorf("schema import failed - %v", err)
	}

	if err = WriteSchema(mutations, WithSchemaPath(ctx.OutputDir)); err != nil {
		return fmt.Errorf("schema writing failed - %v", err)
	}

	if err = FormatFile(ctx); err != nil {
		return fmt.Errorf("format file failed - %v", err)
	}

	color.Green.Println("Generating schemas successfully")

	return err
}

func validateData() error {
	if VarStringDSN == "" {
		return errors.New("the dsn cannot be empty, please set it by \"-d\"")
	} else if !strings.HasPrefix(VarStringDSN, "mysql") && !strings.HasPrefix(VarStringDSN, "postgres") {
		return errors.New("the database is not supported")
	}
	return nil
}
