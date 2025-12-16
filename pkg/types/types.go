package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// BitBool is a boolean that maps to BIT(1) in database
type BitBool bool

// Scan implements the Scanner interface.
func (b *BitBool) Scan(value interface{}) error {
	if value == nil {
		*b = false
		return nil
	}

	switch v := value.(type) {
	case []uint8:
		if len(v) > 0 {
			*b = BitBool(v[0] == 1)
		} else {
			*b = false
		}
	case int64:
		*b = BitBool(v == 1)
	case bool:
		*b = BitBool(v)
	default:
		return errors.New("incompatible type for BitBool")
	}
	return nil
}

// Value implements the driver Valuer interface.
func (b BitBool) Value() (driver.Value, error) {
	if b {
		return int64(1), nil
	}
	return int64(0), nil
}

// QueryClauses 实现 GORM 软删除查询子句
func (BitBool) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{BitBoolQueryClause{Field: f}}
}

// DeleteClauses 实现 GORM 软删除删除子句
func (BitBool) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{BitBoolDeleteClause{Field: f}}
}

// UpdateClauses 实现 GORM 软删除更新子句
func (BitBool) UpdateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{BitBoolUpdateClause{Field: f}}
}

// BitBoolQueryClause 查询子句 - 自动过滤已删除记录
type BitBoolQueryClause struct {
	Field *schema.Field
}

func (sd BitBoolQueryClause) Name() string {
	return ""
}

func (sd BitBoolQueryClause) Build(clause.Builder) {
}

func (sd BitBoolQueryClause) MergeClause(*clause.Clause) {
}

func (sd BitBoolQueryClause) ModifyStatement(stmt *gorm.Statement) {
	if _, ok := stmt.Clauses["soft_delete_enabled"]; !ok && !stmt.Statement.Unscoped {
		// 检查字段是否启用了软删除 (通过 softDelete tag)
		if sd.Field != nil && sd.Field.TagSettings != nil {
			if _, hasSoftDelete := sd.Field.TagSettings["SOFTDELETE"]; hasSoftDelete {
				// 添加 WHERE deleted = 0 条件
				stmt.AddClause(clause.Where{Exprs: []clause.Expression{
					clause.Eq{
						Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName},
						Value:  false, // BitBool false = 0
					},
				}})
				stmt.Clauses["soft_delete_enabled"] = clause.Clause{}
			}
		}
	}
}

// BitBoolDeleteClause 删除子句 - 软删除时设置 deleted = 1
type BitBoolDeleteClause struct {
	Field *schema.Field
}

func (sd BitBoolDeleteClause) Name() string {
	return ""
}

func (sd BitBoolDeleteClause) Build(clause.Builder) {
}

func (sd BitBoolDeleteClause) MergeClause(*clause.Clause) {
}

func (sd BitBoolDeleteClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.Len() == 0 && !stmt.Statement.Unscoped {
		// 设置 deleted = 1
		set := clause.Set{{
			Column: clause.Column{Name: sd.Field.DBName},
			Value:  true, // BitBool true = 1
		}}
		stmt.SetColumn(sd.Field.DBName, true, true)
		stmt.AddClause(set)

		// 添加 WHERE 条件（基于主键）
		if stmt.Schema != nil {
			_, queryValues := schema.GetIdentityFieldValuesMap(stmt.Context, stmt.ReflectValue, stmt.Schema.PrimaryFields)
			column, values := schema.ToQueryValues(stmt.Table, stmt.Schema.PrimaryFieldDBNames, queryValues)

			if len(values) > 0 {
				stmt.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
			}
		}

		// 添加已删除过滤
		BitBoolQueryClause(sd).ModifyStatement(stmt)

		stmt.AddClauseIfNotExists(clause.Update{})
		stmt.Build(stmt.DB.Callback().Update().Clauses...)
	}
}

// BitBoolUpdateClause 更新子句 - 更新时自动过滤已删除记录
type BitBoolUpdateClause struct {
	Field *schema.Field
}

func (sd BitBoolUpdateClause) Name() string {
	return ""
}

func (sd BitBoolUpdateClause) Build(clause.Builder) {
}

func (sd BitBoolUpdateClause) MergeClause(*clause.Clause) {
}

func (sd BitBoolUpdateClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.Len() == 0 && !stmt.Statement.Unscoped {
		BitBoolQueryClause(sd).ModifyStatement(stmt)
	}
}

func NewBitBool(b bool) BitBool {
	return BitBool(b)
}

// IntListFromCSV handles comma-separated integer lists from MyBatis IntegerListTypeHandler.
// Supports both "1,2,3" format and JSON "[1,2,3]" format.
type IntListFromCSV []int

func (l *IntListFromCSV) Scan(value interface{}) error {
	if value == nil {
		*l = nil
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return errors.New("incompatible type for IntListFromCSV")
	}

	if len(data) == 0 {
		*l = nil
		return nil
	}

	str := strings.TrimSpace(string(data))
	if str == "" {
		*l = nil
		return nil
	}

	// Try JSON format first
	if strings.HasPrefix(str, "[") {
		var result []int
		if err := json.Unmarshal(data, &result); err == nil {
			*l = result
			return nil
		}
	}

	// Parse as comma-separated
	parts := strings.Split(str, ",")
	result := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		i, err := strconv.Atoi(p)
		if err != nil {
			return err
		}
		result = append(result, i)
	}
	*l = result
	return nil
}

func (l IntListFromCSV) Value() (driver.Value, error) {
	if len(l) == 0 {
		return "", nil
	}
	parts := make([]string, len(l))
	for i, v := range l {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, ","), nil
}

func (l IntListFromCSV) MarshalJSON() ([]byte, error) {
	return json.Marshal([]int(l))
}

// StringListFromCSV handles comma-separated string lists.
type StringListFromCSV []string

func (l *StringListFromCSV) Scan(value interface{}) error {
	if value == nil {
		*l = nil
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return errors.New("incompatible type for StringListFromCSV")
	}

	if len(data) == 0 {
		*l = nil
		return nil
	}

	str := strings.TrimSpace(string(data))
	if str == "" {
		*l = nil
		return nil
	}

	// Try JSON format first
	if strings.HasPrefix(str, "[") {
		var result []string
		if err := json.Unmarshal(data, &result); err == nil {
			*l = result
			return nil
		}
	}

	// Parse as comma-separated
	parts := strings.Split(str, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	*l = result
	return nil
}

func (l StringListFromCSV) Value() (driver.Value, error) {
	if len(l) == 0 {
		return "", nil
	}
	return strings.Join(l, ","), nil
}

func (l StringListFromCSV) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(l))
}
