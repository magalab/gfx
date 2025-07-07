package genmodel

import (
	"log"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gtag"
	"github.com/magalab/gfx/internal/consts"
)

const (
	CGenModelConfig         = `gfcli.gen.model`
	CGenModelUsage          = `gfx gen model [OPTION]`
	CGenModelBrief          = `gen model template for startup`
	CGenModelEg             = `gfx gen model`
	CGenModelBriefDstFolder = `destination folder path storing automatically generated go files. default: internal/model`
	CGenModelName           = `destination model name`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`CGenModelConfig`:         CGenModelConfig,
		`CGenModelUsage`:          CGenModelUsage,
		`CGenModelBrief`:          CGenModelBrief,
		`CGenModelEg`:             CGenModelEg,
		`CGenModelBriefDstFolder`: CGenModelBriefDstFolder,
		`CGenModelName`:           CGenModelName,
	})
}

type (
	CGenModel      struct{}
	CGenModelInput struct {
		g.Meta    `name:"model" config:"{CGenModelConfig}" usage:"{CGenModelUsage}" brief:"{CGenModelBrief}" eg:"{CGenModelEg}"`
		DstFolder string `short:"d" name:"dstFolder" brief:"{CGenModelBriefDstFolder}" d:"internal/model"`
		ModelName string `short:"s" name:"modelName" brief:"{CGenModelName}" v:"required"`
	}
	CGenModelOutput struct{}
)

func (c *CGenModel) Model(ctx g.Ctx, in CGenModelInput) (out *CGenModelOutput, err error) {
	pwd := gfile.Pwd()
	goModPath := gfile.Join(pwd, "go.mod")
	importPath := ""
	if gfile.Exists(goModPath) {
		match, _ := gregex.MatchString(`^module\s+(.+)\s*`, gfile.GetContents(goModPath))
		importPath = gstr.Trim(match[1]) + "/internal/model/entity"
	} else {
		return nil, gerror.New("go.mod file not found")
	}

	filePath := gfile.Join(in.DstFolder, in.ModelName+".go")

	if !gfile.Exists(filePath) {
		content := gstr.ReplaceByMap(consts.TemplateGenModel, g.MapStrStr{
			"{ModelName}":  gstr.CaseCamel(in.ModelName),
			"{ImportPath}": importPath,
		})
		if err = gfile.PutContents(filePath, gstr.TrimLeft(content)); err != nil {
			return nil, err
		}
	}

	log.Println(`done!`)

	return nil, nil

}
