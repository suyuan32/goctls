package gitlab

import (
	_ "embed"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"path/filepath"
)

//go:embed gitlab.tpl
var gitlabTpl string

var (
	VarStringOutputDir string
)

func Gen(_ *cobra.Command, _ []string) (err error) {
	abs, err := filepath.Abs(VarStringOutputDir)
	if err != nil {
		return errors.Wrap(err, "dir not found")
	}

	color.Green.Println("Generating...")

	err = fileutil.WriteStringToFile(filepath.Join(abs, ".gitlab-ci.yml"), gitlabTpl, false)
	if err != nil {
		return errors.Wrap(err, "write file failed")
	}

	return err
}
