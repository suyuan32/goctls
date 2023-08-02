package cicd

import (
	"github.com/suyuan32/goctls/cicd/drone"
	"github.com/suyuan32/goctls/internal/cobrax"
)

var (
	CicdCmd  = cobrax.NewCommand("cicd")
	DroneCmd = cobrax.NewCommand("drone", cobrax.WithRunE(drone.GenDrone))
)

func init() {
	var (
		droneCmdFlags = DroneCmd.Flags()
	)

	droneCmdFlags.StringVarP(&drone.VarDroneName, "drone_name", "d")
	droneCmdFlags.StringVarPWithDefaultValue(&drone.VarGitGoPrivate, "go_private", "g", "gitee.com")
	droneCmdFlags.StringVarP(&drone.VarServiceName, "service_name", "s")
	droneCmdFlags.StringVarP(&drone.VarServiceType, "service_type", "x")
	droneCmdFlags.StringVarPWithDefaultValue(&drone.VarGitBranch, "git_branch", "b", "master")
	droneCmdFlags.StringVarP(&drone.VarRegistry, "registry", "r")
	droneCmdFlags.StringVarP(&drone.VarRepo, "repo", "o")
	droneCmdFlags.StringVarP(&drone.VarEtcYaml, "etc_yaml", "e")

	CicdCmd.AddCommand(DroneCmd)
}
