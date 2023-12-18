package provider

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/taliesins/terraform-provider-hyperv/api"
	hyperv_winrm "github.com/taliesins/terraform-provider-hyperv/api/hyperv-winrm"

	"github.com/dylanmei/iso8601"
	pool "github.com/jolestar/go-commons-pool/v2"
	winrm "github.com/masterzen/winrm"
	winrm_helper "github.com/taliesins/terraform-provider-hyperv/api/winrm-helper"
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

	KrbRealm  string
	KrbSpn    string
	KrbConfig string
	KrbCCache string

	NTLM bool

	TLSServerName string
	CACert        []byte
	Cert          []byte
	Key           []byte

	ScriptPath string
	Timeout    string
}

// HypervWinRmClient() returns a new client for configuring hyperv.
func (c *Config) Client() (comm api.Client, err error) {
	log.Printf("[INFO][hyperv] HyperV HypervWinRmClient configured for HyperV API operations using:\n"+
		"  Host: %s\n"+
		"  Port: %d\n"+
		"  User: %s\n"+
		"  Password: %t\n"+
		"  HTTPS: %t\n"+
		"  Insecure: %t\n"+
		"  NTLM: %t\n"+
		"  KrbRealm: %s\n"+
		"  KrbSpn: %s\n"+
		"  KrbConfig: %s\n"+
		"  KrbCCache: %s\n"+
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
		c.KrbRealm,
		c.KrbSpn,
		c.KrbConfig,
		c.KrbCCache,
		c.TLSServerName,
		c.CACert != nil,
		c.Cert != nil,
		c.Key != nil,
		c.ScriptPath,
		c.Timeout,
	)

	hyperVProvider, err := getHypervProvider(c)

	if err != nil {
		return nil, err
	}

	return hyperVProvider.Client, nil
}

// New creates a new communicator implementation over WinRM.
func GetWinrmClient(config *Config) (winrmClient *winrm.Client, err error) {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	endpoint, err := parseEndpoint(addr, config.HTTPS, config.Insecure, config.TLSServerName, config.CACert, config.Cert, config.Key, config.Timeout)
	if err != nil {
		return nil, err
	}

	params := winrm.DefaultParameters

	if config.KrbRealm != "" {
		proto := "http"
		if config.HTTPS {
			proto = "https"
		}

		params.TransportDecorator = func() winrm.Transporter {
			return &winrm.ClientKerberos{
				Username:  config.User,
				Password:  config.Password,
				Hostname:  config.Host,
				Port:      config.Port,
				Proto:     proto,
				Realm:     config.KrbRealm,
				SPN:       config.KrbSpn,
				KrbConf:   config.KrbConfig,
				KrbCCache: config.KrbCCache,
			}
		}
	} else if config.NTLM {
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

func getHypervProvider(config *Config) (hypervProvider *api.Provider, err error) {
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

	winrmHelperProvider, err := winrm_helper.New(&winrm_helper.ClientConfig{
		WinRmClientPool:  winRmClientPool,
		Vars:             "",
		ElevatedUser:     config.User,
		ElevatedPassword: config.Password,
	})

	if err != nil {
		return nil, err
	}

	return hyperv_winrm.New(&hyperv_winrm.ClientConfig{
		WinRmClient: winrmHelperProvider.Client,
	})
}
