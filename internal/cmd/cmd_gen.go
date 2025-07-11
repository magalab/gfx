package cmd

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gtag"
)

var (
	Gen = cGen{}
)

type cGen struct {
	g.Meta `name:"gen" brief:"{cGenBrief}" dc:"{cGenDc}"`
	cGenLogic
	cGenApi
	cGenModel
	cGenDao
}

const (
	cGenBrief = `automatically generate go files for logic`
	cGenDc    = `
The "gen" command is designed for multiple generating purposes. 
It's currently supporting generating go files for ORM models, protobuf and protobuf entity files.
Please use "gfx gen logic -h" for specified type help.
`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`cGenBrief`: cGenBrief,
		`cGenDc`:    cGenDc,
	})
}
