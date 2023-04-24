---
subcategory: "VPN"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_l2_vpn_session"
description: A resource to configure a L2VPN VPN session.
---

# nsxt_policy_l2vpn_session

This resource provides a method for the management of a L2VPN VPN session.

This resource is applicable to NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_l2_vpn_session" "test" {
  display_name      = "L2 VPN Session"
  description       = "Terraform-provisioned L2 VPN Tunnel"
  service_path      = nsxt_policy_l2_vpn_service.test.path
  transport_tunnels = [nsxt_policy_ipsec_vpn_session.ipsec_vpn_session_for_l2vpn.path]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `service_path` - (Required) The path of the L2 VPN service for the VPN session.
* `transport_tunnels` - (Required) List of transport tunnels paths for redundancy. L2VPN supports only `AES_GCM_128` encryption algorithm for IPSec tunnel profile.
* `direction` - (Optional) The traffic direction apply to the MSS clamping. `BOTH` or `NONE`.
* `max_segment_size` - (Optional) Maximum amount of data the host will accept in a TCP segment. Value should be between `108` and `8860`. If not specified then the value would be the automatic calculated MSS value.
* `local_address` - (Optional) IP Address of the local tunnel port. This property only applies in `CLIENT` mode.
* `peer_address` - (Optional) IP Address of the peer tunnel port. This property only applies in `CLIENT` mode.
* `protocol` - (Optional) Encapsulation protocol used by the tunnel. `GRE` is the only supported value.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_l2_vpn_session.test UUID
```

The above command imports L2 VPN session named `test` with the NSX ID `UUID`.
