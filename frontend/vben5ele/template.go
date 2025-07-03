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

package vben5ele

import (
	_ "embed"
)

var (
	//go:embed api.tpl
	apiTpl string

	//go:embed model.tpl
	modelTpl string

	//go:embed schemas.tpl
	dataTpl string

	//go:embed index.tpl
	indexTpl string

	//go:embed locale.tpl
	localeTpl string

	//go:embed statusrender.tpl
	statusRenderTpl string

	//go:embed staterender.tpl
	stateRenderTpl string

	//go:embed modal_index.tpl
	modalIndexTpl string

	//go:embed form.tpl
	formTpl string
)
