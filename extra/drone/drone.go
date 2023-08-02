package drone

import (
	_ "embed"
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/suyuan32/goctls/util/pathx"
	"os"
	"strings"
	"text/template"
)

var (
	//go:embed drone.tpl
	DroneTpl string

	//go:embed dockerfile.tpl
	DockerfileTpl string
)

var (
	VarDroneName    string
	VarGitGoPrivate string
	VarServiceName  string
	VarServiceType  string
	VarGitBranch    string
	VarRegistry     string
	VarRepo         string
	VarEtcYaml      string
)

type Drone struct {
	//步骤三
	DroneName    string
	GitGoPrivate string
	ServiceName  string
	ServiceType  string
	GitBranch    string
	Registry     string
	Repo         string
}

type Dockerfile struct {
	EtcYaml string
}

func GenDrone(_ *cobra.Command, _ []string) error {
	fmt.Println(color.Green.Render("verifying params..."))

	// 校验模版逻辑
	etcyaml := VarEtcYaml
	if len(etcyaml) == 0 {
		return fmt.Errorf("etcyaml is empty!")
	}

	dronename := VarDroneName
	if len(dronename) == 0 {
		dronename = "drone-greet"
	}

	goprivate := VarGitGoPrivate
	if len(strings.Split(goprivate, ".")) <= 1 {
		return fmt.Errorf("error go private!")
	}

	serviceName := VarServiceName
	if len(strings.Split(serviceName, ".go")) != 1 {
		return fmt.Errorf("please ignore suffix .go!")
	}

	serviceType := VarServiceType
	if len(serviceType) == 0 {
		// build happy 😄
		serviceType = "happy"
	}
	gitBranch := VarGitBranch
	if len(gitBranch) == 0 {
		gitBranch = "main"
	}
	registry := VarRegistry
	if len(registry) == 0 {
		return fmt.Errorf("registry is empty!")
	}

	repo := VarRepo
	if len(repo) == 0 {
		return fmt.Errorf("repo is empty!")
	}

	fmt.Println(color.Green.Render("loading template..."))

	// 创建 .drone.yml 前面的点是drone默认加载程序，如果脱离本框架会无法找到路径
	droneFile, err := os.Create(".drone.yml")
	if err != nil {
		return err
	}

	dockerfileFile, err := os.Create("Dockerfile")
	if err != nil {
		return err
	}

	// 加载模板
	droneTpl, err := pathx.LoadTemplate("drone", "drone.tpl", DroneTpl)
	if err != nil {
		return err
	}

	dockerfileTpl, err := pathx.LoadTemplate("dockerfile", "dockerfile.tpl", DockerfileTpl)
	if err != nil {
		return err
	}

	// 渲染模板
	t := template.Must(template.New("drone").Parse(droneTpl))
	t.Execute(droneFile, Drone{
		DroneName:    dronename,
		GitGoPrivate: goprivate,
		ServiceName:  serviceName,
		ServiceType:  serviceType,
		GitBranch:    gitBranch,
		Registry:     registry,
		Repo:         repo,
	})

	t1 := template.Must(template.New("dockerfile").Parse(dockerfileTpl))
	t1.Execute(dockerfileFile, Dockerfile{
		EtcYaml: etcyaml,
	})

	fmt.Println(color.Green.Render("Done."))
	return nil
}
