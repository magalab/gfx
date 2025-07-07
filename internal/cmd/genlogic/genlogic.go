package genlogic

import (
	"fmt"
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
	CGenLogicConfig         = `gfcli.gen.logic`
	CGenLogicUsage          = `gfx gen logic [OPTION]`
	CGenLogicBrief          = `gen logic template for startup`
	CGenLogicEg             = `gfx gen logic`
	CGenLogicBriefDstFolder = `destination folder path storing automatically generated go files. default: internal/logic`
	CGenLogicServiceName    = `destination service name`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`CGenLogicConfig`:         CGenLogicConfig,
		`CGenLogicUsage`:          CGenLogicUsage,
		`CGenLogicBrief`:          CGenLogicBrief,
		`CGenLogicEg`:             CGenLogicEg,
		`CGenLogicBriefDstFolder`: CGenLogicBriefDstFolder,
		`CGenLogicServiceName`:    CGenLogicServiceName,
	})
}

type (
	CGenLogic      struct{}
	CGenLogicInput struct {
		g.Meta      `name:"logic" config:"{CGenLogicConfig}" usage:"{CGenLogicUsage}" brief:"{CGenLogicBrief}" eg:"{CGenLogicEg}"`
		DstFolder   string `short:"d" name:"dstFolder" brief:"{CGenLogicBriefDstFolder}" d:"internal/logic"`
		ServiceName string `short:"s" name:"serviceName" brief:"{CGenLogicServiceName}"`
	}
	CGenLogicOutput struct{}
)

func (c *CGenLogic) Logic(ctx g.Ctx, in CGenLogicInput) (out *CGenLogicOutput, err error) {
	pwd := gfile.Pwd()
	goModPath := gfile.Join(pwd, "go.mod")
	importPath := ""
	modelImportPath := ""
	if gfile.Exists(goModPath) {
		match, _ := gregex.MatchString(`^module\s+(.+)\s*`, gfile.GetContents(goModPath))
		importPath = gstr.Trim(match[1]) + "/internal/service"
		modelImportPath = gstr.Trim(match[1]) + "/internal/model"
	} else {
		return nil, gerror.New("go.mod file not found")
	}

	dstLogicFolderPath := gfile.Join(in.DstFolder, in.ServiceName)

	serviceFilePath := gfile.Join(dstLogicFolderPath, in.ServiceName+".go")
	crud := []string{"add", "get", "list", "update", "delete"}

	if !gfile.Exists(serviceFilePath) {
		content := gstr.ReplaceByMap(consts.TemplateGenLogicNew, g.MapStrStr{
			"{PackageName}": in.ServiceName,
			"{Service}":     gstr.CaseCamel(in.ServiceName),

			"{ImportPath}":      importPath,
			"{ModelImportPath}": modelImportPath,
		})
		if err = gfile.PutContents(serviceFilePath, gstr.TrimLeft(content)); err != nil {
			return nil, err
		}
	}
	for _, v := range crud {
		returnType := "error"
		returnRes := "gerror.NewCode(gcode.CodeNotImplemented)"
		importPath := modelImportPath
		if v == "get" {
			returnType = fmt.Sprintf("*model.%sItem, error", gstr.CaseCamel(in.ServiceName))
			returnRes = "nil, gerror.NewCode(gcode.CodeNotImplemented)"
		}
		if v == "list" {
			returnType = fmt.Sprintf("[]*model.%sItem, error", gstr.CaseCamel(in.ServiceName))
			returnRes = "nil, gerror.NewCode(gcode.CodeNotImplemented)"
		}
		servicePath := gfile.Join(dstLogicFolderPath, in.ServiceName+fmt.Sprintf("_%s.go", v))
		content := gstr.ReplaceByMap(consts.TemplateGenLogicNewPlaceholder, g.MapStrStr{
			"{PackageName}": in.ServiceName,
			"{Service}":     gstr.CaseCamel(in.ServiceName),
			"{Method}":      gstr.CaseCamel(v),
			"{ReturnType}":  returnType,
			"{ReturnRes}":   returnRes,
			"{ImportPath}":  importPath,
		})
		if err = gfile.PutContents(servicePath, gstr.TrimLeft(content)); err != nil {
			return nil, err
		}
	}

	log.Println(`done!`)

	return nil, nil

}
