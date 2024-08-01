package docker

import "github.com/suyuan32/goctls/internal/cobrax"

var (
	varServiceName    string
	varServiceType    string
	varStringBase     string
	varIntPort        int
	varStringHome     string
	varStringRemote   string
	varStringBranch   string
	varStringImage    string
	varStringTZ       string
	varBoolChina      bool
	varStringAuthor   string
	varBoolLocalBuild bool

	// Cmd describes a docker command.
	Cmd = cobrax.NewCommand("docker", cobrax.WithRunE(dockerCommand))
)

func init() {
	dockerCmdFlags := Cmd.Flags()
	dockerCmdFlags.StringVarP(&varServiceName, "service_name", "s")
	dockerCmdFlags.StringVarPWithDefaultValue(&varServiceType, "service_type", "t", "rpc")
	dockerCmdFlags.StringVarPWithDefaultValue(&varStringBase, "base", "a", "alpine:3.20")
	dockerCmdFlags.IntVarP(&varIntPort, "port", "p")
	dockerCmdFlags.StringVarP(&varStringHome, "home", "m")
	dockerCmdFlags.StringVarP(&varStringRemote, "remote", "r")
	dockerCmdFlags.StringVarP(&varStringBranch, "branch", "b")
	dockerCmdFlags.BoolVarP(&varBoolChina, "china", "c")
	dockerCmdFlags.BoolVarP(&varBoolLocalBuild, "local_build", "l")
	dockerCmdFlags.StringVarPWithDefaultValue(&varStringImage, "image", "i", "golang:1.22.5-alpine3.20")
	dockerCmdFlags.StringVarP(&varStringTZ, "tz", "z")
	dockerCmdFlags.StringVarPWithDefaultValue(&varStringAuthor, "author", "u", "example@example.com")
}
