package project

import (
	"github.com/suyuan32/goctls/internal/cobrax"
	"github.com/suyuan32/goctls/project/upgrade"
)

var (
	Cmd        = cobrax.NewCommand("project")
	upgradeCmd = cobrax.NewCommand("upgrade", cobrax.WithRunE(upgrade.UpgradeProject))
)

func init() {
	upgradeCmdFlag := upgradeCmd.Flags()

	upgradeCmdFlag.BoolVarP(&upgrade.VarBoolUpgradeMakefile, "makefile", "m")

	Cmd.AddCommand(upgradeCmd)
}
