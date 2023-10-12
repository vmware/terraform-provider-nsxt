---
subcategory: "DNS"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_gateway_dns_forwarder"
description: A resource to configure DNS Forwarder on Gateway.
---

# nsxt_policy_gateway_dns_forwarder

This resource provides a method for the management of DNS Forwarder on Tier0 or Tier1 Gateway.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

~> **NOTE:** Only one DNS Forwarder can be configured per Gateway.

~> **NOTE:** DNS Forwarder is only supported on ACTIVE-STANDBY Tier0 Gateways.

## Example Usage

```hcl
resource "nsxt_policy_gateway_dns_forwarder" "test" {
  display_name = "test"
  description  = "Terraform provisioned Zone"
  gateway_path = nsxt_policy_tier1_gateway.test.path
  listener_ip  = "122.30.0.13"
  enabled      = true
  log_level    = "DEBUG"
  cache_size   = 2048

  default_forwarder_zone_path      = nsxt_policy_dns_forwarder_zone.default.path
  conditional_forwarder_zone_paths = [nsxt_policy_dns_forwarder_zone.oranges.path, nsxt_policy_dns_forwarder_zone.apples.path]
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_gateway_dns_forwarder" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "test"
  description  = "Terraform provisioned Zone"
  gateway_path = nsxt_policy_tier1_gateway.test.path
  listener_ip  = "122.30.0.13"
  enabled      = true
  log_level    = "DEBUG"
  cache_size   = 2048

  default_forwarder_zone_path      = nsxt_policy_dns_forwarder_zone.default.path
  conditional_forwarder_zone_paths = [nsxt_policy_dns_forwarder_zone.oranges.path, nsxt_policy_dns_forwarder_zone.apples.path]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `gateway_path` - (Required) Path of Tier0 or Tier1 Gateway.
* `listener_ip` - (Required) IP address on which the DNS Forwarder listens.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to
* `default_forwarder_zone_path` - (Required) Path of Default Forwarder Zone.
* `conditional_forwarder_zone_paths` - (Optional) List of conditional (FQDN) Zone Paths (Maximum 5 zones).
* `enabled` - (Optional) Flag to indicate whether this DNS Forwarder is enabled. Defaults to `true`.
* `log_level` - (Optional) Log Level for related messages, one of `DEBUG`, `INFO`, `WARNING`, `ERROR`, `FATAL`. Defaults to `INFO`.
* `cache_size` - (Optional) Cache size in KB.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_gateway_dns_forwarder.test GATEWAY-PATH
```
The above command imports Dns Forwarder named `test` for NSX Gateway `GATEWAY-PATH`. Note that in order to support both Tier0 and Tier1 Gateways, a full Gateway path is expected here, rather than the usual ID.
