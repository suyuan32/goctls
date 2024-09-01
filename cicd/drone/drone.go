package drone

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/suyuan32/goctls/util/pathx"
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
	color.Green.Println("verifying parameters...")

	// 校验模版逻辑
	etcYaml := VarEtcYaml
	if len(etcYaml) == 0 {
		return fmt.Errorf("config file not found, please check the etc path and yaml file")
	}

	droneName := VarDroneName
	if len(droneName) == 0 {
		droneName = "drone-greet"
	}

	goPrivate := VarGitGoPrivate
	if len(strings.Split(goPrivate, ".")) <= 1 {
		return fmt.Errorf("wrong private repository address, set like: gitee.com, github.com, gitlab.com")
	}

	serviceName := VarServiceName
	if len(serviceName) < 1 {
		return fmt.Errorf("service name is empty, please set it")
	}
	serviceName = strings.TrimSuffix(serviceName, ".go")

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
		return fmt.Errorf("registry is empty, please set your docker registry address such as \"registry.cn-beijing.aliyuncs.com\"")
	}

	repo := VarRepo
	if len(repo) == 0 {
		return fmt.Errorf("repo is empty, please set your docker repo address such as \"registry.cn-hangzhou.aliyuncs.com/simple_admin/core-api-docker:v1.1.0\" ")
	}

	color.Green.Render("loading template...")

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
		DroneName:    droneName,
		GitGoPrivate: goPrivate,
		ServiceName:  serviceName,
		ServiceType:  serviceType,
		GitBranch:    gitBranch,
		Registry:     registry,
		Repo:         repo,
	})

	t1 := template.Must(template.New("dockerfile").Parse(dockerfileTpl))
	t1.Execute(dockerfileFile, Dockerfile{
		EtcYaml: etcYaml,
	})

	color.Green.Println("Done.")
	return nil
}
