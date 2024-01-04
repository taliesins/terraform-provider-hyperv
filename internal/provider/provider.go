package provider

import (
	"fmt"
	"io/ioutil"
	"os"

	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	DefaultHost = "127.0.0.1"

	DefaultUseHTTPS = true

	DefaultAllowInsecure = false

	DefaultAllowNTLM = true

	DefaultKerberosRealm = ""

	DefaultKerberosServicePrincipalName = ""

	DefaultKerberosConfig = "/etc/krb5.conf"

	DefaultKerberosCredentialCache = ""

	DefaultTLSServerName = ""

	// DefaultUser is used if there is no user given
	DefaultUser = "Administrator"

	// DefaultPort is used if there is no port given
	DefaultPort = 5986

	DefaultCACertFile = ""

	DefaultCertFile = ""

	DefaultKeyFile = ""

	// DefaultScriptPath is used as the path to copy the file to
	// for remote execution if not provided otherwise.
	DefaultScriptPath = "C:/Temp/terraform_%RAND%.cmd"

	// DefaultTimeout is used if there is no timeout given
	DefaultTimeoutString = "30s"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

// Provider returns a terraform.ResourceProvider.
func New(version string, commit string) func() *schema.Provider {
	return func() *schema.Provider {
		provider := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"user": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("HYPERV_USER", DefaultUser),
					Description: "The username to use when HyperV api calls are made. Generally this is Administrator. It can also be sourced from the `HYPERV_USER` environment variable otherwise defaults to `Administrator.",
				},

				"password": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("HYPERV_PASSWORD", ""),
					Description: "The password associated with the username to use for HyperV api calls. It can also be sourced from the `HYPERV_PASSWORD` environment variable`.",
				},

				"host": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("HYPERV_HOST", DefaultHost),
					Description: "The host to run HyperV api calls against. It can also be sourced from the `HYPERV_HOST` environment variable otherwise defaults to `127.0.0.1`.",
				},

				"port": {
					Type:        schema.TypeInt,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("HYPERV_PORT", DefaultPort),
					Description: "The port to run HyperV api calls against. It can also be sourced from the `HYPERV_PORT` environment variable otherwise defaults to `5986`.",
				},

				"https": {
					Type:        schema.TypeBool,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("HYPERV_HTTPS", DefaultUseHTTPS),
					Description: "Should https be used for HyperV api calls. It can also be sourced from `HYPERV_HTTPS` environment variable otherwise defaults to `true`.",
				},

				"insecure": {
					Type:        schema.TypeBool,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("HYPERV_INSECURE", DefaultAllowInsecure),
					Description: "Skips TLS Verification for HyperV api calls. Generally this is used for self-signed certificates. Should only be used if absolutely needed. Can also be set via setting the `HYPERV_INSECURE` environment variable to `true` otherwise defaults to `false`.",
				},

				"use_ntlm": {
					Type:        schema.TypeBool,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("HYPERV_USE_NTLM", DefaultAllowNTLM),
					Description: "Use NTLM for authentication for HyperV api calls. Can also be set via setting the `HYPERV_USE_NTLM` environment variable to `true` otherwise defaults to `true`.",
				},

				"kerberos_realm": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("HYPERV_KERBEROS_REALM", DefaultKerberosRealm),
					Description: "Use Kerberos Realm for authentication for HyperV api calls. Can also be set via setting the `HYPERV_KERBEROS_REALM` environment variable otherwise defaults to empty string.",
				},

				"kerberos_service_principal_name": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("HYPERV_KERBEROS_SERVICE_PRINCIPAL_NAME", DefaultKerberosServicePrincipalName),
					Description: "Use Kerberos Service Principal Name for authentication for HyperV api calls. Can also be set via setting the `HYPERV_KERBEROS_SERVICE_PRINCIPAL_NAME` environment variable otherwise defaults to empty string.",
				},

				"kerberos_config": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{"HYPERV_KERBEROS_CONFIG", "KRB5_CONFIG"}, DefaultKerberosConfig),
					Description: "Use Kerberos Config for authentication for HyperV api calls. Can also be set via setting the `HYPERV_KERBEROS_CONFIG` or `KRB5_CONFIG` environment variable otherwise defaults to `/etc/krb5.conf`.",
				},

				"kerberos_credential_cache": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{"HYPERV_KERBEROS_CREDENTIAL_CACHE", "KRB5CCNAME"}, DefaultKerberosCredentialCache),
					Description: "Use Kerberos Credential Cache for authentication for HyperV api calls. Can also be set via setting the `HYPERV_KERBEROS_CREDENTIAL_CACHE` or `KRB5CCNAME` environment variable otherwise defaults to empty string.",
				},

				"tls_server_name": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("HYPERV_TLS_SERVER_NAME", DefaultTLSServerName),
					Description: "The TLS server name for the host used for HyperV api calls. It can also be sourced from the `HYPERV_TLS_SERVER_NAME` environment variable otherwise defaults to empty string.",
				},

				"cacert_path": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("HYPERV_CACERT_PATH", DefaultCACertFile),
					Description: "The path to the ca certificates to use for HyperV api calls. Can also be sourced from the `HYPERV_CACERT_PATH` environment variable otherwise defaults to empty string.",
				},

				"cert_path": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("HYPERV_CERT_PATH", DefaultCertFile),
					Description: "The path to the certificate to use for authentication for HyperV api calls. Can also be sourced from the `HYPERV_CERT_PATH` environment variable otherwise defaults to empty string.",
				},

				"key_path": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("HYPERV_KEY_PATH", DefaultKeyFile),
					Description: "The path to the certificate private key to use for authentication for HyperV api calls. Can also be sourced from the `HYPERV_KEY_PATH` environment variable otherwise defaults to empty string.",
				},

				"script_path": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("HYPERV_SCRIPT_PATH", DefaultScriptPath),
					Description: "The path used to copy scripts meant for remote execution for HyperV api calls. Can also be sourced from the `HYPERV_SCRIPT_PATH` environment variable otherwise defaults to `C:/Temp/terraform_%RAND%.cmd`.",
				},

				"timeout": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("HYPERV_TIMEOUT", DefaultTimeoutString),
					Description: "The timeout to wait for the connection to become available for HyperV api calls. Should be provided as a string like 30s or 5m. Can also be sourced from the `HYPERV_TIMEOUT` environment variable otherwise defaults to `30s`.",
				},
			},

			ResourcesMap: map[string]*schema.Resource{
				"hyperv_network_switch":   resourceHyperVNetworkSwitch(),
				"hyperv_machine_instance": resourceHyperVMachineInstance(),
				"hyperv_vhd":              resourceHyperVVhd(),
				"hyperv_iso_image":        resourceHyperVIsoImage(),
			},
			DataSourcesMap: map[string]*schema.Resource{
				"hyperv_network_switch":   dataSourceHyperVNetworkSwitch(),
				"hyperv_machine_instance": dataSourceHyperVMachineInstance(),
				"hyperv_vhd":              dataSourceHyperVVhd(),
			},
		}

		provider.ConfigureContextFunc = configure(version, commit, provider)

		return provider
	}
}

func configure(version string, commit string, provider *schema.Provider) func(context context.Context, resourceData *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(context context.Context, resourceData *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics
		var err error
		var cacert []byte = nil
		cacertPath := resourceData.Get("cacert_path").(string)
		if cacertPath != "" {
			if _, err := os.Stat(cacertPath); os.IsNotExist(err) {
				return nil, diag.FromErr(fmt.Errorf("cacertPath does not exist - %s", cacertPath))
			}

			cacert, err = ioutil.ReadFile(cacertPath)
			if err != nil {
				return nil, diag.FromErr(err)
			}
		}

		var cert []byte = nil
		certPath := resourceData.Get("cert_path").(string)
		if certPath != "" {
			if _, err := os.Stat(certPath); os.IsNotExist(err) {
				return nil, diag.FromErr(fmt.Errorf("certPath does not exist - %s", certPath))
			}

			cert, err = ioutil.ReadFile(certPath)
			if err != nil {
				return nil, diag.FromErr(err)
			}
		}

		var key []byte = nil
		keyPath := resourceData.Get("key_path").(string)
		if keyPath != "" {
			if _, err := os.Stat(keyPath); os.IsNotExist(err) {
				return nil, diag.FromErr(fmt.Errorf("keyPath does not exist - %s", keyPath))
			}

			key, err = ioutil.ReadFile(keyPath)
			if err != nil {
				return nil, diag.FromErr(err)
			}
		}

		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}

		config := Config{
			Version:          version,
			Commit:           commit,
			TerraformVersion: terraformVersion,
			User:             resourceData.Get("user").(string),
			Password:         resourceData.Get("password").(string),
			Host:             resourceData.Get("host").(string),
			Port:             resourceData.Get("port").(int),
			HTTPS:            resourceData.Get("https").(bool),
			CACert:           cacert,
			Cert:             cert,
			Key:              key,
			Insecure:         resourceData.Get("insecure").(bool),
			NTLM:             resourceData.Get("use_ntlm").(bool),
			KrbRealm:         resourceData.Get("kerberos_realm").(string),
			KrbSpn:           resourceData.Get("kerberos_service_principal_name").(string),
			KrbConfig:        resourceData.Get("kerberos_config").(string),
			KrbCCache:        resourceData.Get("kerberos_credential_cache").(string),
			TLSServerName:    resourceData.Get("tls_server_name").(string),
			ScriptPath:       resourceData.Get("script_path").(string),
			Timeout:          resourceData.Get("timeout").(string),
		}

		client, err := config.Client()
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return client, diags
	}
}
