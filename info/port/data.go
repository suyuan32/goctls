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

package port

import (
	"os"

	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/suyuan32/goctls/util/env"
)

var envInfo table.Writer
var lang = env.IsChinaEnv()

// portInfo show the port usage across the simple admin
func portInfo() string {
	color.Green.Println("PORT")
	envInfo = table.NewWriter()
	envInfo.SetOutputMirror(os.Stdout)
	if lang {
		envInfo.AppendHeader(table.Row{"端口", "服务"})

	} else {
		envInfo.AppendHeader(table.Row{"Port", "Service"})
	}
	envInfo.AppendRows([]table.Row{
		{9100, "core_api"},
		{9101, "core_rpc"},
		{9102, "file_api"},
		{9103, "member_api"},
		{9104, "member_rpc"},
		{9105, "job_rpc"},
		{9106, "mcms_rpc"},
	})
	return envInfo.Render()
}
