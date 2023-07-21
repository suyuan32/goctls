package bug

import (
	"github.com/spf13/cobra"
	"github.com/suyuan32/goctls/internal/cobrax"
)

// Cmd describes a bug command.
var Cmd = cobrax.NewCommand("bug", cobrax.WithRunE(cobra.NoArgs), cobrax.WithArgs(cobra.NoArgs))
