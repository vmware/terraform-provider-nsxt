---
subcategory: "Beta"
page_title: "NSXT: nsxt_proxy_config"
description: A data source to read internet proxy configuration information.
---

# nsxt_proxy_config

This data source provides information about internet proxy configuration in NSX Manager.

**ℹ️ Singleton Resource**: This data source reads global proxy configuration which exists as a single instance per NSX-T environment.

This data source is applicable to NSX Manager (version 4.2.0 onwards).

## Example Usage

```hcl
data "nsxt_proxy_config" "internet_proxy" {}

output "proxy_status" {
  value = {
    enabled = data.nsxt_proxy_config.internet_proxy.enabled
    host    = data.nsxt_proxy_config.internet_proxy.host
    port    = data.nsxt_proxy_config.internet_proxy.port
    scheme  = data.nsxt_proxy_config.internet_proxy.scheme
  }
}
```

## Argument Reference

This data source does not require any arguments as it reads the singleton proxy configuration.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the proxy configuration.
* `display_name` - Display name (managed by NSX).
* `description` - Description of the resource.
* `path` - NSX path of the resource.
* `enabled` - Whether proxy is enabled.
* `scheme` - Proxy scheme (HTTP or HTTPS).
* `host` - Proxy server host (IP address or FQDN).
* `port` - Proxy server port.
* `username` - Username for authentication.
* `certificate_id` - Certificate ID for HTTPS proxy.
* `test_connection_url` - URL used to test connectivity.
