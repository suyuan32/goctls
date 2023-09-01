package upgrade

import (
	"errors"
	"fmt"
	"github.com/suyuan32/goctls/rpc/execx"
	"github.com/suyuan32/goctls/util/ctx"
	"os"
)

const (
	goZeroMod = "github.com/zeromicro/go-zero"
	adminTool = "github.com/suyuan32/simple-admin-tools"
)

var errInvalidGoMod = errors.New("it's only working for go module")

func editMod(zeroVersion, toolVersion string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	isGoMod, _ := ctx.IsGoMod(wd)
	if !isGoMod {
		return nil
	}

	mod := fmt.Sprintf("%s@%s", goZeroMod, zeroVersion)

	err = addRequire(mod)
	if err != nil {
		return err
	}

	// add replace
	mod = fmt.Sprintf("%s@%s=%s@%s", goZeroMod, zeroVersion, adminTool, toolVersion)

	err = addReplace(mod)
	if err != nil {
		return err
	}

	return nil
}

func addRequire(mod string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	isGoMod, _ := ctx.IsGoMod(wd)
	if !isGoMod {
		return errInvalidGoMod
	}

	_, err = execx.Run("go mod edit -require "+mod, wd)
	return err
}

func addReplace(mod string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	isGoMod, _ := ctx.IsGoMod(wd)
	if !isGoMod {
		return errInvalidGoMod
	}

	_, err = execx.Run("go mod edit -replace "+mod, wd)
	return err
}

func tidy() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	isGoMod, _ := ctx.IsGoMod(wd)
	if !isGoMod {
		return nil
	}

	_, err = execx.Run("go mod tidy", wd)
	return err
}
