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
	for _, v := range config.OldGoZeroVersion {
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
		Repo string
	}{
		{
			Repo: "github.com/suyuan32/simple-admin-common",
		},
		{
			Repo: "github.com/suyuan32/simple-admin-core",
		},
		{
			Repo: "github.com/suyuan32/simple-admin-message-center",
		},
		{
			Repo: "github.com/suyuan32/simple-admin-job",
		},
	}

	for _, v := range deps {
		if strings.Contains(data, v.Repo) {
			_, err = execx.Run(fmt.Sprintf("go mod edit -require=%s@%s", v.Repo,
				config.CoreVersion), workDir)
			if err != nil {
				return err
			}
		}
	}

	return err
}
