package gfcmd

import (
	"context"
	"log"
	"runtime"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/magalab/gfx/internal/cmd"
)

const cliFolderName = `hack`

type Command struct {
	*gcmd.Command
}

// Run starts running the command according the command line arguments and options.
func (c *Command) Run(ctx context.Context) {
	defer func() {
		if exception := recover(); exception != nil {
			if err, ok := exception.(error); ok {
				log.Print(err.Error())
			} else {
				panic(gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception))
			}
		}
	}()

	// CLI configuration, using the `hack/config.yaml` in priority.
	if path, _ := gfile.Search(cliFolderName); path != "" {
		if adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile); ok {
			if err := adapter.SetPath(path); err != nil {
				log.Fatal(err)
			}
		}
	}

	// zsh alias "git fetch" conflicts checks.
	handleZshAlias()

	// -y option checks.

	// just run.
	if err := c.RunWithError(ctx); err != nil {
		// Exit with error message and exit code 1.
		// It is very important to exit the command process with code 1.
		log.Fatalf(`%+v`, err)
	}
}

// GetCommand retrieves and returns the root command of CLI `gf`.
func GetCommand(ctx context.Context) (*Command, error) {
	root, err := gcmd.NewFromObject(cmd.GFX)
	if err != nil {
		return nil, err
	}
	err = root.AddObject(
		cmd.Gen,
	)
	if err != nil {
		return nil, err
	}
	command := &Command{
		root,
	}
	return command, nil
}

// zsh alias "git fetch" conflicts checks.
func handleZshAlias() {
	if runtime.GOOS == "windows" {
		return
	}
	if home, err := gfile.Home(); err == nil {
		zshPath := gfile.Join(home, ".zshrc")
		if gfile.Exists(zshPath) {
			var (
				aliasCommand = `alias gfx=gfx`
				content      = gfile.GetContents(zshPath)
			)
			if !gstr.Contains(content, aliasCommand) {
				_ = gfile.PutContentsAppend(zshPath, "\n"+aliasCommand+"\n")
			}
		}
	}
}
