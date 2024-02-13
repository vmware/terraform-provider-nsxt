---
subcategory: "VPN"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_ipsec_vpn_service"
description: A resource to configure a IPSec VPN Service.
---

# nsxt_policy_ipsec_vpn_service

This resource provides a method for the management of a IPSec VPN Service.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_ipsec_vpn_service" "test" {
  display_name  = "ipsec-vpn-service1"
  description   = "Terraform provisioned IPSec VPN service"
  gateway_path  = data.nsxt_policy_tier0_gateway.gw1.path
  enabled       = true
  ha_sync       = true
  ike_log_level = "INFO"
  bypass_rule {
    sources      = ["192.168.10.0/24"]
    destinations = ["192.169.10.0/24"]
    action       = "BYPASS"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `locale_service_path` - (Deprecated) Path of the gateway locale service associated with the IPSec VPN Service.
* `gateway_path` - (Optional) Path of gateway associated with the IPSec VPN Service. Note that at least one of `gateway_path` and `locale_service_path` must be specified for the IPSec VPN Service object.
* `enabled` - (Optional) Whether this IPSec VPN Service is enabled. Default is `true`.
* `ha_sync` - (Optional) Enable/Disable IPSec VPN service HA state sync. Default is `true`.
* `ike_log_level` - (Optional) Set of algorithms to be used for message digest during IKE negotiation. Value is one of `DEBUG`, `INFO`, `WARN`, `ERROR` and `EMERGENCY`. Default is `INFO`. 
* `bypass_rule` - (Optional) Set the bypass rules for this IPSec VPN Service.
  * `sources` - (Required) List of source subnets. Subnet format is ipv4 CIDR.
  * `destinations` - (Required) List of distination subnets. Subnet format is ipv4 CIDR.
  * `action` - (Required) `PROTECT` or `BYPASS`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_ipsec_vpn_service.test PATH
```

The above command imports IPSec VPN Service named `test` that corresponds to NSX IPSec VPN service with policy path `PATH`.
