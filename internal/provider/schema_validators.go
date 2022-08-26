package provider

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func stringKeyInMap(valid interface{}, ignoreCase bool) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		mapType := reflect.ValueOf(valid)
		if mapType.Kind() != reflect.Map {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "not a map!",
			})

			return diags
		}

		mapKeyString, ok := i.(string)
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("expected type of %s to be string", i),
			})

			return diags
		}

		if ignoreCase {
			mapKeyString = strings.ToLower(mapKeyString)
		}

		mapKeyType := reflect.ValueOf(mapKeyString)
		mapValueType := mapType.MapIndex(mapKeyType)

		if !mapValueType.IsValid() {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("expected %s to be one of %v mapKeyString, got %s", i, valid, mapKeyString),
			})

			return diags
		}

		return diags
	}
}

func IntInSlice(valid []int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		value, ok := i.(int)
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("expected type of %s to be int", i),
			})

			return diags
		}

		for _, validValue := range valid {
			if value == validValue {
				return diags
			}
		}

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("expected %s to be one of %v, got %v", i, valid, value),
		})

		return diags
	}
}

func IntBetween(min, max int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		v, ok := i.(int)
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("expected type of %s to be int", i),
			})

			return diags
		}

		if v < min || v > max {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("expected %s to be in the range (%d - %d), got %d", i, min, max, v),
			})
		}

		return diags
	}
}

func ValueOrIntBetween(value, min, max int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		v, ok := i.(int)
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("expected type of %s to be int", i),
			})

			return diags
		}

		if v == value {
			return diags
		}

		if v < min || v > max {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("expected %s to be in the range (%d - %d), got %d", i, min, max, v),
			})

			return diags
		}

		return diags
	}
}
