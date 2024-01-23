package env

import (
	"fmt"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/suyuan32/goctls/util/pathx"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/suyuan32/goctls/pkg/env"
	"github.com/suyuan32/goctls/pkg/goswagger"
	"github.com/suyuan32/goctls/pkg/protoc"
	"github.com/suyuan32/goctls/pkg/protocgengo"
	"github.com/suyuan32/goctls/pkg/protocgengogrpc"
	"github.com/suyuan32/goctls/util/console"
)

type bin struct {
	name   string
	exists bool
	get    func(cacheDir string) (string, error)
}

var bins = []bin{
	{
		name:   "protoc",
		exists: protoc.Exists(),
		get:    protoc.Install,
	},
	{
		name:   "protoc-gen-go",
		exists: protocgengo.Exists(),
		get:    protocgengo.Install,
	},
	{
		name:   "protoc-gen-go-grpc",
		exists: protocgengogrpc.Exists(),
		get:    protocgengogrpc.Install,
	},
	{
		name:   "swagger",
		exists: goswagger.Exists(),
		get:    goswagger.Install,
	},
}

func check(_ *cobra.Command, _ []string) error {
	return Prepare(boolVarInstall, boolVarForce, boolVarVerbose, boolVarClearCache)
}

func Prepare(install, force, verbose, clear bool) error {
	if clear {
		err := clearCache()
		if err != nil {
			return err
		}
	}

	log := console.NewColorConsole(verbose)
	pending := true
	log.Info("[goctl-env]: preparing to check env")
	defer func() {
		if p := recover(); p != nil {
			log.Error("%+v", p)
			return
		}
		if pending {
			log.Success("\n[goctl-env]: congratulations! your goctl environment is ready!")
		} else {
			log.Error(`
[goctl-env]: check env finish, some dependencies is not found in PATH, you can execute
command 'goctl env check --install' to install it, for details, please execute command 
'goctl env check --help'`)
		}
	}()
	for _, e := range bins {
		time.Sleep(200 * time.Millisecond)
		log.Info("")
		log.Info("[goctl-env]: looking up %q", e.name)
		if e.exists && !clear {
			log.Success("[goctl-env]: %q is installed", e.name)
			continue
		}
		log.Warning("[goctl-env]: %q is not found in PATH", e.name)
		if install {
			install := func() {
				log.Info("[goctl-env]: preparing to install %q", e.name)
				path, err := e.get(env.Get(env.GoctlCache))
				if err != nil {
					log.Error("[goctl-env]: an error interrupted the installation: %+v", err)
					pending = false
				} else {
					log.Success("[goctl-env]: %q is already installed in %q", e.name, path)
				}
			}
			if force {
				install()
				continue
			}
			console.Info("[goctl-env]: do you want to install %q [y: YES, n: No]", e.name)
			for {
				var in string
				fmt.Scanln(&in)
				var brk bool
				switch {
				case strings.EqualFold(in, "y"):
					install()
					brk = true
				case strings.EqualFold(in, "n"):
					pending = false
					console.Info("[goctl-env]: %q installation is ignored", e.name)
					brk = true
				default:
					console.Error("[goctl-env]: invalid input, input 'y' for yes, 'n' for no")
				}
				if brk {
					break
				}
			}
		} else {
			pending = false
		}
	}
	return nil
}

func clearCache() error {
	dir, err := pathx.GetCacheDir()
	if err != nil {
		return err
	}

	files, err := fileutil.ListFileNames(dir)
	if err != nil {
		return err
	}

	for _, v := range files {
		err := fileutil.RemoveFile(filepath.Join(dir, v))
		if err != nil {
			return err
		}
	}

	return nil
}
