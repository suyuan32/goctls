package drone

import (
	"github.com/suyuan32/goctls/internal/cobrax"
)

var (
	// Cmd describes an api command.
	Cmd = cobrax.NewCommand("drone", cobrax.WithRunE(GenDrone))
)

func init() {
	var (
		droneCmdFlags = Cmd.Flags()
	)

	droneCmdFlags.StringVarP(&VarDroneName, "drone_name", "d")
	droneCmdFlags.StringVarPWithDefaultValue(&VarGitGoPrivate, "go_private", "g", "gitee.com")
	droneCmdFlags.StringVarP(&VarServiceName, "service_name", "s")
	droneCmdFlags.StringVarP(&VarServiceType, "service_type", "x")
	droneCmdFlags.StringVarPWithDefaultValue(&VarGitBranch, "git_branch", "b", "master")
	droneCmdFlags.StringVarP(&VarRegistry, "registry", "r")
	droneCmdFlags.StringVarP(&VarRepo, "repo", "o")
	droneCmdFlags.StringVarP(&VarEtcYaml, "etc_yaml", "e")
}
