package provider

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
The string is limited to 15 characters.

For ISO 9660 discs, the volume name can use the following characters:

"A" through "Z"
"0" through "9"
"_" (underscore)

For Joliet and UDF discs, the volume name can use the following characters:

"a" through "z"
"A" through "Z"
"0" through "9"
"." (period)
"_" (underscore)
*/
func AllowedIsoVolumeName() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		validIso9660VolumeNameRegex := regexp.MustCompile(`^[A-Z0-9_]*$`)
		// validJolietAndUdfVolumeNameRegex := regexp.MustCompile(`^[a-zA-Z0-9_\.]*$`)

		value, ok := v.(string)
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("expected type of %s to be string", v),
			})

			return diags
		}

		if len(value) > 15 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("expected length of %s to be 15 characters or less", v),
			})

			return diags
		}

		if !validIso9660VolumeNameRegex.MatchString(value) {
			diags = append(diags, diag.Errorf("%q must only use characters `A`-`Z`, `0`-`9` or `_`", value)...)
		}

		return diags
	}
}

func StringKeyInMap(valid interface{}, ignoreCase bool) schema.SchemaValidateDiagFunc {
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

func IsDivisibleBy(logicalSize int) schema.SchemaValidateDiagFunc {
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

		if v%logicalSize != 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("expected %d to be perfectly divisible by %d, maybe use %d instead", v, logicalSize, int(math.Ceil(float64(v)/float64(logicalSize)))*logicalSize),
			})
		}

		return diags
	}
}
