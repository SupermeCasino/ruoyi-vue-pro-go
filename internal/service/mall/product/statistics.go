package product

import (
	"context"
	"fmt"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// ProductStatisticsService 商品统计服务接口
type ProductStatisticsService interface {
	// GetProductStatisticsRankPage 获得商品统计排行榜分页
	GetProductStatisticsRankPage(ctx context.Context, reqVO *req.ProductStatisticsReqVO, pageParam *pagination.PageParam) (*pagination.PageResult[interface{}], error)
	// GetProductStatisticsAnalyse 获得商品统计分析（含环比对照）
	GetProductStatisticsAnalyse(ctx context.Context, reqVO *req.ProductStatisticsReqVO) (*resp.DataComparisonRespVO[resp.ProductStatisticsRespVO], error)
	// GetProductStatisticsList 获得商品统计列表
	GetProductStatisticsList(ctx context.Context, reqVO *req.ProductStatisticsReqVO) ([]*resp.ProductStatisticsRespVO, error)
	// StatisticsProduct 统计指定天数的商品数据
	StatisticsProduct(ctx context.Context, days int) (string, error)
}

// ProductStatisticsRepository 商品统计数据访问接口
type ProductStatisticsRepository interface {
	// GetByDateRange 根据日期范围获取统计数据
	GetByDateRange(ctx context.Context, beginTime, endTime time.Time) ([]*resp.ProductStatisticsRespVO, error)
	// GetSummaryByDateRange 根据日期范围获取汇总统计数据
	GetSummaryByDateRange(ctx context.Context, beginTime, endTime time.Time) (*resp.ProductStatisticsRespVO, error)
	// GetPageGroupBySpuId 分页获取按 SPU 分组的统计数据
	GetPageGroupBySpuId(ctx context.Context, reqVO *req.ProductStatisticsReqVO, pageParam *pagination.PageParam) (*pagination.PageResult[*resp.ProductStatisticsRespVO], error)
	// CountByDateRange 统计指定日期范围内的记录数
	CountByDateRange(ctx context.Context, beginTime, endTime time.Time) (int64, error)
	// StatisticsProductByDateRange 统计指定日期范围内的商品数据并入库
	StatisticsProductByDateRange(ctx context.Context, date time.Time, beginTime, endTime time.Time) error
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
// 对应 Java: ProductStatisticsServiceImpl.getProductStatisticsRankPage
func (s *ProductStatisticsServiceImpl) GetProductStatisticsRankPage(ctx context.Context, reqVO *req.ProductStatisticsReqVO, pageParam *pagination.PageParam) (*pagination.PageResult[interface{}], error) {
	// 调用仓储层分页查询（按 SPU 分组聚合）
	pageResult, err := s.productStatisticsRepo.GetPageGroupBySpuId(ctx, reqVO, pageParam)
	if err != nil {
		return nil, err
	}

	// 转换为 interface{} slice
	var resultList []interface{}
	for _, item := range pageResult.List {
		resultList = append(resultList, item)
	}

	return &pagination.PageResult[interface{}]{
		List:  resultList,
		Total: pageResult.Total,
	}, nil
}

// GetProductStatisticsAnalyse 获得商品统计分析
// 对应 Java: ProductStatisticsServiceImpl.getProductStatisticsAnalyse
func (s *ProductStatisticsServiceImpl) GetProductStatisticsAnalyse(ctx context.Context, reqVO *req.ProductStatisticsReqVO) (*resp.DataComparisonRespVO[resp.ProductStatisticsRespVO], error) {
	beginTime := reqVO.Times[0]
	endTime := reqVO.Times[1]

	// 1. 统计数据：查询当前时间范围的汇总数据
	value, err := s.productStatisticsRepo.GetSummaryByDateRange(ctx, beginTime, endTime)
	if err != nil {
		return nil, err
	}
	if value == nil {
		value = &resp.ProductStatisticsRespVO{}
	}

	// 2. 对照数据：环比，时长一致
	duration := endTime.Sub(beginTime)
	referenceBeginTime := beginTime.Add(-duration)
	referenceEndTime := beginTime

	reference, err := s.productStatisticsRepo.GetSummaryByDateRange(ctx, referenceBeginTime, referenceEndTime)
	if err != nil {
		return nil, err
	}
	if reference == nil {
		reference = &resp.ProductStatisticsRespVO{}
	}

	return &resp.DataComparisonRespVO[resp.ProductStatisticsRespVO]{
		Summary:    value,
		Comparison: reference,
	}, nil
}

// GetProductStatisticsList 获得商品统计列表
// 对应 Java: ProductStatisticsServiceImpl.getProductStatisticsList
func (s *ProductStatisticsServiceImpl) GetProductStatisticsList(ctx context.Context, reqVO *req.ProductStatisticsReqVO) ([]*resp.ProductStatisticsRespVO, error) {
	return s.productStatisticsRepo.GetByDateRange(ctx, reqVO.Times[0], reqVO.Times[1])
}

// StatisticsProduct 统计指定天数的商品数据
// 对应 Java: ProductStatisticsServiceImpl.statisticsProduct
func (s *ProductStatisticsServiceImpl) StatisticsProduct(ctx context.Context, days int) (string, error) {
	today := time.Now()
	var results []string

	// 遍历指定天数，逐天统计
	for day := 1; day <= days; day++ {
		date := today.AddDate(0, 0, -day)
		result, err := s.statisticsProductByDate(ctx, date)
		if err != nil {
			return "", err
		}
		results = append(results, result)
	}

	// 合并结果
	var output string
	for _, r := range results {
		if output != "" {
			output += "\n"
		}
		output += r
	}
	return output, nil
}

// statisticsProductByDate 统计指定日期的商品数据
// 对应 Java: ProductStatisticsServiceImpl.statisticsProduct(LocalDateTime date)
func (s *ProductStatisticsServiceImpl) statisticsProductByDate(ctx context.Context, date time.Time) (string, error) {
	// 1. 处理统计时间范围（当天的开始和结束）
	beginTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endTime := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())
	dateStr := date.Format("2006-01-02")

	// 2. 检查该日是否已经统计过
	count, err := s.productStatisticsRepo.CountByDateRange(ctx, beginTime, endTime)
	if err != nil {
		return "", err
	}
	if count > 0 {
		return fmt.Sprintf("%s 数据已存在，如果需要重新统计，请先删除对应的数据", dateStr), nil
	}

	// 3. 执行统计并入库
	startTime := time.Now()
	err = s.productStatisticsRepo.StatisticsProductByDateRange(ctx, date, beginTime, endTime)
	if err != nil {
		return "", err
	}
	elapsed := time.Since(startTime)

	return fmt.Sprintf("%s 统计完成，耗时 %v", dateStr, elapsed), nil
}
