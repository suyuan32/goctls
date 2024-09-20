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

package frontend

import (
	"github.com/suyuan32/goctls/frontend/vben"
	"github.com/suyuan32/goctls/internal/cobrax"
)

var (
	// Cmd describes an api command.
	Cmd     = cobrax.NewCommand("frontend")
	VbenCmd = cobrax.NewCommand("vben", cobrax.WithRunE(vben.GenCRUDLogic))
)

func init() {
	vbenCmdFlags := VbenCmd.Flags()

	vbenCmdFlags.StringVarPWithDefaultValue(&vben.VarStringOutput, "output", "o", "./")
	vbenCmdFlags.StringVarP(&vben.VarStringApiFile, "api_file", "a")
	vbenCmdFlags.StringVarPWithDefaultValue(&vben.VarStringFolderName, "folder_name", "f", "sys")
	vbenCmdFlags.StringVarP(&vben.VarStringSubFolder, "sub_folder", "s")
	vbenCmdFlags.StringVarPWithDefaultValue(&vben.VarStringApiPrefix, "prefix", "p", "sys-api")
	vbenCmdFlags.StringVarP(&vben.VarStringModelName, "model_name", "m")
	vbenCmdFlags.StringVarPWithDefaultValue(&vben.VarStringFormType, "form_type", "t", "drawer")
	vbenCmdFlags.BoolVarP(&vben.VarBoolOverwrite, "overwrite", "w")
	vbenCmdFlags.StringVar(&vben.VarStringModelChineseName, "model_chinese_name")
	vbenCmdFlags.StringVar(&vben.VarStringModelEnglishName, "model_english_name")

	Cmd.AddCommand(VbenCmd)
}
