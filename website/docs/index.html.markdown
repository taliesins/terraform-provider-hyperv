---
layout: "hyperv"
page_title: "Provider: HyperV"
description: |-
  The HyperV provider is used to interact with the many resources supported by HyperV. The provider needs to be configured with credentials for the HyperV API.
---

# HyperV Provider

The HyperV provider is used to interact with the many resources supported by HyperV. The provider needs to be configured with credentials for HyperV API.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure HyperV
provider "hyperv" {
  user            = "Administator"
  password        = "P@ssw0rd"
  host            = "127.0.0.1"
  port            = 5986
  https           = true
  insecure        = false
  use_ntlm        = true
  tls_server_name = ""
  cacert_path     = ""
  cert_path       = ""
  key_path        = ""
  script_path     = "C:/Temp/terraform_%RAND%.cmd"
  timeout         = "30s"
}

# Create a switch
resource "hyperv_network_switch" "dmz" {
}

# Create a vhd
resource "hyperv_vhd" "webserver" {
}

# Create a machine
resource "hyperv_machine_instance" "webserver" {
}
```

## Argument Reference

The following arguments are supported:

* `user` - (Optional) `Administrator` (default). The username to use when HyperV api calls are made. Generally this is Administrator. It can also be sourced from the `HYPERV_USERNAME` environment variable.

* `password` - (Required) The password associated with the username to use for HyperV api calls. It can also be sourced from the `HYPERV_PASSWORD` environment variable.

* `host` - (Optional) `127.0.0.1` (default). The host to run HyperV api calls against. It can also be sourced from the `HYPERV_HOST` environment variable.

* `port` - (Optional) `5986` (default). The port to run HyperV api calls against. It can also be sourced from the `HYPERV_PORT` environment variable.

* `https` - (Optional) `true` (default). Should https be used for HyperV api calls. It can also be sourced from `HYPERV_HTTPS` environment variable.

* `insecure` - (Optional) `false` (default). Skips TLS Verification for HyperV api calls. Generally this is used for self-signed certificates. Should only be used if absolutely needed. Can also be set via setting the `HYPERV_INSECURE` environment variable to `true`.

* `use_ntlm` - (Optional) `true` (default). Use NTLM for authentication for HyperV api calls. Can also be set via setting the `HYPERV_USE_NTLM` environment variable to `true`.

* `tls_server_name` - (Optional) empty (default). The TLS server name for the host used for HyperV api calls. It can also be sourced from the `HYPERV_TLS_SERVER_NAME` environment variable. Defaults to empty string.

* `cacert_path` - (Optional) empty (default). The path to the ca certificates to use for HyperV api calls. Can also be sourced from the `HYPERV_CACERT_PATH` environment variable.

* `cert_path` - (Optional) empty (default). The path to the certificate to use for authentication for HyperV api calls. Can also be sourced from the `HYPERV_CERT_PATH` environment variable.

* `key_path` - (Optional) empty (default). The path to the certificate private key to use for authentication for HyperV api calls. Can also be sourced from the `HYPERV_KEY_PATH` environment variable.

* `script_path` - (Optional) `C:/Temp/terraform_%RAND%.cmd` (default). The path used to copy scripts meant for remote execution for HyperV api calls. Can also be sourced from the `HYPERV_SCRIPT_PATH` environment variable.

* `timeout` - (Optional) `30s` (default). The timeout to wait for the connection to become available for HyperV api calls. This defaults to 5 minutes. Should be provided as a string like 30s or 5m. Can also be sourced from the `HYPERV_TIMEOUT` environment variable.

## Testing

Credentials must be provided via the `HYPERV_USERNAME`, `HYPERV_PASSWORD` environment variables and the host to run on via `HYPERV_HOST`, in order to run acceptance tests.
