package tmpl

import _ "embed"

var (
	//go:embed authorty.tpl
	AuthorityTmpl string

	//go:embed authorty_tenant.tpl
	AuthorityTenantTmpl string

	//go:embed data_perm.tpl
	DataPermTmpl string
)
