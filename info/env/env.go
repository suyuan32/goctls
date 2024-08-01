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

package env

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var (
	ServiceName string
	ShowList    bool
)

// ShowEnv is used to show the environment variable usages.
func ShowEnv(_ *cobra.Command, _ []string) error {
	if ShowList {
		getServiceList()
		return nil
	}

	if lang {
		color.Green.Println("Simple Admin的环境变量")
		color.Red.Println("注意： 环境变量的优先级大于配置文件")
	} else {
		color.Green.Println("Simple Admin's environment variables")
		color.Red.Println("Notice: Environment variables have priority over configuration files")
	}

	switch ServiceName {
	case "core":
		toolEnvInfo()
		authEnvInfo()
		crosEnvInfo()
		apiEnvInfo()
		rpcEnvInfo()
		logEnvInfo()
		databaseEnvInfo()
		redisEnvInfo()
		i18nEnvInfo()
		captchaEnvInfo()
	case "fms":
		fmsEnvInfo()
	}

	return nil
}
