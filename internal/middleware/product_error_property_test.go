package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wxlbd/ruoyi-mall-go/internal/model/product"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/logger"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// 初始化测试环境
func init() {
	// 初始化测试用的简单logger
	logger.Log = zap.NewNop() // 使用空logger避免日志输出
}

// TestProductErrorResponseFormat 属性测试：验证错误响应格式一致性
// 属性4: 错误响应格式一致性 - 验证需求 2.5, 7.1, 7.2
func TestProductErrorResponseFormat(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)

	// 测试用例：验证所有商品模块错误码都返回正确的响应格式
	testCases := []struct {
		name        string
		bizError    *errors.BizError
		expectedMsg string
	}{
		{
			name:        "商品SPU不存在错误",
			bizError:    product.ErrSpuNotExists,
			expectedMsg: "商品 SPU 不存在",
		},
		{
			name:        "商品SKU不存在错误",
			bizError:    product.ErrSkuNotExists,
			expectedMsg: "商品 SKU 不存在",
		},
		{
			name:        "商品分类不存在错误",
			bizError:    product.ErrCategoryNotExists,
			expectedMsg: "商品分类不存在",
		},
		{
			name:        "商品品牌不存在错误",
			bizError:    product.ErrBrandNotExists,
			expectedMsg: "品牌不存在",
		},
		{
			name:        "商品SKU库存不足错误",
			bizError:    product.ErrSkuStockNotEnough,
			expectedMsg: "商品 SKU 库存不足",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建测试路由
			router := gin.New()
			router.Use(ProductErrorHandler())

			// 创建测试端点，模拟抛出商品模块错误
			router.GET("/test", func(c *gin.Context) {
				c.Error(tc.bizError)
				c.Abort()
			})

			// 发送测试请求
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 验证响应状态码
			assert.Equal(t, http.StatusOK, w.Code, "响应状态码应该是200")

			// 解析响应体
			var resp response.Result[any]
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err, "响应体应该是有效的JSON")

			// 验证错误响应格式与Java版本一致
			assert.Equal(t, tc.bizError.Code, resp.Code, "错误码应该与预期一致")
			assert.Equal(t, tc.expectedMsg, resp.Msg, "错误信息应该与预期一致")
			assert.Nil(t, resp.Data, "错误响应的data字段应该为null")

			// 验证响应头
			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		})
	}
}

// TestProductErrorCodeRange 属性测试：验证商品模块错误码范围识别
func TestProductErrorCodeRange(t *testing.T) {
	// 测试商品模块错误码范围识别的正确性
	testCases := []struct {
		name     string
		code     int
		expected bool
	}{
		// 商品模块错误码（1008000000-1008999999）
		{"商品模块最小错误码", 1008000000, true},
		{"商品分类错误码", 1008001000, true},
		{"商品品牌错误码", 1008002000, true},
		{"商品属性错误码", 1008003000, true},
		{"商品SPU错误码", 1008005000, true},
		{"商品SKU错误码", 1008006000, true},
		{"商品模块最大错误码", 1008999999, true},

		// 非商品模块错误码
		{"系统错误码", 500, false},
		{"参数错误码", 400, false},
		{"用户模块错误码", 1001001000, false},
		{"订单模块错误码", 1004001000, false},
		{"支付模块错误码", 1007001000, false},
		{"超出商品模块范围", 1009000000, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isProductError(tc.code)
			assert.Equal(t, tc.expected, result,
				"错误码 %d 的识别结果应该是 %v", tc.code, tc.expected)
		})
	}
}

// TestHandleProductErrorConsistency 属性测试：验证错误处理一致性
func TestHandleProductErrorConsistency(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 测试不同类型错误的处理一致性
	testCases := []struct {
		name           string
		inputError     error
		expectedCode   int
		expectedStatus int
	}{
		{
			name:           "商品模块业务错误",
			inputError:     product.ErrSpuNotExists,
			expectedCode:   1008005000,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "系统参数错误",
			inputError:     errors.ErrParam,
			expectedCode:   errors.ParamErrCode,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "系统服务器错误",
			inputError:     fmt.Errorf("系统内部错误"),
			expectedCode:   errors.ServerErrCode,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建测试上下文
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)

			// 调用错误处理函数
			HandleProductError(c, tc.inputError)

			// 验证响应状态码
			assert.Equal(t, tc.expectedStatus, w.Code)

			// 解析响应体
			var resp response.Result[any]
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			// 验证错误码
			assert.Equal(t, tc.expectedCode, resp.Code)
			assert.NotEmpty(t, resp.Msg, "错误信息不应该为空")
		})
	}
}

// TestProductNotFoundErrorMapping 属性测试：验证商品不存在错误映射
func TestProductNotFoundErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 测试不同商品类型的不存在错误映射
	testCases := []struct {
		name         string
		productType  string
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "SPU不存在",
			productType:  "spu",
			expectedCode: 1008005000,
			expectedMsg:  "商品 SPU 不存在",
		},
		{
			name:         "SKU不存在",
			productType:  "sku",
			expectedCode: 1008006000,
			expectedMsg:  "商品 SKU 不存在",
		},
		{
			name:         "分类不存在",
			productType:  "category",
			expectedCode: 1008001000,
			expectedMsg:  "商品分类不存在",
		},
		{
			name:         "品牌不存在",
			productType:  "brand",
			expectedCode: 1008002000,
			expectedMsg:  "品牌不存在",
		},
		{
			name:         "未知类型默认为SPU",
			productType:  "unknown",
			expectedCode: 1008005000,
			expectedMsg:  "商品 SPU 不存在",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建测试上下文
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)

			// 调用商品不存在错误处理函数
			ProductNotFoundError(c, tc.productType, 123)

			// 验证响应
			assert.Equal(t, http.StatusOK, w.Code)

			var resp response.Result[any]
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedCode, resp.Code)
			assert.Equal(t, tc.expectedMsg, resp.Msg)
		})
	}
}

// TestValidateProductParamsConsistency 属性测试：验证参数验证一致性
func TestValidateProductParamsConsistency(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 测试参数验证的一致性
	testCases := []struct {
		name        string
		inputError  error
		expectValid bool
	}{
		{
			name:        "无错误应该返回true",
			inputError:  nil,
			expectValid: true,
		},
		{
			name:        "有错误应该返回false",
			inputError:  fmt.Errorf("参数错误"),
			expectValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建测试上下文
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)

			// 调用参数验证函数
			result := ValidateProductParams(c, tc.inputError)

			// 验证返回值
			assert.Equal(t, tc.expectValid, result)

			if !tc.expectValid {
				// 验证错误响应
				assert.Equal(t, http.StatusOK, w.Code)

				var resp response.Result[any]
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)

				assert.Equal(t, errors.ParamErrCode, resp.Code)
				assert.Contains(t, resp.Msg, "参数错误")
			}
		})
	}
}

// BenchmarkProductErrorHandler 性能基准测试：验证错误处理性能
func BenchmarkProductErrorHandler(b *testing.B) {
	gin.SetMode(gin.TestMode)

	// 创建测试路由
	router := gin.New()
	router.Use(ProductErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		c.Error(product.ErrSpuNotExists)
		c.Abort()
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
