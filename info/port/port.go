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
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// ShowPort is used to show port usages across the simple admin.
func ShowPort(_ *cobra.Command, _ []string) error {
	if lang {
		color.Green.Println("Simple Admin的端口使用情况")
	} else {
		color.Green.Println("Simple Admin's port usage")
	}
	portInfo()
	return nil
}
