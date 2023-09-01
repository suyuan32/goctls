package upgrade

import (
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/suyuan32/goctls/config"
	"github.com/suyuan32/goctls/rpc/execx"
	"path/filepath"
	"strings"
)

func upgradeDependencies(workDir string) error {
	// drop old replace
	oldVersion := []string{"v1.5.2", "v1.5.3", "v1.5.4"}
	for _, v := range oldVersion {
		_, err := execx.Run(fmt.Sprintf("go mod edit -dropreplace github.com/zeromicro/go-zero@%s", v), workDir)
		if err != nil {
			return errors.New("failed to drop old replace")
		}
	}

	data, err := fileutil.ReadFileToString(filepath.Join(workDir, "go.mod"))
	if err != nil {
		return err
	}

	err = upgradeOfficialDependencies(data, workDir)
	if err != nil {
		return err
	}

	err = tidy()
	if err != nil {
		return err
	}

	return nil
}

func upgradeOfficialDependencies(data, workDir string) (err error) {
	deps := []struct {
		Name string
		Repo string
	}{
		{
			Name: "simple-admin-common",
			Repo: "github.com/suyuan32/simple-admin-common",
		},
		{
			Name: "simple-admin-core",
			Repo: "github.com/suyuan32/simple-admin-core",
		},
		{
			Name: "simple-admin-message-center",
			Repo: "github.com/suyuan32/simple-admin-message-center",
		},
	}

	for _, v := range deps {
		if strings.Contains(data, v.Name) {
			_, err = execx.Run(fmt.Sprintf("go mod edit -require=%s@%s", v.Repo,
				config.CoreVersion), workDir)
			if err != nil {
				return err
			}
		}
	}

	return err
}
