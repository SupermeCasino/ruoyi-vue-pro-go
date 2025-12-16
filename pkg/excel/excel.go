package excel

import (
	"fmt"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// WriteExcel 写入 Excel 并输出到 HTTP 响应
func WriteExcel(c *gin.Context, fileName string, sheetName string, data interface{}) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 创建 Sheet
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// 写入 header 和 data
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return fmt.Errorf("data must be a slice")
	}

	// 没有任何数据
	if val.Len() == 0 {
		return writeResponse(c, fileName, f)
	}

	// 获取结构体字段
	elemType := val.Index(0).Type()
	// 如果是指针，获取其指向的元素类型
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	var headers []string
	var fields []string

	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		// 优先使用 label 标签，其次 json 标签
		header := field.Tag.Get("label")
		if header == "" {
			continue // 如果没有 label 标签，跳过
		}
		headers = append(headers, header)
		fields = append(fields, field.Name)
	}

	// 写入 Header (第一行)
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// 写入数据
	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)
		// 如果是指针，获取其指向的元素
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}

		for j, fieldName := range fields {
			fieldVal := item.FieldByName(fieldName)
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)

			// 处理时间格式
			if fieldVal.Type() == reflect.TypeOf(time.Time{}) {
				f.SetCellValue(sheetName, cell, fieldVal.Interface().(time.Time).Format("2006-01-02 15:04:05"))
			} else if fieldVal.Type() == reflect.TypeOf(&time.Time{}) {
				if !fieldVal.IsNil() {
					f.SetCellValue(sheetName, cell, fieldVal.Interface().(*time.Time).Format("2006-01-02 15:04:05"))
				}
			} else {
				f.SetCellValue(sheetName, cell, fieldVal.Interface())
			}
		}
	}

	f.SetActiveSheet(index)
	return writeResponse(c, fileName, f)
}

func writeResponse(c *gin.Context, fileName string, f *excelize.File) error {
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	return f.Write(c.Writer)
}
