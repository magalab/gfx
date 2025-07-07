package consts

const TemplateGenDao = `
// GetBy{Service}Id 
func (d *{service}Dao) GetBy{Service}Id(ctx context.Context, req *model.{Service}GetReq) (*model.{Service}Model, error){
	var item *model.{Service}Model
	if err := d.Ctx(ctx).
		Where(d.Columns().{Service}Id, req.{Service}Id).
		Scan(&item); err != nil {
		return nil, err
	}
	return item, nil 
}

// UpdateBy{Service}Id
func (d *{service}Dao) UpdateBy{Service}Id(ctx context.Context, req *model.{Service}UpdateReq) (error){
	updater := g.Map{
		// TODO
	}
	if _,  err := d.Ctx(ctx).
		Where(d.Columns().{Service}Id, req.{Service}Id).
		UpdateAndGetAffected(updater); err != nil {
		return err
	}
	return nil 
}

// DeleteBy{Service}Id
func (d *{service}Dao) DeleteBy{Service}(ctx context.Context, req *model.{Service}DeleteReq) (error){
	updater := g.Map{
		d.Columns().DeletedAt: gtime.Now().Unix(),
	}
	if _,  err := d.Ctx(ctx).
		Where(d.Columns().{Service}Id, req.{Service}Id).
		UpdateAndGetAffected(updater); err != nil {
		return err
	}
	return nil 
}

// {Service}List
func (d *{service}Dao) {Service}List(ctx context.Context, req *model.{Service}ListReq) ([]*model.{Service}Model, int, error){	
	db := d.Ctx(ctx).
		Where(d.Columns().{Service}Id, req.{Service}Id)
	count, err := db.Count()
	if err != nil {
		return nil, 0, err
	}
	if count == 0 {
		return nil, 0, nil
	}
	var items []*model.{Service}Model
	// TODO order
	if err = db.Page(req.PageNum, req.PageSize).
		Scan(&items); err != nil {
		return nil, 0, err
	}

	return items, count, nil
}

`
