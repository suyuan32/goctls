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

package importschema

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/gookit/color"
	"github.com/spf13/cobra"

	"github.com/suyuan32/goctls/extra/ent/importschema/mux"
)

var (
	VarStringDSN           string
	VarStringOutputDir     string
	VarStringTables        string
	VarStringExcludeTables string
	VarBoolAutoMixin       bool
	VarBoolPluralTable     bool
)

type GenContext struct {
	Dsn           string
	OutputDir     string
	Tables        []string
	ExcludeTables []string
	AutoMixin     bool
	PluralTable   bool
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
		} else {
			return errors.New("failed to find ent directory")
		}
	} else if !strings.HasSuffix(VarStringOutputDir, "/") {
		VarStringOutputDir = VarStringOutputDir + "/ent/schema/"

		outputDir, err = filepath.Abs(VarStringOutputDir)
		if err != nil {
			return err
		}
	} else {
		VarStringOutputDir = VarStringOutputDir + "ent/schema/"

		outputDir, err = filepath.Abs(VarStringOutputDir)
		if err != nil {
			return err
		}
	}

	ctx := &GenContext{}
	ctx.OutputDir = outputDir
	ctx.Dsn = VarStringDSN
	ctx.AutoMixin = VarBoolAutoMixin
	ctx.PluralTable = VarBoolPluralTable
	if len(VarStringTables) == 0 {
		ctx.Tables = nil
	} else {
		ctx.Tables = strings.Split(VarStringTables, ",")
	}

	if len(VarStringExcludeTables) == 0 {
		ctx.ExcludeTables = nil
	} else {
		ctx.ExcludeTables = strings.Split(VarStringExcludeTables, ",")
	}

	drv, err := mux.Default.OpenImport(ctx.Dsn)
	if err != nil {
		return fmt.Errorf("failed to create import driver - %v", err)
	}

	i, err := NewImport(
		WithTables(ctx.Tables),
		WithExcludedTables(ctx.ExcludeTables),
		WithDriver(drv),
	)
	if err != nil {
		return fmt.Errorf("create importer failed: %v", err)
	}

	mutations, err := i.SchemaMutations(context.Background())
	if err != nil {
		return fmt.Errorf("schema import failed - %v", err)
	}

	if len(mutations) == 0 {
		return errors.New("table not found: " + VarStringTables)
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
