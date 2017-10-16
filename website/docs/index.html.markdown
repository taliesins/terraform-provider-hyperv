---
layout: "hyperv"
page_title: "Provider: HyperV"
sidebar_current: "docs-hyperv-index"
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
  user            = "..."
  password        = "..."
  host            = "..."
  port            = "..."
  https           = "..."
  insecure        = "..."
  tls_server_name = "..."
  cacert_path     = "..."
  cert_path       = "..."
  key_path        = "..."
  script_path     = "..."
  timeout         = "..."
}

# Create a switch
resource "hyperv_network_switch" "dmz" {
}

# Create a machine
resource "hyperv_machine_instance" "webserver" {
}
```

## Argument Reference

The following arguments are supported:

* `user` - (Optional) The username to use when HyperV api calls are made. Generally this is Administrator. It can also be sourced from the `HYPERV_USERNAME` environment variable.

* `password` - (Optional) The password associated with the username to use for HyperV api calls. It can also be sourced from the `HYPERV_PASSWORD` environment variable.

* `host` - (Optional) The host to run HyperV api calls against. It can also be sourced from the `HYPERV_HOST` environment variable.

* `port` - (Optional) The port to run HyperV api calls against. It can also be sourced from the `HYPERV_PORT` environment variable.

* `https` - (Optional) Should https be used for HyperV api calls. It can also be sourced from `HYPERV_HTTPS` environment variable.

* `insecure` - (Optional) Skips TLS Verification for HyperV api calls. Generally this is used for self-signed certificates. Should only be used if absolutely needed. Can also be set via setting the `HYPERV_INSECURE` environment variable to `true`.

* `tls_server_name` - (Optional) The TLS server name for the host used for HyperV api calls. It can also be sourced from the `HYPERV_TLS_SERVER_NAME` environment variable. Defaults to empty string.

* `cacert_path` - (Optional) The path to the ca certificates to use for HyperV api calls. Can also be sourced from the `HYPERV_CACERT_PATH` environment variable.

* `cert_path` - (Optional) The path to the certificate to use for authentication for HyperV api calls. Can also be sourced from the `HYPERV_CERT_PATH` environment variable.

* `key_path` - (Optional) The path to the certificate private key to use for authentication for HyperV api calls. Can also be sourced from the `HYPERV_KEY_PATH` environment variable.

* `script_path` - (Optional) The path used to copy scripts meant for remote execution for HyperV api calls. Can also be sourced from the `HYPERV_SCRIPT_PATH` environment variable.

* `timeout` - (Optional) The timeout to wait for the connection to become available for HyperV api calls. This defaults to 5 minutes. Should be provided as a string like 30s or 5m. Can also be sourced from the `HYPERV_TIMEOUT` environment variable.

## Testing

Credentials must be provided via the `HYPERV_USERNAME`, `HYPERV_PASSWORD` environment variables and the host to run on via `HYPERV_HOST`, in order to run acceptance tests.
