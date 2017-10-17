package hyperv

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
	"github.com/masterzen/winrm"
	"github.com/dylanmei/iso8601"
	"github.com/taliesins/terraform-provider-hyperv/api"
	"strconv"
)

type Config struct {
	User          	string
	Password      	string
	Host  	      	string
	Port	      		int
	HTTPS	      	bool
	Insecure      	bool
	TLSServerName 	string
	CACert     		[]byte
	Key    			[]byte
	Cert     		[]byte
	ScriptPath 		string
	Timeout 			string
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
			"  TLSServerName: %t\n"+
			"  CACert: %t\n"+
			"  Cert: %t\n"+
			"  Key: %t\n"+
			"  ScriptPath: %t\n"+
			"  Timeout: %t",
		c.Host,
		c.Port,
		c.User,
		c.Password != "",
		c.HTTPS,
		c.Insecure,
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
func getWinrmClient(config *Config) (winrmClient *winrm.Client, err error) {
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

	//if config.TransportDecorator != nil {
	//	params.TransportDecorator = config.TransportDecorator
	//}

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
		return nil, fmt.Errorf("Couldn't convert \"\" to an address.")
	}
	if !strings.Contains(addr, ":") || (strings.HasPrefix(addr, "[") && strings.HasSuffix(addr, "]")) {
		host = addr
		port = 5985
	} else {
		shost, sport, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, fmt.Errorf("Couldn't convert \"%s\" to an address.", addr)
		}
		// Check for IPv6 addresses and reformat appropriately
		host = ipFormat(shost)
		port, err = strconv.Atoi(sport)
		if err != nil {
			return nil, fmt.Errorf("Couldn't convert \"%s\" to a port number.", sport)
		}
	}

	timeoutDuration, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, fmt.Errorf("Couldn't convert \"%s\" to a duration.", timeout)
	}

	return &winrm.Endpoint{
		Host:          	host,
		Port:          	port,
		HTTPS:         	https,
		Insecure:      	insecure,
		TLSServerName: 	tlsServerName,
		Cert:			cert,
		Key:				key,
		CACert:        	caCert,
		Timeout:       	timeoutDuration,
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
	winrmClient, err := getWinrmClient(config)

	if err != nil {
		return hypervClient, err
	}

	hypervClient = &api.HypervClient{
		ElevatedPassword: config.Password,
		ElevatedUser:     config.User,
		Vars:             "",
		WinrmClient:      winrmClient,
	}

	return hypervClient, err
}