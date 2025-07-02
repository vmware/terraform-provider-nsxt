---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_transit_gateway_ipsec_vpn_service"
description: A resource to configure a Transit Gateway IPSecVpnService.
---

# nsxt_policy_transit_gateway_ipsec_vpn_service

This resource provides a method for the management of a Transit Gateway IPSecVpnService.

This resource is applicable to centralised transit gateway.

## Example Usage

```hcl
resource "nsxt_policy_transit_gateway_ipsec_vpn_service" "test" {
  display_name  = "test"
  description   = "Terraform provisioned IPSecVpnService"
  parent_path   = data.nsxt_policy_transit_gateway.tgw1.path
  enabled       = true
  ha_sync       = true
  ike_log_level = "DEBUG"

}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `parent_path` - (Required) Path of parent object. VPN is supported only on Centralized Transit Gateway and requires ACTIVE-STANDBY HA mode.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_transit_gateway_ipsec_vpn_service.test PATH
```

The above command imports transit gateway IPSecVpnService named `test` with the NSX path `PATH`.
