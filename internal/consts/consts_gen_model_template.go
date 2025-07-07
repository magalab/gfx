package consts

const TemplateGenModel = `
package model 

import (
	"{ImportPath}"
)

type {ModelName}Model struct {
	entity.{ModelName}
}

type {ModelName}AddReq struct {
  
}

type {ModelName}DeleteReq struct {
  
}

type {ModelName}UpdateReq struct {
  
}

type {ModelName}ListReq struct {
  
}

type {ModelName}GetReq struct {
  
}

type {ModelName}Item struct {

}

func (m *{ModelName}Model) ToItem() *{ModelName}Item {

    return &{ModelName}Item{
        
    }
}

`
