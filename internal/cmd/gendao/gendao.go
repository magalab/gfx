package gendao

import (
	"log"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gtag"
	"github.com/magalab/gfx/internal/consts"
)

const (
	CGenDaoConfig         = `gfcli.gen.dao`
	CGenDaoUsage          = `gfx gen dao [OPTION]`
	CGenDaoBrief          = `gen dao template for startup`
	CGenDaoEg             = `gfx gen dao`
	CGenDaoBriefDstFolder = `destination folder path storing automatically generated go files. default: internal/dao`
	CGenDaoName           = `destination dao name`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`CGenDaoConfig`:         CGenDaoConfig,
		`CGenDaoUsage`:          CGenDaoUsage,
		`CGenDaoBrief`:          CGenDaoBrief,
		`CGenDaoEg`:             CGenDaoEg,
		`CGenDaoBriefDstFolder`: CGenDaoBriefDstFolder,
		`CGenDaoName`:           CGenDaoName,
	})
}

type (
	CGenDao      struct{}
	CGenDaoInput struct {
		g.Meta    `name:"dao" config:"{CGenDaoConfig}" usage:"{CGenDaoUsage}" brief:"{CGenDaoBrief}" eg:"{CGenDaoEg}"`
		DstFolder string `short:"d" name:"dstFolder" brief:"{CGenDaoBriefDstFolder}" d:"internal/dao"`
		DaoName   string `short:"s" name:"daoName" brief:"{CGenDaoName}" v:"required"`
	}
	CGenDaoOutput struct{}
)

func (c *CGenDao) Dao(ctx g.Ctx, in CGenDaoInput) (out *CGenDaoOutput, err error) {
	filePath := gfile.Join(in.DstFolder, in.DaoName+".go")
	if !gfile.Exists(filePath) {
		return nil, gerror.New("you should `gf gen dao` or `make dao` first")
	}
	content := gstr.ReplaceByMap(consts.TemplateGenDao, g.MapStrStr{
		"{Service}": gstr.CaseCamel(in.DaoName),
		"{service}": gstr.CaseCamelLower(in.DaoName),
	})
	if err = gfile.PutContentsAppend(filePath, gstr.TrimLeft(content)); err != nil {
		return nil, err
	}

	log.Println(`done!`)

	return nil, nil

}
