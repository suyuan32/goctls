package upgrade

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/suyuan32/goctls/rpc/execx"
)

// upgrade gets the latest goctl by
// go install github.com/suyuan32/goctls@latest
func upgrade(_ *cobra.Command, _ []string) error {
	cmd := `GO111MODULE=on GOPROXY=https://goproxy.cn/,direct go install github.com/suyuan32/goctls@latest`
	if runtime.GOOS == "windows" {
		cmd = `set GOPROXY=https://goproxy.cn,direct && go install github.com/suyuan32/goctls@latest`
	}
	info, err := execx.Run(cmd, "")
	if err != nil {
		return err
	}

	fmt.Print(info)
	return nil
}
