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

package template

import _ "embed"

var (
	//go:embed tmpl/pagination.tmpl
	PaginationTmpl string

	//go:embed tmpl/set_not_nil.tmpl
	NotNilTmpl string

	//go:embed tmpl/set_or_clear.tmpl
	SetOrClearTmpl string

	//go:embed tmpl/cache.tmpl
	CacheTmpl string

	//go:embed tmpl/cache_zero.tmpl
	CacheZeroTmpl string

	//go:embed tmpl/tenant_privacy.tmpl
	TenantPrivacyTmpl string
)
