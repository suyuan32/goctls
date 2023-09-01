package upgrade

import (
	"errors"
	"fmt"
	"github.com/suyuan32/goctls/config"
	"github.com/suyuan32/goctls/rpc/execx"
	"strings"
)

func upgradeDependencies(workdir string) error {
	// drop old replace
	oldVersion := []string{"v1.5.2", "v1.5.3", "v1.5.4"}
	for _, v := range oldVersion {
		_, err := execx.Run(fmt.Sprintf("go mod edit -dropreplace github.com/zeromicro/go-zero@%s", v), workdir)
		if err != nil {
			return errors.New("failed to drop old replace")
		}
	}

	data, err := execx.Run("go list -json", workdir)
	if err != nil {
		return err
	}

	if strings.Contains(data, "simple-admin-common") {
		_, err = execx.Run(fmt.Sprintf("go mod edit -require=github.com/suyuan32/simple-admin-common@%s",
			config.CoreVersion), workdir)
		if err != nil {
			return err
		}
	}

	if strings.Contains(data, "simple-admin-core") {
		_, err = execx.Run(fmt.Sprintf("go mod edit -require=github.com/suyuan32/simple-admin-core@%s",
			config.CoreVersion), workdir)
		if err != nil {
			return err
		}
	}

	err = tidy()
	if err != nil {
		return err
	}

	return nil
}
