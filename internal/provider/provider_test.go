//go:build integration
// +build integration

package provider

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "0.0.0"

	// goreleaser can also pass the specific commit if you want
	commit string = ""
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"hyperv": func() (*schema.Provider, error) {
		return New(version, commit)(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New(version, commit)().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func escapeForHcl(value string) string {
	return strings.ReplaceAll(value, "\\", "\\\\")
}

func randInt() int {
	rand.Seed(time.Now().UnixNano())
	min := 100
	max := 999
	return rand.Intn(max-min+1) + min
}
