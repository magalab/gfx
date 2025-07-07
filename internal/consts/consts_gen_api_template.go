package consts

const TemplateGenApiNew = `
package {ApiVersion}

import (

	"github.com/gogf/gf/v2/frame/g"

	"{ImportPath}"
)


`

const TemplateGenApiNewPlaceholder = `
type {Service}{Method}Req struct {
  	g.Meta ` + "`path:\"/{Path}\" tags:\"{Service}\" method:\"{method}\" sm:\"{sm}\"`" + `
  	*model.{Service}{Method}Req
}
type {Service}{Method}Res struct {
	{Data}
}

`
