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

package rpcconfgen

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var (
	VarStringServiceName string
	VarStringRpcDir      string
	VarIntPort           int
	VarStringApiDir      string
)

type GenContext struct {
	ServiceName string
	ClientPath  string
	RpcDir      string
	Port        int
	ApiDir      string
}

func Gen(_ *cobra.Command, _ []string) error {
	if VarStringRpcDir == "" {
		return errors.New("rpc dir cannot be empty")
	}

	return DoGen(&GenContext{
		ServiceName: VarStringServiceName,
		RpcDir:      VarStringRpcDir,
		Port:        VarIntPort,
		ApiDir:      VarStringApiDir,
	})
}

func DoGen(g *GenContext) error {
	color.Green.Println("Generating...")

	rpcPath, err := filepath.Abs(g.RpcDir)
	if err != nil {
		return err
	}

	apiPath, err := filepath.Abs(g.ApiDir)
	if err != nil {
		return err
	}

	rpcGoMod, err := fileutil.ReadFileByLine(filepath.Join(rpcPath, "go.mod"))
	if err != nil {
		return err
	}

	packageSrc := strings.TrimSpace(strings.Split(rpcGoMod[0], " ")[1])

	if packageSrc == "" {
		return errors.New("failed to get module src")
	}

	etcFiles, _ := fileutil.ListFileNames(filepath.Join(apiPath, "etc"))

	for _, file := range etcFiles {
		fileutil.WriteStringToFile(filepath.Join(apiPath, "etc", file), fmt.Sprintf("\n\n%sRpc:\n  Target: k8s://default/%s-rpc-svc:%d\n  Enabled: true\n", g.ServiceName, strings.ToLower(g.ServiceName), g.Port), true)
	}

	configFilePath := filepath.Join(apiPath, "internal", "config", "config.go")

	configFileStr, _ := fileutil.ReadFileToString(configFilePath)

	configInsertIdx := strings.Index(configFileStr, "Config struct {")
	if configInsertIdx == -1 {
		return errors.New("failed to find insert place in config")
	}

	fileutil.WriteStringToFile(configFilePath, configFileStr[:configInsertIdx+16]+fmt.Sprintf("\t%s    zrpc.RpcClientConf", g.ServiceName)+configFileStr[configInsertIdx+16:], false)

	svcStrBuilder := strings.Builder{}

	svcFilePath := filepath.Join(apiPath, "internal", "svc", "service_context.go")

	svcFileStr, _ := fileutil.ReadFileToString(svcFilePath)

	importInsertIdx := strings.Index(svcFileStr, "import (")

	svcStrBuilder.WriteString(svcFileStr[:importInsertIdx+8])

	if importInsertIdx == -1 {
		return errors.New("failed to find import place in service context.go")
	}

	svcStrBuilder.WriteString(fmt.Sprintf("\n\t \"%s/%sclient\"", packageSrc, strings.ToLower(g.ServiceName)))

	svcTypeInsertIdx := strings.Index(svcFileStr, "type ServiceContext struct {")

	if svcTypeInsertIdx == -1 {
		return errors.New("failed to find insert place in svc")
	}

	svcStrBuilder.WriteString(svcFileStr[importInsertIdx+8 : svcTypeInsertIdx+29])

	svcStrBuilder.WriteString(fmt.Sprintf("\n\t%sRpc   %sclient.%s", g.ServiceName, strings.ToLower(g.ServiceName), g.ServiceName))

	svcReturnInsertIdx := strings.Index(svcFileStr, "&ServiceContext{")

	if svcReturnInsertIdx == -1 {
		return errors.New("failed to find return insert place in svc")
	}

	svcStrBuilder.WriteString(svcFileStr[svcTypeInsertIdx+29 : svcReturnInsertIdx+16])

	svcStrBuilder.WriteString(fmt.Sprintf("\n\t\t%sRpc:   %sclient.New%s(zrpc.NewClientIfEnable(c.%sRpc)),", g.ServiceName, strings.ToLower(g.ServiceName), g.ServiceName, g.ServiceName))

	svcStrBuilder.WriteString(svcFileStr[svcReturnInsertIdx+16:])

	err = fileutil.WriteStringToFile(svcFilePath, svcStrBuilder.String(), false)
	if err != nil {
		return err
	}

	return nil
}
