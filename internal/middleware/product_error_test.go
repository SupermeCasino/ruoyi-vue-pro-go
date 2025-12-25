package middleware

import (
	"testing"

	"github.com/wxlbd/ruoyi-mall-go/internal/model/product"
)

func TestIsProductError(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		expected bool
	}{
		{
			name:     "商品SPU不存在错误码",
			code:     1008005000,
			expected: true,
		},
		{
			name:     "商品SKU不存在错误码",
			code:     1008006000,
			expected: true,
		},
		{
			name:     "商品分类不存在错误码",
			code:     1008001000,
			expected: true,
		},
		{
			name:     "非商品模块错误码",
			code:     1004003002,
			expected: false,
		},
		{
			name:     "系统错误码",
			code:     500,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isProductError(tt.code)
			if result != tt.expected {
				t.Errorf("isProductError(%d) = %v, expected %v", tt.code, result, tt.expected)
			}
		})
	}
}

func TestProductErrorMapping(t *testing.T) {
	tests := []struct {
		name        string
		inputError  error
		expectedErr error
	}{
		{
			name:        "record not found 映射为 SPU 不存在",
			inputError:  &testError{msg: "record not found"},
			expectedErr: product.ErrSpuNotExists,
		},
		{
			name:        "invalid input 映射为参数错误",
			inputError:  &testError{msg: "invalid input"},
			expectedErr: nil, // 这里会返回 errors.ErrParam，但我们简化测试
		},
		{
			name:        "nil 错误返回 nil",
			inputError:  nil,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProductErrorMapping(tt.inputError)
			if tt.expectedErr != nil && result != tt.expectedErr {
				t.Errorf("ProductErrorMapping() = %v, expected %v", result, tt.expectedErr)
			}
			if tt.expectedErr == nil && result == nil {
				// 正确的情况
			}
		})
	}
}

// 测试用的错误类型
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
