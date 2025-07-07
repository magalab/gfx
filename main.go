package main

import (
	"log"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/magalab/gfx/gfcmd"
)

func main() {
	var (
		ctx          = gctx.GetInitCtx()
		command, err = gfcmd.GetCommand(ctx)
	)

	if err != nil {
		log.Fatalf(`%+v`, err)
	}
	if command == nil {
		panic(gerror.New(`retrieve root command failed for "gf"`))
	}
	command.Run(ctx)
}
