package service

import (
	"context"
	"errors"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
)

type DictService struct {
	q *query.Query
}

func NewDictService(q *query.Query) *DictService {
	return &DictService{
		q: q,
	}
}

// --- DictType ---

func (s *DictService) CreateDictType(ctx context.Context, req *req.DictTypeSaveReq) (int64, error) {
	t := s.q.SystemDictType
	// Check if type exists
	count, err := t.WithContext(ctx).Where(t.Type.Eq(req.Type)).Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("字典类型已存在")
	}

	dictType := &model.SystemDictType{
		Name:   req.Name,
		Type:   req.Type,
		Status: int32(req.Status),
		Remark: req.Remark,
	}
	err = t.WithContext(ctx).Create(dictType)
	return dictType.ID, err
}

func (s *DictService) UpdateDictType(ctx context.Context, req *req.DictTypeSaveReq) error {
	t := s.q.SystemDictType
	// Check existence
	_, err := t.WithContext(ctx).Where(t.ID.Eq(req.ID)).First()
	if err != nil {
		return errors.New("字典类型不存在")
	}
	// Check if type conflict (if type changed)
	count, err := t.WithContext(ctx).Where(t.Type.Eq(req.Type), t.ID.Neq(req.ID)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("字典类型已存在")
	}

	_, err = t.WithContext(ctx).Where(t.ID.Eq(req.ID)).Updates(&model.SystemDictType{
		Name:   req.Name,
		Type:   req.Type,
		Status: int32(req.Status),
		Remark: req.Remark,
	})
	return err
}

func (s *DictService) DeleteDictType(ctx context.Context, id int64) error {
	t := s.q.SystemDictType
	_, err := t.WithContext(ctx).Where(t.ID.Eq(id)).Delete()
	return err
}

func (s *DictService) GetDictType(ctx context.Context, id int64) (*resp.DictTypeRespVO, error) {
	t := s.q.SystemDictType
	item, err := t.WithContext(ctx).Where(t.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return &resp.DictTypeRespVO{
		ID:         item.ID,
		Name:       item.Name,
		Type:       item.Type,
		Status:     item.Status,
		Remark:     item.Remark,
		CreateTime: item.CreatedAt,
	}, nil
}

func (s *DictService) GetDictTypePage(ctx context.Context, req *req.DictTypePageReq) (*core.PageResult[*resp.DictTypeRespVO], error) {
	t := s.q.SystemDictType
	qb := t.WithContext(ctx)

	if req.Name != "" {
		qb = qb.Where(t.Name.Like("%" + req.Name + "%"))
	}
	if req.Type != "" {
		qb = qb.Where(t.Type.Like("%" + req.Type + "%"))
	}
	if req.Status != nil {
		qb = qb.Where(t.Status.Eq(int32(*req.Status)))
	}

	total, err := qb.Count()
	if err != nil {
		return nil, err
	}

	list, err := qb.Order(t.ID.Desc()).Offset(req.GetOffset()).Limit(req.PageSize).Find()
	if err != nil {
		return nil, err
	}

	var data []*resp.DictTypeRespVO
	for _, item := range list {
		data = append(data, &resp.DictTypeRespVO{
			ID:         item.ID,
			Name:       item.Name,
			Type:       item.Type,
			Status:     item.Status,
			Remark:     item.Remark,
			CreateTime: item.CreatedAt,
		})
	}

	return &core.PageResult[*resp.DictTypeRespVO]{
		List:  data,
		Total: total,
	}, nil
}

func (s *DictService) GetSimpleDictTypeList(ctx context.Context) ([]*resp.DictTypeSimpleRespVO, error) {
	t := s.q.SystemDictType
	list, err := t.WithContext(ctx).Order(t.ID).Find() // Return all, frontend filters? Or Java returns all. Java returns all actually.
	if err != nil {
		return nil, err
	}

	var res []*resp.DictTypeSimpleRespVO
	for _, item := range list {
		res = append(res, &resp.DictTypeSimpleRespVO{
			ID:   item.ID,
			Type: item.Type,
			Name: item.Name,
		})
	}
	return res, nil
}

// --- DictData ---

func (s *DictService) CreateDictData(ctx context.Context, req *req.DictDataSaveReq) (int64, error) {
	d := s.q.SystemDictData
	dictData := &model.SystemDictData{
		Sort:      req.Sort,
		Label:     req.Label,
		Value:     req.Value,
		DictType:  req.DictType,
		Status:    int32(req.Status),
		ColorType: req.ColorType,
		CssClass:  req.CssClass,
		Remark:    req.Remark,
	}
	err := d.WithContext(ctx).Create(dictData)
	return dictData.ID, err
}

func (s *DictService) UpdateDictData(ctx context.Context, req *req.DictDataSaveReq) error {
	d := s.q.SystemDictData
	_, err := d.WithContext(ctx).Where(d.ID.Eq(req.ID)).Updates(&model.SystemDictData{
		Sort:      req.Sort,
		Label:     req.Label,
		Value:     req.Value,
		DictType:  req.DictType,
		Status:    int32(req.Status),
		ColorType: req.ColorType,
		CssClass:  req.CssClass,
		Remark:    req.Remark,
	})
	return err
}

func (s *DictService) DeleteDictData(ctx context.Context, id int64) error {
	d := s.q.SystemDictData
	_, err := d.WithContext(ctx).Where(d.ID.Eq(id)).Delete()
	return err
}

func (s *DictService) GetDictData(ctx context.Context, id int64) (*resp.DictDataRespVO, error) {
	d := s.q.SystemDictData
	item, err := d.WithContext(ctx).Where(d.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return &resp.DictDataRespVO{
		ID:         item.ID,
		Sort:       item.Sort,
		Label:      item.Label,
		Value:      item.Value,
		DictType:   item.DictType,
		Status:     item.Status,
		ColorType:  item.ColorType,
		CssClass:   item.CssClass,
		Remark:     item.Remark,
		CreateTime: item.CreatedAt,
	}, nil
}

func (s *DictService) GetDictDataPage(ctx context.Context, req *req.DictDataPageReq) (*core.PageResult[*resp.DictDataRespVO], error) {
	d := s.q.SystemDictData
	qb := d.WithContext(ctx)

	if req.Label != "" {
		qb = qb.Where(d.Label.Like("%" + req.Label + "%"))
	}
	if req.DictType != "" {
		qb = qb.Where(d.DictType.Eq(req.DictType))
	}
	if req.Status != nil {
		qb = qb.Where(d.Status.Eq(int32(*req.Status)))
	}

	total, err := qb.Count()
	if err != nil {
		return nil, err
	}

	list, err := qb.Order(d.DictType, d.Sort).Offset(req.GetOffset()).Limit(req.PageSize).Find()
	if err != nil {
		return nil, err
	}

	var data []*resp.DictDataRespVO
	for _, item := range list {
		data = append(data, &resp.DictDataRespVO{
			ID:         item.ID,
			Sort:       item.Sort,
			Label:      item.Label,
			Value:      item.Value,
			DictType:   item.DictType,
			Status:     item.Status,
			ColorType:  item.ColorType,
			CssClass:   item.CssClass,
			Remark:     item.Remark,
			CreateTime: item.CreatedAt,
		})
	}

	return &core.PageResult[*resp.DictDataRespVO]{
		List:  data,
		Total: total,
	}, nil
}

func (s *DictService) GetSimpleDictDataList(ctx context.Context) ([]*resp.DictDataSimpleRespVO, error) {
	d := s.q.SystemDictData
	list, err := d.WithContext(ctx).Where(d.Status.Eq(0)).Order(d.DictType, d.Sort).Find()
	if err != nil {
		return nil, err
	}

	var res []*resp.DictDataSimpleRespVO
	for _, item := range list {
		res = append(res, &resp.DictDataSimpleRespVO{
			DictType:  item.DictType,
			Value:     item.Value,
			Label:     item.Label,
			ColorType: item.ColorType,
			CssClass:  item.CssClass,
		})
	}
	return res, nil
}
