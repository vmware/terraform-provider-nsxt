---
subcategory: "Policy - Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_l2vpn_vpn_session"
description: A resource to configure a L2VPN VPN session.
---

# nsxt_policy_l2vpn_vpn_session

This resource provides a method for the management of a L2VPN VPN session.

This resource is applicable to NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_l2vpn_vpn_session" "test" {
    display_name      = "L2 VPN Session"
    description       = "Terraform-provisioned L2 VPN Tunnel"
    transport_tunnels = [nsxt_policy_ipsec_vpn_session.ipsec_vpn_session_for_l2vpn.path]
}
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `transport_tunnels` - (Optional) List of transport tunnels paths for redundancy.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_l2vpn_vpn_session.test UUID
```

The above command imports l2vpn VPN  session named `test` with the NSX L2VPN VPN Ike session ID `UUID`.
