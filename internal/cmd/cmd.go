package cmd

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gtag"
)

var GFX = cGFX{}

type cGFX struct {
	g.Meta `name:"gfx" ad:"{cGFXAd}"`
}

const (
	cGFXAd = `
ADDITIONAL
    Use "gfx COMMAND -h" for details about a command.
`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`cGFXAd`: cGFXAd,
	})
}

type cGFInput struct {
	g.Meta  `name:"gfx"`
	Yes     bool `short:"y" name:"yes"     brief:"all yes for all command without prompt ask"   orphan:"true"`
	Version bool `short:"v" name:"version" brief:"show version information of current binary"   orphan:"true"`
	Debug   bool `short:"d" name:"debug"   brief:"show internal detailed debugging information" orphan:"true"`
}

type cGFOutput struct{}
