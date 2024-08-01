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

package i18n

import (
	"errors"

	"github.com/gookit/color"

	"github.com/spf13/cobra"

	"github.com/suyuan32/goctls/extra/i18n/api"
)

var (
	// VarStringTarget describes the target.
	VarStringTarget string
	// VarStringModelName describes the model name
	VarStringModelName string
	// VarStringModelNameZh describes the model's Chinese translation
	VarStringModelNameZh string
	// VarStringOutputDir describes the output directory
	VarStringOutputDir string
)

func Gen(_ *cobra.Command, _ []string) error {
	err := Validate()
	if err != nil {
		return err
	}
	return DoGen()
}

func DoGen() error {
	switch VarStringTarget {
	case "api":
		ctx := &api.GenContext{
			Target:      VarStringTarget,
			ModelName:   VarStringModelName,
			ModelNameZh: VarStringModelNameZh,
			OutputDir:   VarStringOutputDir,
		}
		return api.GenApiI18n(ctx)
	}

	color.Green.Println("Generate successfully")
	return errors.New("invalid target, try \"api\"")
}

func Validate() error {
	if VarStringTarget == "" {
		return errors.New("the target cannot be empty, use --target to set it")
	} else if VarStringModelName == "" {
		return errors.New("the model name cannot be empty, use --model_name to set it")
	}
	return nil
}
