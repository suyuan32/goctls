package env

import "github.com/suyuan32/goctls/internal/cobrax"

var (
	sliceVarWriteValue []string
	boolVarForce       bool
	boolVarVerbose     bool
	boolVarInstall     bool
	boolVarClearCache  bool

	// Cmd describes an env command.
	Cmd        = cobrax.NewCommand("env", cobrax.WithRunE(write))
	installCmd = cobrax.NewCommand("install", cobrax.WithRunE(install))
	checkCmd   = cobrax.NewCommand("check", cobrax.WithRunE(check))
)

func init() {
	// The root command flags
	Cmd.Flags().StringSliceVarP(&sliceVarWriteValue, "write", "w")
	Cmd.PersistentFlags().BoolVarP(&boolVarForce, "force", "f")
	Cmd.PersistentFlags().BoolVarP(&boolVarVerbose, "verbose", "v")

	// The sub-command flags
	checkCmd.Flags().BoolVarP(&boolVarInstall, "install", "i")
	checkCmd.Flags().BoolVarP(&boolVarClearCache, "clear_cache", "c")

	// Add sub-command
	Cmd.AddCommand(checkCmd, installCmd)
}
