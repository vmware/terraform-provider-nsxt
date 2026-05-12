---
subcategory: "Beta"
page_title: "NSXT: nsxt_proxy_config"
description: |-
  Manages internet proxy configuration.
---

# nsxt_proxy_config

This resource provides a means to configure internet proxy settings in NSX Manager (System > General Settings > Internet Proxy Server).

**⚠️ Singleton Resource**: This resource represents global proxy configuration and should only be declared **once per NSX-T environment**. The resource uses a fixed ID managed by NSX.

Deleting this resource disables the proxy configuration in NSX rather than removing it. The configuration persists in NSX.
This resource is applicable to NSX Manager (version 4.2.0 onwards).

## Example Usage

### Basic HTTP Proxy (Disabled)

```hcl
resource "nsxt_proxy_config" "test" {
  enabled = false
  scheme  = "HTTP"
  host    = "proxy.example.com"
  port    = 8080
}
```

### HTTP Proxy with Authentication

```hcl
resource "nsxt_proxy_config" "internet_proxy" {
  enabled             = true
  scheme              = "HTTP"
  host                = "proxy.company.com"
  port                = 8080
  username            = "proxyuser"
  password            = var.proxy_password
  test_connection_url = "https://www.vmware.com"

  tag {
    scope = "environment"
    tag   = "production"
  }
}
```

### HTTPS Proxy with Certificate

```hcl
data "nsxt_policy_certificate" "proxy_cert" {
  display_name = "proxy-server-cert"
}

resource "nsxt_proxy_config" "secure_proxy" {
  enabled             = true
  scheme              = "HTTPS"
  host                = "secure-proxy.company.com"
  port                = 3128
  certificate_id      = data.nsxt_policy_certificate.proxy_cert.id
  test_connection_url = "https://www.vmware.com"
}
```

## Argument Reference

The following arguments are supported:

* `enabled` - (Optional) Enable proxy configuration. Default is `false`. When `enabled` is set to `true`, NSX validates connectivity to the proxy server. The test will fail if the proxy is unreachable.
* `scheme` - (Optional) Proxy scheme. Valid values are `HTTP` and `HTTPS`. Default is `HTTP`. When using HTTPS scheme, you may need to provide a `certificate_id`. The certificate must be imported to NSX trust management first.
* `host` - (Optional) Proxy server host (IP address or FQDN). Required when `enabled` is `true`.
* `port` - (Optional) Proxy server port. Default is `3128`. Valid range: 1-65535.
* `username` - (Optional) Username for proxy authentication.
* `password` - (Optional) Password for proxy authentication. This field is sensitive. The password is write-only and not returned by the NSX API. It won't appear in state refresh operations.
* `certificate_id` - (Optional) Certificate ID for HTTPS proxy (from trust-management API). Required when `scheme` is `HTTPS`.
* `test_connection_url` - (Optional) URL to test proxy connectivity. Default is `https://www.vmware.com`.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the proxy configuration (always "TelemetryConfigIdentifier").
* `nsx_id` - NSX ID of the proxy configuration.
* `display_name` - Display name of the proxy configuration (managed by NSX).
* `description` - Description (not supported by NSX for proxy config).
* `path` - The NSX path of the resource (`/api/v1/proxy/config`).
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.

## Importing

An existing proxy configuration can be imported, for example:

```bash
terraform import nsxt_proxy_config.internet_proxy TelemetryConfigIdentifier
```

The above command imports the proxy configuration with the ID `TelemetryConfigIdentifier` (the singleton ID).
