---
subcategory: "VPN"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_l2_vpn_service"
description: A resource to configure a L2 VPN Service.
---

# nsxt_policy_l2_vpn_service

This resource provides a method for the management of a L2 VPN Service.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_l2_vpn_service" "test" {
  display_name        = "l2-vpn-service1"
  description         = "Terraform provisioned L2 VPN service"
  locale_service_path = data.nsxt_policy_gateway_locale_service.test.path
  enable_hub          = true
  mode                = "SERVER"
  encap_ip_pool       = ["192.168.10.0/24"]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `locale_service_path` - (Deprecated) Path of the locale service associated with the L2 VPN Service.
* `gateway_path` - (Optional) Path of gateway associated with the L2 VPN Service. Note that at least one of `gateway_path` and `locale_service_path` must be specified for the L2 VPN Service object.
* `enable_hub` - (Optional) This property applies only in `SERVER` mode. If set to true, traffic from any client will be replicated to all other clients. If set to false, traffic received from clients is only replicated to the local VPN endpoint. Default is `true`.
* `mode` - (Optional) Specify an L2VPN service mode as SERVER or CLIENT. Value is one of `SERVER`, `CLIENT`. Default is `SERVER`.
* `encap_ip_pool` - (Optional) IP Pool to allocate local and peer endpoint IPs. Format is ipv4 CIDR block.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_l2_vpn_service.test PATH
```

The above command imports L2 VPN Service named `test` that corresponds to NSX L2 VPN service with policy path `PATH`.
