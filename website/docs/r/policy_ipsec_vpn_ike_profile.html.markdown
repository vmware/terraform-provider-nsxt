---
subcategory: "VPN"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_ipsec_vpn_ike_profile"
description: A resource to configure a IPSec VPN IKE Profile.
---

# nsxt_policy_ipsec_vpn_ike_profile

This resource provides a method for the management of a IPSec VPN IKE Profile.

This resource is applicable to NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_ipsec_vpn_ike_profile" "test" {
  display_name          = "test"
  description           = "Terraform provisioned IPSec VPN Ike Profile"
  encryption_algorithms = ["AES_128"]
  digest_algorithms     = ["SHA2_256"]
  dh_groups             = ["GROUP14"]
  ike_version           = "IKE_V2"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `ike_version` - (Optional) Internet Key Exchange(IKE) protocol version to be used, one of `IKE_V1`, `IKE_V2`, `IKE_FLEX`. `IKE_FLEX` will initiate IKE-V2 and respond to both `IKE_V1` and `IKE_V2`.
* `encryption_algorithms` - (Required) Set of encryption algorithms to be used during IKE negotiation.
* `digest_algorithms` - (Required) Set of algorithms to be used for message digest during IKE negotiation. Default is `SHA2_256`.
* `dh_groups` - (Required) Diffie-Hellman group to be used if PFS is enabled. Default is GROUP14.
* `sa_life_time` - (Optional) Life time for security association.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_ipsec_vpn_ike_profile.test UUID
```

The above command imports IPSec VPN IKE Profile named `test` with the NSX ID `UUID`.
