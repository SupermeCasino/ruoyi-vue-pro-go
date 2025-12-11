package service

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/pkg/core"
	"context"
	"time"
)

// ProductStatisticsService 商品统计服务接口
type ProductStatisticsService interface {
	GetProductStatisticsRankPage(ctx context.Context, reqVO *req.ProductStatisticsReqVO, pageParam *core.PageParam) (*core.PageResult[interface{}], error)
	GetProductStatisticsAnalyse(ctx context.Context, reqVO *req.ProductStatisticsReqVO) (*resp.DataComparisonRespVO[resp.ProductStatisticsRespVO], error)
	GetProductStatisticsList(ctx context.Context, reqVO *req.ProductStatisticsReqVO) ([]*resp.ProductStatisticsRespVO, error)
	StatisticsProduct(ctx context.Context, days int) (string, error)
}

// ProductStatisticsRepository 商品统计数据访问接口
type ProductStatisticsRepository interface {
	GetByDateRange(ctx context.Context, beginTime, endTime time.Time) ([]*resp.ProductStatisticsRespVO, error)
}

// ProductStatisticsServiceImpl 商品统计服务实现
type ProductStatisticsServiceImpl struct {
	productStatisticsRepo ProductStatisticsRepository
}

// NewProductStatisticsService 创建商品统计服务
func NewProductStatisticsService(repo ProductStatisticsRepository) ProductStatisticsService {
	return &ProductStatisticsServiceImpl{
		productStatisticsRepo: repo,
	}
}

// GetProductStatisticsRankPage 获得商品统计排行榜分页
func (s *ProductStatisticsServiceImpl) GetProductStatisticsRankPage(ctx context.Context, reqVO *req.ProductStatisticsReqVO, pageParam *core.PageParam) (*core.PageResult[interface{}], error) {
	// 获取所有数据（排序交给前端或后续处理，这里先获取范围内的数据并聚合）
	// 注意：这里假设 Repo GetByDateRange 返回的是聚合后的数据，或者我们需要在内存中聚合
	// 如果 Repo 只是返回明细，我们需要自己聚合。
	// 根据 Java 逻辑，它是查 DB 聚合。我们先假设 Repo 提供了聚合查询，或者我们先查出来再手动分页
	list, err := s.productStatisticsRepo.GetByDateRange(ctx, reqVO.Times[0], reqVO.Times[1])
	if err != nil {
		return nil, err
	}

	// 内存分页
	total := int64(len(list))
	start := (pageParam.PageNo - 1) * pageParam.PageSize
	end := start + pageParam.PageSize
	if start > int(total) {
		start = int(total)
	}
	if end > int(total) {
		end = int(total)
	}

	// 转换为 interface{} slice
	var resultList []interface{}
	for _, item := range list[start:end] {
		resultList = append(resultList, item)
	}

	return &core.PageResult[interface{}]{
		List:  resultList,
		Total: total,
	}, nil
}

// GetProductStatisticsAnalyse 获得商品统计分析
func (s *ProductStatisticsServiceImpl) GetProductStatisticsAnalyse(ctx context.Context, reqVO *req.ProductStatisticsReqVO) (*resp.DataComparisonRespVO[resp.ProductStatisticsRespVO], error) {
	// 1. 查询当前时间范围的数据
	list, err := s.productStatisticsRepo.GetByDateRange(ctx, reqVO.Times[0], reqVO.Times[1])
	if err != nil {
		return nil, err
	}

	// 聚合数据
	summary := &resp.ProductStatisticsRespVO{}
	for _, item := range list {
		summary.BuyCount += item.BuyCount
		summary.BuyPrice += item.BuyPrice
		summary.BrowseCount += item.BrowseCount
		summary.FavoriteCount += item.FavoriteCount
		summary.CommentCount += item.CommentCount
	}

	// 2. 查询对比时间范围的数据 (环比，时长一致)
	duration := reqVO.Times[1].Sub(reqVO.Times[0])
	compareBeginTime := reqVO.Times[0].Add(-duration)
	compareEndTime := reqVO.Times[0]
	compareList, err := s.productStatisticsRepo.GetByDateRange(ctx, compareBeginTime, compareEndTime)
	if err != nil {
		return nil, err
	}

	// 聚合对比数据
	compareSummary := &resp.ProductStatisticsRespVO{}
	for _, item := range compareList {
		compareSummary.BuyCount += item.BuyCount
		compareSummary.BuyPrice += item.BuyPrice
		compareSummary.BrowseCount += item.BrowseCount
		compareSummary.FavoriteCount += item.FavoriteCount
		compareSummary.CommentCount += item.CommentCount
	}

	return &resp.DataComparisonRespVO[resp.ProductStatisticsRespVO]{
		Summary:    summary,
		Comparison: compareSummary,
	}, nil
}

// GetProductStatisticsList 获得商品统计列表
func (s *ProductStatisticsServiceImpl) GetProductStatisticsList(ctx context.Context, reqVO *req.ProductStatisticsReqVO) ([]*resp.ProductStatisticsRespVO, error) {
	return s.productStatisticsRepo.GetByDateRange(ctx, reqVO.Times[0], reqVO.Times[1])
}

// StatisticsProduct 统计指定天数的商品数据
func (s *ProductStatisticsServiceImpl) StatisticsProduct(ctx context.Context, days int) (string, error) {
	// TODO: 实现每日统计逻辑
	return "success", nil
}
