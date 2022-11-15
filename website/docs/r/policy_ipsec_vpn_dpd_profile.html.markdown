---
subcategory: "VPN"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_ipsec_vpn_dpd_profile"
description: A resource to configure a IPSec VPN DPD Profile.
---

# nsxt_policy_ipsec_vpn_dpd_profile

This resource provides a method for the management of a IPSec VPN Dead Peer Detection Profile.

This resource is applicable to NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_ipsec_vpn_dpd_profile" "test" {
  display_name       = "test"
  description        = "Terraform provisioned IPSec VPN DPD Profile"
  dpd_probe_mode     = "ON_DEMAND"
  dpd_probe_interval = 120
  enabled            = true
  retry_count        = 8
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `dpd_probe_mode` - (Optional) DPD Probe Mode, one of `PERIODIC`, `ON_DEMAND`. Periodic mode is used to query the liveliness of the peer at regular intervals (`dpd_probe_interval` seconds). On-demand mode would trigger DPD message if there is traffic to send to the peer AND the peer was idle for `dpd_probe_interval` seconds. Default is `PERIODIC`.
* `dpd_probe_interval` - (Optional) Probe interval in seconds. Default is 60.
* `enabled` - (Optional) Enabled Dead Peer Detection. Default is True.
* `retry_count` - (Optional) Maximum number of DPD retry attempts. Default is 10.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_ipsec_vpn_dpd_profile.test UUID
```

The above command imports IPSec VPN DPD Profile named `test` with the NSX ID `UUID`.
