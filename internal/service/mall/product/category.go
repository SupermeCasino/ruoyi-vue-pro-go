package product

import (
	"context"

	product "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	productModel "github.com/wxlbd/ruoyi-mall-go/internal/model/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"

	"github.com/samber/lo"
)

type ProductCategoryService struct {
	q *query.Query
}

func NewProductCategoryService(q *query.Query) *ProductCategoryService {
	return &ProductCategoryService{q: q}
}

// CreateCategory 创建商品分类
func (s *ProductCategoryService) CreateCategory(ctx context.Context, req *product.ProductCategoryCreateReq) (int64, error) {
	// 校验父分类
	if err := s.validateParentCategory(ctx, req.ParentID); err != nil {
		return 0, err
	}

	category := &productModel.ProductCategory{
		ParentID: req.ParentID,
		Name:     req.Name,
		PicURL:   req.PicURL,
		Sort:     req.Sort,
		Status:   req.Status,
	}
	err := s.q.ProductCategory.WithContext(ctx).Create(category)
	return category.ID, err
}

// UpdateCategory 更新商品分类
func (s *ProductCategoryService) UpdateCategory(ctx context.Context, req *product.ProductCategoryUpdateReq) error {
	// 校验存在
	if err := s.ValidateCategory(ctx, req.ID); err != nil {
		return err
	}
	// 校验父分类
	if err := s.validateParentCategory(ctx, req.ParentID); err != nil {
		return err
	}
	// 校验不能设置自己为父分类
	if req.ID == req.ParentID {
		return errors.NewBizError(1006001004, "不能设置自己为父分类")
	}

	u := s.q.ProductCategory
	_, err := u.WithContext(ctx).Where(u.ID.Eq(req.ID)).Updates(&productModel.ProductCategory{
		ParentID: req.ParentID,
		Name:     req.Name,
		PicURL:   req.PicURL,
		Sort:     req.Sort,
		Status:   req.Status,
	})
	return err
}

// DeleteCategory 删除商品分类
func (s *ProductCategoryService) DeleteCategory(ctx context.Context, id int64) error {
	// 校验存在
	if err := s.ValidateCategory(ctx, id); err != nil {
		return err
	}
	// 校验是否有子分类
	count, err := s.q.ProductCategory.WithContext(ctx).Where(s.q.ProductCategory.ParentID.Eq(id)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.NewBizError(1006001001, "存在子分类，无法删除")
	}
	// 校验是否绑定了 SPU
	spuCount, err := s.q.ProductSpu.WithContext(ctx).Where(s.q.ProductSpu.CategoryID.Eq(id)).Count()
	if err != nil {
		return err
	}
	if spuCount > 0 {
		return errors.NewBizError(1006001004, "存在商品绑定，无法删除")
	}

	_, err = s.q.ProductCategory.WithContext(ctx).Where(s.q.ProductCategory.ID.Eq(id)).Delete()
	return err
}

// GetCategory 获得商品分类
func (s *ProductCategoryService) GetCategory(ctx context.Context, id int64) (*product.ProductCategoryResp, error) {
	u := s.q.ProductCategory
	category, err := u.WithContext(ctx).Where(u.ID.Eq(id)).First()
	if err != nil {
		return nil, nil
	}
	return s.convertResp(category), nil
}

// GetCategoryList 获得商品分类列表
func (s *ProductCategoryService) GetCategoryList(ctx context.Context, req *product.ProductCategoryListReq) ([]*product.ProductCategoryResp, error) {
	u := s.q.ProductCategory
	q := u.WithContext(ctx)
	if req.Name != "" {
		q = q.Where(u.Name.Like("%" + req.Name + "%"))
	}
	if req.ParentID != nil {
		q = q.Where(u.ParentID.Eq(*req.ParentID))
	}
	if req.Status != nil {
		q = q.Where(u.Status.Eq(*req.Status))
	}
	list, err := q.Order(u.Sort.Asc(), u.ID.Asc()).Find() // Sort asc (Java parity)
	if err != nil {
		return nil, err
	}
	return lo.Map(list, func(item *productModel.ProductCategory, _ int) *product.ProductCategoryResp {
		return s.convertResp(item)
	}), nil
}

// GetEnableCategoryList 获得开启状态的商品分类列表
func (s *ProductCategoryService) GetEnableCategoryList(ctx context.Context) ([]*productModel.ProductCategory, error) {
	u := s.q.ProductCategory
	return u.WithContext(ctx).Where(u.Status.Eq(0)).Order(u.Sort.Asc()).Find()
}

// GetEnableCategoryListByIds 获得开启状态的商品分类列表，指定编号
func (s *ProductCategoryService) GetEnableCategoryListByIds(ctx context.Context, ids []int64) ([]*productModel.ProductCategory, error) {
	if len(ids) == 0 {
		return []*productModel.ProductCategory{}, nil
	}
	u := s.q.ProductCategory
	return u.WithContext(ctx).Where(u.ID.In(ids...), u.Status.Eq(0)).Order(u.Sort.Asc()).Find()
}

// GetCategoryAndChildrenIds 获得分类及其所有子分类的编号
func (s *ProductCategoryService) GetCategoryAndChildrenIds(ctx context.Context, categoryID int64) ([]int64, error) {
	if categoryID == 0 {
		return []int64{}, nil
	}

	categoryIds := []int64{categoryID}

	// 获取该分类的所有子分类
	u := s.q.ProductCategory
	children, err := u.WithContext(ctx).Where(u.ParentID.Eq(categoryID), u.Status.Eq(consts.CommonStatusDisable)).Find()
	if err != nil {
		return nil, err
	}

	for _, child := range children {
		categoryIds = append(categoryIds, child.ID)
	}

	return categoryIds, nil
}

func (s *ProductCategoryService) validateParentCategory(ctx context.Context, parentId int64) error {
	if parentId == 0 {
		return nil
	}
	// 父分类必须存在
	u := s.q.ProductCategory
	parent, err := u.WithContext(ctx).Where(u.ID.Eq(parentId)).First()
	if err != nil {
		return errors.NewBizError(1006001002, "父分类不存在")
	}
	// 父分类不能是二级分类 (即 parentId != 0) -> 意味着只能创建二级分类 (parentId 指向一级), 不能创建三级
	// Logic: If parent's ParentID is NOT 0, it means parent is ALREADY a child (Level 2).
	// So we cannot add a child to it.
	if parent.ParentID != 0 {
		return errors.NewBizError(1006001003, "父分类不能是二级分类")
	}
	return nil
}

func (s *ProductCategoryService) ValidateCategory(ctx context.Context, id int64) error {
	u := s.q.ProductCategory
	count, err := u.WithContext(ctx).Where(u.ID.Eq(id)).Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.NewBizError(1006001000, "分类不存在")
	}
	return nil
}

// ValidateCategoryLevel 校验分类层级（只允许二级分类绑定商品）
func (s *ProductCategoryService) ValidateCategoryLevel(ctx context.Context, id int64) error {
	u := s.q.ProductCategory
	category, err := u.WithContext(ctx).Where(u.ID.Eq(id)).First()
	if err != nil {
		return errors.NewBizError(1006001000, "分类不存在")
	}
	// 只有二级分类（ParentID != 0）才能绑定商品
	if category.ParentID == 0 {
		return errors.NewBizError(1006001005, "只能在二级分类下创建商品")
	}
	return nil
}

func (s *ProductCategoryService) convertResp(item *productModel.ProductCategory) *product.ProductCategoryResp {
	return &product.ProductCategoryResp{
		ID:          item.ID,
		ParentID:    item.ParentID,
		Name:        item.Name,
		PicURL:      item.PicURL,
		Sort:        item.Sort,
		Status:      item.Status,
		Description: item.Description,
		CreateTime:  item.CreateTime,
	}
}
