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

package logviewer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var (
	VarStringFilePath         string
	VarIntMessageCapacity     int
	VarStringWorkspaceSetting string
	VarStringWorkspace        string
	VarStringLogType          string
	VarBoolResetWorkspace     bool
	VarBoolList               bool
	VarStringRemoveConfig     string
)

func Gen(_ *cobra.Command, _ []string) (err error) {
	logData := ""
	if VarStringFilePath != "" {
		filePath, err := filepath.Abs(VarStringFilePath)
		if err != nil {
			return err
		}
		logData, err = fileutil.ReadFileToString(filePath)
		if err != nil {
			return err
		}
	} else if VarStringWorkspaceSetting != "" {
		workspaceData := strings.Split(VarStringWorkspaceSetting, ",")
		if len(workspaceData) < 2 {
			return errors.New("wrong workspace setting, make sure the format: \"name,dir-path\"")
		}
		if names, err := fileutil.ListFileNames(workspaceData[1]); err != nil {
			return errors.Join(err, errors.New("failed to access the path"))
		} else {
			if len(names) < 5 {
				return errors.Join(err, errors.New("the folder does not contain the log files"))
			}
		}

		userDir, err := user.Current()
		if err != nil {
			return errors.Join(err, errors.New("failed to get the user directory to store data"))
		}

		configFile := filepath.Join(userDir.HomeDir, ".goctls/log_workspace_config.txt")
		if !fileutil.IsExist(configFile) {
			err = fileutil.CreateDir(filepath.Join(userDir.HomeDir, ".goctls") + "/")
			if err != nil {
				return err
			}
			fileutil.CreateFile(configFile)
		}

		configData, err := fileutil.ReadFileToString(configFile)
		if err != nil {
			return errors.Join(err, errors.New("failed to read config file"))
		}

		workspaceConfigData := strings.Split(configData, "\n")
		isAppend := true

		for i := 0; i < len(workspaceConfigData); i++ {
			if strings.Contains(workspaceConfigData[i], workspaceData[0]) {
				workspaceConfigData[i] = VarStringWorkspaceSetting
				isAppend = false
			}
		}

		if isAppend {
			err = fileutil.WriteStringToFile(configFile, VarStringWorkspaceSetting+"\n", true)
			if err != nil {
				return errors.Join(err, errors.New("failed to write workspace data to file"))
			}
		} else {
			err = fileutil.WriteStringToFile(configFile, strings.Join(workspaceConfigData, "\n"), true)
			if err != nil {
				return errors.Join(err, errors.New("failed to write workspace data to file"))
			}
		}

		color.Green.Println("set workspace successfully")
		return nil
	} else if VarStringWorkspace != "" {
		configFile, err := getWorkspaceConfigDir()
		if err != nil {
			return err
		}

		configData, err := fileutil.ReadFileToString(configFile)
		if err != nil {
			return errors.Join(err, errors.New("failed to read config file"))
		}

		workspaceConfigData := strings.Split(configData, "\n")
		for i := 0; i < len(workspaceConfigData); i++ {
			if strings.Contains(workspaceConfigData[i], VarStringWorkspace) {
				workspaceData := strings.Split(workspaceConfigData[i], ",")
				logData, err = fileutil.ReadFileToString(filepath.Join(workspaceData[1], VarStringLogType+".log"))
				if err != nil {
					return errors.Join(err, errors.New("failed to read log file"))
				}
			}
		}
	} else if VarBoolResetWorkspace {
		configFile, err := getWorkspaceConfigDir()
		if err != nil {
			return err
		}

		err = fileutil.ClearFile(configFile)
		if err != nil {
			return errors.Join(err, errors.New("failed to reset config file"))
		}

		color.Green.Println("Reset workspace configuration successfully")
	} else if VarBoolList {
		configFile, err := getWorkspaceConfigDir()
		if err != nil {
			return err
		}

		configData, err := fileutil.ReadFileToString(configFile)
		if err != nil {
			return err
		}

		fmt.Println(strings.ReplaceAll(configData, ",", "    "))
	} else if VarStringRemoveConfig != "" {
		configFile, err := getWorkspaceConfigDir()
		if err != nil {
			return err
		}

		err = removeConfig(VarStringRemoveConfig, configFile)
		if err != nil {
			return err
		}

		color.Green.Println(fmt.Sprintf("Remove %s configuration successfully", VarStringRemoveConfig))
		return nil
	}

	err = prettierJsonData(logData, VarIntMessageCapacity)

	return err
}

func prettierJsonData(data string, capacity int) error {
	messageData := strings.Split(data, "\n")
	messageData = messageData[:len(messageData)-1]

	var messageDataCut []string
	if len(messageData) < capacity {
		messageDataCut = messageData
	} else {
		messageDataCut = messageData[len(messageData)-capacity:]
	}

	for i := len(messageDataCut) - 1; i >= 0; i-- {
		if len(messageDataCut[i]) < 2 {
			continue
		}

		tmp, err := beautifyJsonData(messageDataCut[i])
		if err != nil {
			return err
		}
		fmt.Println(tmp)
	}

	return nil
}

func beautifyJsonData(data string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(data), "", "    "); err != nil {
		return "", err
	}

	result := strings.ReplaceAll(prettyJSON.String(), "\\n", "\n    ")
	result = strings.ReplaceAll(result, "\\t", "    ")
	result = strings.ReplaceAll(result, "\"@timestamp\"", color.Green.Sprintf("\"@timestamp\""))
	result = strings.ReplaceAll(result, "\"content\"", color.Red.Sprintf("\"content\""))
	return result, nil
}

func removeConfig(target, configPath string) error {
	removeTarget := strings.Split(target, ",")
	originalData, err := fileutil.ReadFileToString(configPath)
	if err != nil {
		return err
	}

	originConfigs := strings.Split(originalData, "\n")

	var output []string

	for _, v := range originConfigs {
		isRemove := false
		for _, v1 := range removeTarget {
			if strings.Contains(v, v1) {
				isRemove = true
			}
		}

		if !isRemove {
			output = append(output, v)
		}
	}

	err = fileutil.WriteStringToFile(configPath, strings.Join(output, "\n"), false)
	if err != nil {
		return err
	}

	return nil
}
