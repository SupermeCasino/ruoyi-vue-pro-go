package promotion

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"testing/quick"
)

// ExpectedConstantLocation defines where specific constants should be located
type ExpectedConstantLocation struct {
	ConstantPrefix string // Prefix of constant names (e.g., "CouponStatus", "ActivityStatus")
	ExpectedFile   string // Expected file name (e.g., "coupon_constants.go")
	Domain         string // Domain description (e.g., "Coupon Status")
}

// getExpectedConstantLocations returns the expected organization of constants
func getExpectedConstantLocations() []ExpectedConstantLocation {
	return []ExpectedConstantLocation{
		// Coupon-related constants should be in coupon_constants.go
		{ConstantPrefix: "CouponStatus", ExpectedFile: "coupon_constants.go", Domain: "Coupon Status"},
		{ConstantPrefix: "CouponTakeType", ExpectedFile: "coupon_constants.go", Domain: "Coupon Take Type"},
		{ConstantPrefix: "CouponValidityType", ExpectedFile: "coupon_constants.go", Domain: "Coupon Validity Type"},

		// Activity-related constants should be in activity_constants.go
		{ConstantPrefix: "ActivityStatus", ExpectedFile: "activity_constants.go", Domain: "Activity Status"},
		{ConstantPrefix: "SeckillActivityStatus", ExpectedFile: "activity_constants.go", Domain: "Seckill Activity Status"},

		// Promotion type constants are in promotion_type.go (existing file)
		{ConstantPrefix: "PromotionType", ExpectedFile: "promotion_type.go", Domain: "Promotion Type"},
		{ConstantPrefix: "CombinationRecordStatus", ExpectedFile: "activity_constants.go", Domain: "Combination Record Status"},
		{ConstantPrefix: "BargainRecordStatus", ExpectedFile: "activity_constants.go", Domain: "Bargain Record Status"},

		// Common promotion constants should be in common_constants.go
		{ConstantPrefix: "ProductScope", ExpectedFile: "common_constants.go", Domain: "Product Scope"},
		{ConstantPrefix: "DiscountType", ExpectedFile: "common_constants.go", Domain: "Discount Type"},
		{ConstantPrefix: "ConditionType", ExpectedFile: "common_constants.go", Domain: "Condition Type"},
		{ConstantPrefix: "CommonStatus", ExpectedFile: "common_constants.go", Domain: "Common Status"},

		// Banner-related constants should be in banner_constants.go
		{ConstantPrefix: "BannerPosition", ExpectedFile: "banner_constants.go", Domain: "Banner Position"},
		{ConstantPrefix: "BannerStatus", ExpectedFile: "banner_constants.go", Domain: "Banner Status"},
		{ConstantPrefix: "BannerPriority", ExpectedFile: "banner_constants.go", Domain: "Banner Priority"},
		{ConstantPrefix: "BannerType", ExpectedFile: "banner_constants.go", Domain: "Banner Type"},
	}
}

// scanPromotionConstants scans all Go files in the promotion package for constant definitions
func scanPromotionConstants() ([]ConstantDefinition, error) {
	var constants []ConstantDefinition

	// Get the current working directory and construct the promotion package path
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	promotionDir := filepath.Join(wd, ".")

	// Walk through all Go files in the promotion directory
	err = filepath.Walk(promotionDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Go files and test files
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Parse the Go file
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		// Extract constants from the AST
		ast.Inspect(node, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.GenDecl:
				if x.Tok == token.CONST {
					for _, spec := range x.Specs {
						if valueSpec, ok := spec.(*ast.ValueSpec); ok {
							for i, name := range valueSpec.Names {
								var value int
								if i < len(valueSpec.Values) {
									if basicLit, ok := valueSpec.Values[i].(*ast.BasicLit); ok {
										// Try to parse as integer
										if parsedVal, parseErr := strconv.Atoi(basicLit.Value); parseErr == nil {
											value = parsedVal
										}
									}
								}

								var comment string
								if x.Doc != nil && len(x.Doc.List) > 0 {
									comment = x.Doc.List[0].Text
								}

								constants = append(constants, ConstantDefinition{
									Name:    name.Name,
									Value:   value,
									File:    filepath.Base(path),
									Package: node.Name.Name,
									Comment: comment,
								})
							}
						}
					}
				}
			}
			return true
		})

		return nil
	})

	return constants, err
}

// TestConstantOrganizationProperty tests Property 2: Constant Organization
// Feature: promotion-constants-refactor, Property 2: For any constant definition in the promotion module, it should be located in the appropriate model layer file
func TestConstantOrganizationProperty(t *testing.T) {
	// Scan all constants in the promotion package
	constants, err := scanPromotionConstants()
	if err != nil {
		t.Fatalf("Failed to scan promotion constants: %v", err)
	}

	expectedLocations := getExpectedConstantLocations()

	// Property-based test: For any constant definition, it should be in the correct file
	property := func() bool {
		for _, constant := range constants {
			// Skip non-promotion constants (like those imported from other packages)
			if !isPromotionConstant(constant.Name) {
				continue
			}

			// Find the expected location for this constant
			expectedFile := findExpectedFile(constant.Name, expectedLocations)
			if expectedFile == "" {
				// If we can't determine expected file, it might be a new constant
				// For now, we'll allow it but log it
				t.Logf("Warning: Unknown constant %s in file %s", constant.Name, constant.File)
				continue
			}

			// Verify the constant is in the expected file
			if constant.File != expectedFile {
				t.Errorf("Constant %s should be in %s but found in %s",
					constant.Name, expectedFile, constant.File)
				return false
			}
		}
		return true
	}

	// Run the property test
	if err := quick.Check(property, &quick.Config{MaxCount: 1}); err != nil {
		t.Errorf("Constant organization property failed: %v", err)
	}
}

// isPromotionConstant checks if a constant name is a promotion-related constant
func isPromotionConstant(name string) bool {
	promotionPrefixes := []string{
		"CouponStatus", "CouponTakeType", "CouponValidityType",
		"ActivityStatus", "PromotionType", "SeckillActivityStatus",
		"CombinationRecordStatus", "BargainRecordStatus",
		"ProductScope", "DiscountType", "ConditionType",
		"BannerPosition", "BannerStatus", "BannerPriority", "BannerType",
		"DefaultPageSize", "MaxPageSize", "MinPrice", "MaxPrice",
		"MinDiscountPercent", "MaxDiscountPercent", "CouponTemplateTakeLimitCountMax",
	}

	for _, prefix := range promotionPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}

	return false
}

// findExpectedFile finds the expected file for a given constant name
func findExpectedFile(constantName string, expectedLocations []ExpectedConstantLocation) string {
	for _, location := range expectedLocations {
		if strings.HasPrefix(constantName, location.ConstantPrefix) {
			return location.ExpectedFile
		}
	}

	// Handle special cases for general constants
	generalConstants := []string{
		"DefaultPageSize", "MaxPageSize", "MinPrice", "MaxPrice",
		"MinDiscountPercent", "MaxDiscountPercent", "CouponTemplateTakeLimitCountMax",
	}

	for _, general := range generalConstants {
		if constantName == general {
			return "common_constants.go"
		}
	}

	return ""
}

// TestConstantOrganizationSpecificCases tests specific organization requirements
func TestConstantOrganizationSpecificCases(t *testing.T) {
	constants, err := scanPromotionConstants()
	if err != nil {
		t.Fatalf("Failed to scan promotion constants: %v", err)
	}

	// Test specific requirements from the spec
	testCases := []struct {
		constantName string
		expectedFile string
		requirement  string
	}{
		// Requirement 1.1: Coupon status constants in model layer
		{"CouponStatusUnused", "coupon_constants.go", "1.1"},
		{"CouponStatusUsed", "coupon_constants.go", "1.1"},
		{"CouponStatusExpired", "coupon_constants.go", "1.1"},

		// Requirement 2.1: Coupon take type constants in model layer
		{"CouponTakeTypeUser", "coupon_constants.go", "2.1"},
		{"CouponTakeTypeAdmin", "coupon_constants.go", "2.1"},
		{"CouponTakeTypeRegister", "coupon_constants.go", "2.1"},

		// Requirement 3.1: Product scope constants in model layer
		{"ProductScopeAll", "common_constants.go", "3.1"},
		{"ProductScopeSpu", "common_constants.go", "3.1"},
		{"ProductScopeCategory", "common_constants.go", "3.1"},

		// Requirement 4.1: Discount type constants in model layer
		{"DiscountTypePrice", "common_constants.go", "4.1"},
		{"DiscountTypePercent", "common_constants.go", "4.1"},

		// Requirement 5.1: Template validity type constants in model layer
		{"CouponValidityTypeDate", "coupon_constants.go", "5.1"},
		{"CouponValidityTypeTerm", "coupon_constants.go", "5.1"},

		// Requirement 6.1: Activity status constants in model layer
		{"ActivityStatusWait", "activity_constants.go", "6.1"},
		{"ActivityStatusRun", "activity_constants.go", "6.1"},
		{"ActivityStatusEnd", "activity_constants.go", "6.1"},
		{"ActivityStatusClose", "activity_constants.go", "6.1"},

		// Requirement 7.1: Common status constants in model layer
		// Note: These use model.CommonStatus* so they're referenced, not defined here

		// Requirement 8.1: Banner position constants in model layer
		{"BannerPositionHome", "banner_constants.go", "8.1"},
		{"BannerPositionSeckill", "banner_constants.go", "8.1"},
		{"BannerPositionCombination", "banner_constants.go", "8.1"},

		// Requirement 10.1: Constants organized in model layer
		{"ConditionTypePrice", "common_constants.go", "10.1"},
		{"ConditionTypeCount", "common_constants.go", "10.1"},
	}

	for _, tc := range testCases {
		t.Run(tc.constantName, func(t *testing.T) {
			found := false
			for _, constant := range constants {
				if constant.Name == tc.constantName {
					found = true
					if constant.File != tc.expectedFile {
						t.Errorf("Requirement %s: Constant %s should be in %s but found in %s",
							tc.requirement, tc.constantName, tc.expectedFile, constant.File)
					}

					// Verify it's in the promotion package (model layer)
					if constant.Package != "promotion" {
						t.Errorf("Requirement %s: Constant %s should be in promotion package but found in %s",
							tc.requirement, tc.constantName, constant.Package)
					}
					break
				}
			}

			if !found {
				t.Errorf("Requirement %s: Constant %s not found in promotion package", tc.requirement, tc.constantName)
			}
		})
	}
}

// TestModelLayerLocation verifies all constants are in the model layer
func TestModelLayerLocation(t *testing.T) {
	constants, err := scanPromotionConstants()
	if err != nil {
		t.Fatalf("Failed to scan promotion constants: %v", err)
	}

	for _, constant := range constants {
		if isPromotionConstant(constant.Name) {
			// Verify it's in the promotion package (which is part of the model layer)
			if constant.Package != "promotion" {
				t.Errorf("Constant %s should be in promotion package (model layer) but found in %s",
					constant.Name, constant.Package)
			}

			// Verify it's in a constants file
			if !strings.HasSuffix(constant.File, "_constants.go") &&
				constant.File != "promotion_type.go" { // promotion_type.go is an exception
				t.Errorf("Constant %s should be in a constants file but found in %s",
					constant.Name, constant.File)
			}
		}
	}
}
