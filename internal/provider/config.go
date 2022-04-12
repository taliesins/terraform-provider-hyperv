package provider

import (
	"context"
	"fmt"
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/dylanmei/iso8601"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	pool "github.com/jolestar/go-commons-pool/v2"
	"github.com/masterzen/winrm"
	"github.com/taliesins/terraform-provider-hyperv/api"
)

type Config struct {
	Version          string
	Commit           string
	TerraformVersion string
	User             string
	Password         string
	Host             string
	Port             int
	HTTPS            bool
	Insecure         bool
	NTLM             bool
	TLSServerName    string
	CACert           []byte
	Key              []byte
	Cert             []byte
	ScriptPath       string
	Timeout          string
}

// HypervClient() returns a new client for configuring hyperv.
func (c *Config) Client() (comm *api.HypervClient, err error) {
	log.Printf("[INFO][hyperv] HyperV HypervClient configured for HyperV API operations using:\n"+
		"  Host: %s\n"+
		"  Port: %d\n"+
		"  User: %s\n"+
		"  Password: %t\n"+
		"  HTTPS: %t\n"+
		"  Insecure: %t\n"+
		"  NTLM: %t\n"+
		"  TLSServerName: %s\n"+
		"  CACert: %t\n"+
		"  Cert: %t\n"+
		"  Key: %t\n"+
		"  ScriptPath: %s\n"+
		"  Timeout: %s",
		c.Host,
		c.Port,
		c.User,
		c.Password != "",
		c.HTTPS,
		c.Insecure,
		c.NTLM,
		c.TLSServerName,
		c.CACert != nil,
		c.Cert != nil,
		c.Key != nil,
		c.ScriptPath,
		c.Timeout,
	)

	return getHypervClient(c)
}

// New creates a new communicator implementation over WinRM.
func GetWinrmClient(config *Config) (winrmClient *winrm.Client, err error) {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	endpoint, err := parseEndpoint(addr, config.HTTPS, config.Insecure, config.TLSServerName, config.CACert, config.Cert, config.Key, config.Timeout)
	if err != nil {
		return nil, err
	}

	params := winrm.NewParameters(
		winrm.DefaultParameters.Timeout,
		winrm.DefaultParameters.Locale,
		winrm.DefaultParameters.EnvelopeSize,
	)

	if config.NTLM {
		params.TransportDecorator = func() winrm.Transporter { return &winrm.ClientNTLM{} }
	}

	if endpoint.Timeout.Seconds() > 0 {
		params.Timeout = iso8601.FormatDuration(endpoint.Timeout)
	}

	winrmClient, err = winrm.NewClientWithParameters(
		endpoint, config.User, config.Password, params)

	if err != nil {
		return nil, err
	}

	return winrmClient, nil
}

func parseEndpoint(addr string, https bool, insecure bool, tlsServerName string, caCert []byte, cert []byte, key []byte, timeout string) (*winrm.Endpoint, error) {
	var host string
	var port int

	if addr == "" {
		return nil, fmt.Errorf("couldn't convert \"\" to an address")
	}
	if !strings.Contains(addr, ":") || (strings.HasPrefix(addr, "[") && strings.HasSuffix(addr, "]")) {
		host = addr
		port = 5985
	} else {
		shost, sport, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, fmt.Errorf("couldn't convert \"%s\" to an address", addr)
		}
		// Check for IPv6 addresses and reformat appropriately
		host = ipFormat(shost)
		port, err = strconv.Atoi(sport)
		if err != nil {
			return nil, fmt.Errorf("couldn't convert \"%s\" to a port number", sport)
		}
	}

	timeoutDuration, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, fmt.Errorf("couldn't convert \"%s\" to a duration", timeout)
	}

	return &winrm.Endpoint{
		Host:          host,
		Port:          port,
		HTTPS:         https,
		Insecure:      insecure,
		TLSServerName: tlsServerName,
		Cert:          cert,
		Key:           key,
		CACert:        caCert,
		Timeout:       timeoutDuration,
	}, nil
}

// ipFormat formats the IP correctly, so we don't provide IPv6 address in an IPv4 format during node communication.
// We return the ip parameter as is if it's an IPv4 address or a hostname.
func ipFormat(ip string) string {
	ipObj := net.ParseIP(ip)
	// Return the ip/host as is if it's either a hostname or an IPv4 address.
	if ipObj == nil || ipObj.To4() != nil {
		return ip
	}

	return fmt.Sprintf("[%s]", ip)
}

func getHypervClient(config *Config) (hypervClient *api.HypervClient, err error) {
	ctx := context.Background()
	factory := pool.NewPooledObjectFactorySimple(
		func(context.Context) (interface{}, error) {
			winrmClient, err := GetWinrmClient(config)

			if err != nil {
				return nil, err
			}

			return winrmClient, nil
		})

	winRmClientPool := pool.NewObjectPoolWithDefaultConfig(ctx, factory)
	winRmClientPool.Config.BlockWhenExhausted = true
	winRmClientPool.Config.MinIdle = 0
	winRmClientPool.Config.MaxIdle = 2
	winRmClientPool.Config.MaxTotal = 5
	winRmClientPool.Config.TimeBetweenEvictionRuns = 10 * time.Second

	hypervClient = &api.HypervClient{
		WinRmClientPool:  winRmClientPool,
		Vars:             "",
		ElevatedUser:     config.User,
		ElevatedPassword: config.Password,
	}

	return hypervClient, err
}

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
