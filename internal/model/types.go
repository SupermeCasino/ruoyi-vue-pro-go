package model

import "github.com/wxlbd/ruoyi-mall-go/pkg/types"

// BitBool is a boolean that maps to BIT(1) in database
type BitBool = types.BitBool

// IntListFromCSV handles comma-separated integer lists from MyBatis IntegerListTypeHandler.
type IntListFromCSV = types.IntListFromCSV

// StringListFromCSV handles comma-separated string lists.
type StringListFromCSV = types.StringListFromCSV

// NewBitBool creates a new BitBool
var NewBitBool = types.NewBitBool
