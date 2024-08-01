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

package info

import (
	"github.com/suyuan32/goctls/info/env"
	"github.com/suyuan32/goctls/info/port"
	"github.com/suyuan32/goctls/internal/cobrax"
)

var (
	// Cmd describes a docker command.
	Cmd = cobrax.NewCommand("info")

	EnvCmd = cobrax.NewCommand("env", cobrax.WithRunE(env.ShowEnv))

	PortCmd = cobrax.NewCommand("port", cobrax.WithRunE(port.ShowPort))
)

func init() {

	var (
		envCmdFlags = EnvCmd.Flags()
	)

	envCmdFlags.StringVarPWithDefaultValue(&env.ServiceName, "service_name", "s", "core")
	envCmdFlags.BoolVarP(&env.ShowList, "list", "l")

	Cmd.AddCommand(EnvCmd)
	Cmd.AddCommand(PortCmd)
}
