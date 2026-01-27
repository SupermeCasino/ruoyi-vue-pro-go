package model

import "github.com/wxlbd/ruoyi-mall-go/pkg/types"

// BitBool is a boolean that maps to BIT(1) in database
type BitBool = types.BitBool

// IntListFromCSV handles comma-separated integer lists from MyBatis IntegerListTypeHandler.
type IntListFromCSV = types.IntListFromCSV

// Int64ListFromCSV handles comma-separated long integer lists.
type Int64ListFromCSV = types.Int64ListFromCSV

// StringListFromCSV handles comma-separated string lists.
type StringListFromCSV = types.StringListFromCSV

// TimeOfDay handles TIME type from database (HH:MM:SS format)
type TimeOfDay = types.TimeOfDay

// FlexInt64 supports deserializing from both number and string
type FlexInt64 = types.FlexInt64

// NewBitBool creates a new BitBool
var NewBitBool = types.NewBitBool
