package upgrade

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/suyuan32/goctls/rpc/execx"
)

// upgrade gets the latest goctl by
// go install github.com/suyuan32/goctls@latest
func upgrade(_ *cobra.Command, _ []string) error {
	cmd := `go install github.com/suyuan32/goctls@latest`
	info, err := execx.Run(cmd, "")
	if err != nil {
		return err
	}

	fmt.Print(info)
	color.Green.Println("Upgrade successfully")
	return nil
}
