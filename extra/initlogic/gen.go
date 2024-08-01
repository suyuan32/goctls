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

package initlogic

import (
	_ "embed"
	"errors"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// VarStringTarget describes the target.
	VarStringTarget string
	// VarStringModelName describes the model name
	VarStringModelName string
	// VarStringRoutePrefix describes the prefix of route path
	VarStringRoutePrefix string
	// VarStringOutputPath describes the output directory
	VarStringOutputPath string
	// VarStringStyle describes the file naming style
	VarStringStyle string
	// VarServiceName describes the service name
	VarServiceName string
)

func Gen(_ *cobra.Command, _ []string) error {
	err := Validate()
	if err != nil {
		return err
	}

	routePrefix := ""
	if VarStringRoutePrefix != "" {
		if strings.HasPrefix(VarStringRoutePrefix, "/") {
			routePrefix = VarStringRoutePrefix
		} else {
			routePrefix = "/" + VarStringRoutePrefix
		}
	}

	ctx := &CoreGenContext{
		Target:      VarStringTarget,
		ModelName:   VarStringModelName,
		Output:      VarStringOutputPath,
		Style:       VarStringStyle,
		ServiceName: VarServiceName,
		RoutePrefix: routePrefix,
	}

	return DoGen(ctx)
}

func DoGen(g *CoreGenContext) error {
	if g.Target == "core" {
		return GenCore(g)
	} else if g.Target == "other" {
		return OtherGen(g)
	}
	return errors.New("invalid target, try \"core\" or \"other\"")
}

func Validate() error {
	if VarStringTarget == "" {
		return errors.New("the target cannot be empty, use --target to set it")
	} else if VarStringModelName == "" {
		return errors.New("the model name cannot be empty, use --model_name to set it")
	}
	return nil
}
