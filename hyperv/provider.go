package hyperv

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

const (
	DefaultHost = "127.0.0.1"

	DefaultUseHTTPS = false

	DefaultAllowInsecure = false

	DefaultTLSServerName = ""

	// DefaultUser is used if there is no user given
	DefaultUser = "Administrator"

	// DefaultPort is used if there is no port given
	DefaultPort = 5985

	DefaultCACertFile = ""

	DefaultCertFile = ""

	DefaultKeyFile = ""

	// DefaultScriptPath is used as the path to copy the file to
	// for remote execution if not provided otherwise.
	DefaultScriptPath = "C:/Temp/terraform_%RAND%.cmd"

	// DefaultTimeout is used if there is no timeout given
	DefaultTimeoutString = "30s"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_USER", DefaultUser),
				Description: "The user name for HyperV API operations.",
			},

			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_PASSWORD", nil),
				Description: "The user password for HyperV API operations.",
			},

			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_HOST", DefaultHost),
				Description: "The HyperV server host for HyperV API operations.",
			},

			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_PORT", DefaultPort),
				Description: "The HyperV server port for HyperV API operations.",
			},

			"https": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_HTTPS", DefaultUseHTTPS),
				Description: "Should https communication be used for HyperV API operations.",
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_INSECURE", DefaultAllowInsecure),
				Description: "Should insecure communication be used for HyperV API operations.",
			},

			"tls_server_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_TLS_SERVER_NAME", DefaultTLSServerName),
				Description: "Should TLS server name be used for HyperV API operations.",
			},

			"cacert_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_CACERT_PATH", DefaultCACertFile),
				Description: "The ca cert to use for HyperV API operations.",
			},

			"cert_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_CERT_PATH", DefaultCertFile),
				Description: "The cert to use for HyperV API operations.",
			},

			"key_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_KEY_PATH", DefaultKeyFile),
				Description: "The cert key to use for HyperV API operations.",
			},

			"script_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_SCRIPT_PATH", DefaultScriptPath),
				Description: "The script path on host for HyperV API operations.",
			},

			"timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_TIMEOUT", DefaultTimeoutString),
				Description: "Timeout for HyperV API operations.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"hyperv_network_switch":   resourceHyperVNetworkSwitch(),
			"hyperv_machine_instance": resourceHyperVMachineInstance(),
			"hyperv_vhd":              resourceHyperVVhd(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var err error = nil
	var cacert []byte = nil
	cacertPath := d.Get("cacert_path").(string)
	if cacertPath != "" {
		if _, err := os.Stat(cacertPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("cacertPath does not exist - %s", cacertPath)
		}

		cacert, err = ioutil.ReadFile(cacertPath)
		if err != nil {
			return nil, err
		}
	}

	var cert []byte = nil
	certPath := d.Get("cert_path").(string)
	if certPath != "" {
		if _, err := os.Stat(certPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("certPath does not exist - %s", certPath)
		}

		cert, err = ioutil.ReadFile(certPath)
		if err != nil {
			return nil, err
		}
	}

	var key []byte = nil
	keyPath := d.Get("key_path").(string)
	if keyPath != "" {
		if _, err := os.Stat(keyPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("keyPath does not exist - %s", keyPath)
		}

		key, err = ioutil.ReadFile(keyPath)
		if err != nil {
			return nil, err
		}
	}

	config := Config{
		User:          d.Get("user").(string),
		Password:      d.Get("password").(string),
		Host:          d.Get("host").(string),
		Port:          d.Get("port").(int),
		HTTPS:         d.Get("https").(bool),
		CACert:        cacert,
		Cert:          cert,
		Key:           key,
		Insecure:      d.Get("insecure").(bool),
		TLSServerName: d.Get("tls_server_name").(string),
		ScriptPath:    d.Get("script_path").(string),
		Timeout:       d.Get("timeout").(string),
	}

	return config.Client()
}
